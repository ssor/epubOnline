package epub

import (
	"testing"

	"encoding/json"

	"fmt"
	"strings"

	"io/ioutil"
	"path"

	"github.com/davecgh/go-spew/spew"
	"github.com/ssor/epubOnline/bom"
)

const (
	bookPath_gcdxy     = "../books_raw/gcdxy.epub"
	bookPath_zrddxcxyj = "../books_raw/zrddxcxyj.epub"

	destPath = "testdata"
)

var (
	books_path = []string{
		bookPath_gcdxy, bookPath_zrddxcxyj,
	}
)

func TestHtml2Text(t *testing.T) {
	html_files := []struct {
		file                      string
		keyword_should_not_exists string
		keyword_should_exists     string
	}{
		{

			"html2text2.html",
			"学习之道:美国公认学习第一书NOBOM",
			"",
		},
		{

			"html2text.html",
			"学习之道:美国公认学习第一书title",
			"次世界冠军赛上，我几近疯狂",
		},
		{
			"chapter_00008.xhtml",
			"1892年波兰文版序言title",
			"种新的波兰文本已成为必要",
		},
		{
			"chapter_9.xhtml",
			"1892年波兰文版序言title",
			"",
		},
	}

	for _, html_file := range html_files {
		bs, err := ioutil.ReadFile(path.Join(destPath, html_file.file))
		if err != nil {
			t.Fatal("ReadFile  err: ", err)
		}
		bs = bom.CleanBom(bs)
		// fmt.Println(string(bs))
		text, err := html2Text(bs)
		if err != nil {
			t.Fatal("html2Text  err: ", err)
		}
		// t.Log(" ******", html_file, " text: ******")
		// t.Log(text)
		if strings.Contains(text, html_file.keyword_should_exists) == false {
			t.Fatal("keyword ", html_file.keyword_should_exists, " should  exists in file ", html_file.file)
		}
		if strings.Contains(text, html_file.keyword_should_not_exists) == true {
			t.Fatal("keyword ", html_file.keyword_should_not_exists, " should not exists in file ", html_file.file)
		}
	}
}

func TestEpub(t *testing.T) {
	for _, book_path := range books_path {
		openBook(book_path, t)
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
		fmt.Printf("%-8s  %-10d %-10d %10s %s\n", nav.Tag, nav.CharactorCountSelf, nav.CharactorCountTotal, nav.Url, strings.Repeat("-", nav.Level*4)+nav.Title)

		if isFileExist(nav.Url) == false {
			t.Fatal(nav.Url, " should exits")
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
