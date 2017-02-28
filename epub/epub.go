package epub

import (
	"io/ioutil"
	"regexp"
	"strings"

	"fmt"

	"os"
	"path"

	"path/filepath"

	"bytes"

	"github.com/davecgh/go-spew/spew"
	"github.com/ssor/epubgo/raw"
	"github.com/ssor/html2text"
)

var ()

// func LoadEpub(bookPath string, extractDir string) (*Epub, error) {
func LoadEpub(bookPath string) (*Epub, error) {
	zipReader, err := raw.NewEpub(bookPath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	epub_files_dir := ""
	// if len(extractDir) > 0 {
	md5, err := caculateMD5Value(bookPath)
	if err != nil {
		return nil, err
	}
	epub_files_dir = md5
	// }

	// epub, err := NewEpub(zipReader, epub_files_dir, extractDir)
	epub, err := NewEpub(zipReader, epub_files_dir)
	if err != nil {
		return nil, err
	}

	return epub, nil
}

// func NewEpub(src *raw.Epub, book_files_dir, extract_dir string) (*Epub, error) {
func NewEpub(src *raw.Epub, book_files_dir string) (*Epub, error) {
	epub := &Epub{
		Navigations: NavigationPointArray{},
		FileDir:     book_files_dir,
		ID:          book_files_dir,
	}
	// spew.Dump(src)

	meta_list := []string{"title", "language", "identifier", "creator", "subject", "description", "publisher", "contributor", "date", "type", "format", "source", "relation", "coverage", "rights"}
	epub.MetaInfo = generateMetaInfo(meta_list, src.Metadata)

	if len(epub.MetaInfo["coverage"]) <= 0 {
		cover_id := getCover(src.MetadataAttr)
		if len(cover_id) > 0 {
			epub.MetaInfo["coverage"] = path.Join(src.GetFileHrefByID(cover_id))
			// epub.MetaInfo["coverage"] = path.Join(book_files_dir, src.GetFileHrefByID(cover_id))
		}
	}

	if len(epub.MetaInfo["coverage"]) > 0 {
		epub.MetaInfo["coverage"] = path.Join(book_files_dir, epub.MetaInfo["coverage"])
	}

	title := epub.MetaInfo["title"]
	if len(title) > 0 {
		underscore_index := strings.Index(title, "_")
		if underscore_index > 0 {
			epub.MetaInfo["title"] = title[:underscore_index]
		}
	}

	epub.Navigations = generateNaviPoints(src.NavPoints(), 1, 1, "", nil)

	err := epub.Navigations.Each(func(nav *NavigationPoint) error {
		bs, err := readFileContent(src, nav.Url)
		if err != nil {
			return err
		}
		text, err := html2Text(bs)
		if err != nil {
			return err
		}
		nav.Text = text

		content := getHtmlContent(bs)

		nav.Url = path.Join(epub.FileDir, nav.Url)
		nav.CharactorCountSelf = len(content)
		nav.CharactorCountTotal = nav.CharactorCountSelf
		epub.CharactorCount += nav.CharactorCountSelf
		return nil
	})
	if err != nil {
		return nil, err
	}

	epub.Navigations.Each(func(nav *NavigationPoint) error {
		nav.CharactorCountTotal = epub.Navigations.SumSubLevelCharactorCount(nav.Tag) + nav.CharactorCountSelf
		return nil
	})

	// return epub.copyFiles(src, extract_dir)
	return epub.copyFiles(src)
}

// func (e *Epub) copyFiles(zipReader *raw.Epub, extract_dir string) (*Epub, error) {
func (e *Epub) copyFiles(zipReader *raw.Epub) (*Epub, error) {
	files := zipReader.Files()

	for _, file := range files {
		full_path := path.Join(e.FileDir, file)
		// full_path := path.Join(extract_dir, e.FileDir, file)
		err := os.MkdirAll(filepath.Dir(full_path), os.ModePerm)
		if err != nil {
			return nil, err
		}
		content, err := readFileContent(zipReader, file)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(full_path, content, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func readFileContent(zipReader *raw.Epub, file string) ([]byte, error) {
	closer, err := zipReader.OpenFile(file)
	if err != nil {
		fmt.Println("[ERR] open file err: ", err)
		return nil, err
	}
	defer closer.Close()
	content, err := ioutil.ReadAll(closer)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func generateMetaInfo(meta_list []string, f func(string) ([]string, error)) map[string]string {
	meta_info := make(map[string]string)

	for _, meta := range meta_list {
		ls, err := f(meta)
		if err != nil {
			fmt.Println("[TIP] get meta info: ", err)
			ls = []string{}
		}
		meta_info[meta] = strings.Join(ls, " ")
	}
	return meta_info
}

func getCover(f func(string) ([]map[string]string, error)) string {

	attributes_meta, err := f("meta")
	if err != nil {
		fmt.Println("[ERR] get meta err: ", err)
		return ""
	}
	spew.Dump(attributes_meta)
	for _, atr := range attributes_meta {
		if name, exists := atr["name"]; exists && name == "cover" {
			// fmt.Println("name -> ", name)
			if content, exists := atr["content"]; exists {
				// fmt.Println("content -> ", content)
				return content
			}
		}
	}
	return ""
}

func generateNaviPoints(nps raw.NavPointArray, IndexInList, level int, tagPre string, points NavigationPointArray) NavigationPointArray {
	if points == nil {
		points = NavigationPointArray{}
	}
	if nps == nil || len(nps) <= 0 {
		return points
	}

	head_np := nps[0]
	point := NewNavigationPoint(head_np, level, IndexInList, tagPre)
	points = append(points, point)

	if head_np.Children() != nil {
		points = append(points, generateNaviPoints(head_np.Children(), 1, level+1, point.Tag, nil)...)
	}

	return generateNaviPoints(nps[1:], IndexInList+1, level, tagPre, points)
}

type Epub struct {
	ID             string               `json:"id"`
	Navigations    NavigationPointArray `json:"navigations"`
	MetaInfo       map[string]string    `json:"meta"`
	CharactorCount int                  `json:"charactor_count"`
	Url            string               `json:"url"`
	FileDir        string               `json:"-"`
}

func (e *Epub) SetCoverageIfEmpty(cover string) {
	if len(e.MetaInfo["coverage"]) <= 0 {
		e.MetaInfo["coverage"] = cover
	}
}

func (e *Epub) SetCoverage(cover string) {
	e.MetaInfo["coverage"] = cover
}
func (e *Epub) Meta(field string) string {
	if v, exists := e.MetaInfo[field]; exists {
		return v
	} else {
		return ""
	}
}

type EpubArray []*Epub

func (ea EpubArray) Find(p func(*Epub) bool) *Epub {
	if ea == nil {
		return nil
	}
	for _, e := range ea {
		if p(e) {
			return e
		}
	}
	return nil
}

func getHtmlContent(bs []byte) []rune {
	src := string(bs)
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	body_head := strings.Index(src, "<body>")
	body_tail := strings.Index(src, "</body>")
	if body_head >= 0 && body_tail >= 0 && body_tail > body_head {
		src = src[body_head+6 : body_tail]
	}
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	src = strings.Replace(src, "\n", "", -1)
	src = strings.TrimSpace(src)
	src_rune := []rune(src)

	return src_rune
}

func html2Text(bs []byte) (string, error) {
	text, err := html2text.FromReader(bytes.NewReader(bs))
	if err != nil {
		return "", err
	}
	return text, nil
}
