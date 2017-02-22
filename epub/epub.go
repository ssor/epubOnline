package epub

import (
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/ssor/epubgo/raw"
)

var ()

func LoadEpub(bookPath string) (*Epub, error) {
	zipReader, err := raw.NewEpub(bookPath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()
	epub, err := NewEpub(zipReader)
	if err != nil {
		return nil, err
	}
	return epub, nil
}

func NewEpub(src *raw.Epub) (*Epub, error) {
	epub := &Epub{
		Navigations: NavigationPointArray{},
	}
	// spew.Dump(src)

	meta_list := []string{"title", "language", "identifier", "creator", "subject", "description", "publisher", "contributor", "date", "type", "format", "source", "relation", "coverage", "rights"}
	epub.MetaInfo = generateMetaInfo(meta_list, src.Metadata)

	if len(epub.MetaInfo["coverage"]) <= 0 {
		cover_id := getCover(src.MetadataAttr)
		epub.MetaInfo["coverage"] = src.GetFileHrefByID(cover_id)
	}

	epub.Navigations = generateNaviPoints(src.NavPoints(), 1, 1, "", nil)

	err := epub.Navigations.Each(func(nav *NavigationPoint) error {
		closer, err := src.OpenFile(nav.Url)
		if err != nil {
			fmt.Println("[ERR] open file err: ", err)
			return err
		}
		defer closer.Close()
		content, err := getHtmlContent(closer)
		if err != nil {
			return err
		}
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

	return epub, nil
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
	Navigations    NavigationPointArray `json:"navigations"`
	MetaInfo       map[string]string    `json:"meta"`
	CharactorCount int                  `json:"charactor_count"`
	Url            string               `json:"url"`
}

func (e *Epub) Meta(field string) string {
	if v, exists := e.MetaInfo[field]; exists {
		return v
	} else {
		return ""
	}
}

type EpubArray []*Epub

func getHtmlContent(reader io.Reader) ([]rune, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	src := string(b)
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

	return src_rune, nil
}
