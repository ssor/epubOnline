package epub

import (
	"testing"

	"encoding/json"

	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
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
