package epub

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ssor/epubgo/raw"
	"github.com/ssor/html2text"
)

var (
	metaList = []string{"title", "language", "identifier", "creator", "subject", "description", "publisher", "contributor", "date", "type", "format", "source", "relation", "coverage", "rights"}
)

// LoadEpub will read epub file, and copy it's content to it's dir(md5), and returns epub's info
func LoadEpub(bookPath string) (*Epub, error) {
	zipReader, err := raw.NewEpub(bookPath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	epubFilesDir := ""
	md5, err := caculateMD5Value(bookPath)
	if err != nil {
		return nil, err
	}
	epubFilesDir = md5

	epub, err := NewEpub(zipReader, epubFilesDir)
	if err != nil {
		return nil, err
	}

	return epub, nil
}

func makeupCoverage(currentCoverage, bookDir string, getFileHrefByID func(string) string, f func(string) ([]map[string]string, error)) (coverage string) {
	if len(currentCoverage) <= 0 {
		coverID := getCoverID(f)
		if len(coverID) > 0 {
			coverage = getFileHrefByID(coverID)
		}
	}

	if len(coverage) > 0 {
		coverage = path.Join(bookDir, coverage)
	}
	return
}

func trimSurplusPartOfTitle(current string) string {
	if len(current) > 0 {
		underscoreIndex := strings.Index(current, "_")
		if underscoreIndex > 0 {
			return current[:underscoreIndex]
		}
	}
	return current
}

// NewEpub get metainfo from raw.Epub, and statistics charactor count, and copy files finally
func NewEpub(src *raw.Epub, bookFilesDir string) (*Epub, error) {
	epub := &Epub{
		Navigations: NavigationPointArray{},
		FileDir:     bookFilesDir,
		ID:          bookFilesDir,
	}
	// spew.Dump(src)

	epub.MetaInfo = generateMetaInfo(metaList, src.Metadata)
	epub.MetaInfo["coverage"] = makeupCoverage(epub.MetaInfo["coverage"], bookFilesDir, src.GetFileHrefByID, src.MetadataAttr)
	epub.MetaInfo["title"] = trimSurplusPartOfTitle(epub.MetaInfo["title"])
	epub.Navigations = generateNaviPoints(src.NavPoints(), 1, 1, "", nil)

	err := epub.Navigations.Each(func(nav *NavigationPoint) error {
		navCharactorCount, err := fillNavigationPoint(nav, epub.FileDir, src.OpenFile)
		if err != nil {
			return err
		}

		epub.CharactorCount += navCharactorCount
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

func fillNavigationPoint(nav *NavigationPoint, fileDir string, fileReader func(name string) (io.ReadCloser, error)) (int, error) {
	bs, err := readFileContent(fileReader, nav.URL)
	if err != nil {
		return 0, err
	}
	text, err := html2Text(bs)
	if err != nil {
		return 0, err
	}
	// fmt.Println("content: ")
	// fmt.Println(string(bs))
	// fmt.Println("text: ")
	// fmt.Println(text)
	nav.Text = text

	content := extractHTMLContent(bs)

	nav.URL = path.Join(fileDir, nav.URL)
	nav.CharactorCountSelf = len(content)
	nav.CharactorCountTotal = nav.CharactorCountSelf
	// epub.CharactorCount += nav.CharactorCountSelf
	return nav.CharactorCountSelf, nil
}

// func (e *Epub) copyFiles(zipReader *raw.Epub, extract_dir string) (*Epub, error) {
func (e *Epub) copyFiles(zipReader *raw.Epub) (*Epub, error) {
	files := zipReader.Files()

	for _, file := range files {
		fullPath := path.Join(e.FileDir, file)
		// fullPath := path.Join(extract_dir, e.FileDir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
		if err != nil {
			return nil, err
		}
		content, err := readFileContent(zipReader.OpenFile, file)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(fullPath, content, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func readFileContent(fileReader func(name string) (io.ReadCloser, error), file string) ([]byte, error) {
	closer, err := fileReader(file)
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

func generateMetaInfo(metaList []string, f func(string) ([]string, error)) map[string]string {
	metaInfo := make(map[string]string)

	for _, meta := range metaList {
		ls, err := f(meta)
		if err != nil {
			fmt.Println("[TIP] get meta info: ", err)
			ls = []string{}
		}
		metaInfo[meta] = strings.Join(ls, " ")
	}
	return metaInfo
}

func getCoverID(f func(string) ([]map[string]string, error)) string {

	attributesMeta, err := f("meta")
	if err != nil {
		fmt.Println("[ERR] get meta err: ", err)
		return ""
	}
	// spew.Dump(attributesMeta)
	for _, atr := range attributesMeta {
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

	headNp := nps[0]
	point := NewNavigationPoint(headNp, level, IndexInList, tagPre)
	points = append(points, point)

	if headNp.Children() != nil {
		points = append(points, generateNaviPoints(headNp.Children(), 1, level+1, point.Tag, nil)...)
	}

	return generateNaviPoints(nps[1:], IndexInList+1, level, tagPre, points)
}

// Epub stands for an epub's info
type Epub struct {
	ID             string               `json:"id"`
	Navigations    NavigationPointArray `json:"navigations"`
	MetaInfo       map[string]string    `json:"meta"`
	CharactorCount int                  `json:"charactor_count"`
	URL            string               `json:"url"`
	FileDir        string               `json:"-"`
}

// SetCoverageIfEmpty set a cover if need
func (e *Epub) SetCoverageIfEmpty(cover string) {
	if len(e.MetaInfo["coverage"]) <= 0 {
		// e.MetaInfo["coverage"] = cover
		e.SetCoverage(cover)
	}
}

// SetCoverage set coverage
func (e *Epub) SetCoverage(cover string) {
	e.MetaInfo["coverage"] = cover
}

// Meta returns field value
func (e *Epub) Meta(field string) string {
	if v, exists := e.MetaInfo[field]; exists {
		return v
	}
	return ""
}

// Array is Epub's list
type Array []*Epub

// Find returns finded epub
func (ea Array) Find(p func(*Epub) bool) *Epub {
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

func extractHTMLContent(bs []byte) []rune {
	src := string(bs)
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	bodyHead := strings.Index(src, "<body>")
	bodyTail := strings.Index(src, "</body>")
	if bodyHead >= 0 && bodyTail >= 0 && bodyTail > bodyHead {
		src = src[bodyHead+6 : bodyTail]
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
	srcRune := []rune(src)

	return srcRune
}

func html2Text(bs []byte) (string, error) {
	text, err := html2text.FromReader(bytes.NewReader(bs))
	if err != nil {
		return "", err
	}
	return text, nil
}
