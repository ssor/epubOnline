package epub

import (
	"testing"

	"encoding/json"

	"fmt"
	"strings"

	"io/ioutil"
	"path"

	"github.com/davecgh/go-spew/spew"
)

const (
	bookPathGCDXY     = "../books_raw/gcdxy.epub"
	bookPathZRDDXCXYJ = "../books_raw/zrddxcxyj.epub"

	destPath = "testdata"
)

var (
	booksPath = []string{
		bookPathGCDXY, bookPathZRDDXCXYJ,
	}
)

func TestHtml2Text(t *testing.T) {
	htmlFiles := []struct {
		file                   string
		keywordShouldNotExists string
		keywordShouldExists    string
	}{
		{

			"utf8.html",
			"学习之道:美国公认学习第一书title",
			"次世界冠军赛上，我几近疯狂",
		},
		{
			"utf8_with_bom.xhtml",
			"1892年波兰文版序言title",
			"种新的波兰文本已成为必要",
		},
	}

	for _, htmlFile := range htmlFiles {
		bs, err := ioutil.ReadFile(path.Join(destPath, htmlFile.file))
		if err != nil {
			t.Fatal("ReadFile  err: ", err)
		}
		text, err := html2Text(bs)
		if err != nil {
			t.Fatal("html2Text  err: ", err)
		}
		// t.Log(" ******", htmlFile, " text: ******")
		// t.Log(text)
		if strings.Contains(text, htmlFile.keywordShouldExists) == false {
			t.Fatal("keyword ", htmlFile.keywordShouldExists, " should  exists in file ", htmlFile.file)
		}
		if strings.Contains(text, htmlFile.keywordShouldNotExists) == true {
			t.Fatal("keyword ", htmlFile.keywordShouldNotExists, " should not exists in file ", htmlFile.file)
		}
	}
}

func TestEpub(t *testing.T) {
	for _, bookPath := range booksPath {
		openBook(bookPath, t)
	}
}

func openBook(bookPath string, t *testing.T) {
	epub, err := LoadEpub(bookPath)
	if err != nil {
		t.Fatal("init epub err: ", err)
	}
	err = MoveEpub(destPath, epub)
	if err != nil {
		t.Fatal("move epub err: ", err)
	}

	fmt.Println("*** ", epub.Meta("title"), " 总字数: ", epub.CharactorCount)
	fmt.Printf("%-8s  %-6s %-6s %-15s %s\n", "", "本页字数", "章节总字数", "文件", "目录名称")
	for _, nav := range epub.Navigations {
		fmt.Printf("%-8s  %-10d %-10d %10s %s\n", nav.Tag, nav.CharactorCountSelf, nav.CharactorCountTotal, nav.URL, strings.Repeat("-", nav.Level*4)+nav.Title)

		if isFileExist(nav.URL) == false {
			t.Fatal(nav.URL, " should exits")
		}
	}
	if len(epub.MetaInfo["coverage"]) > 0 {
		if isFileExist(epub.MetaInfo["coverage"]) == false {
			t.Fatal(epub.MetaInfo["coverage"], " should exits")
		}
	}

}

func TestMapJson(t *testing.T) {

	type example struct {
		Obj   map[string]string
		Title string
	}
	m := make(map[string]string)

	m["a"] = "aaa"
	m["b"] = "bbb"

	spew.Dump(m)
	obj := example{
		Obj:   m,
		Title: "example",
	}

	bs, _ := json.Marshal(obj)
	spew.Dump(string(bs))
}

func TestSwitch(t *testing.T) {
	loop := 1

	switch loop {
	default:
		println("default")
	case 1:
		println("11111111")
	}
}
