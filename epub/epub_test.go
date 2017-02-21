package epub

import (
	"testing"

	"encoding/json"

	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/ssor/epubgo/raw"
)

const (
	bookPath_gcdxy     = "../books/gcdxy.epub"
	bookPath_zrddxcxyj = "../books/zrddxcxyj.epub"
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
	zipReader, err := raw.NewEpub(bookPath)
	if err != nil {
		t.Fatalf("Open(%v) return an error: %v", bookPath, err)
	}
	defer zipReader.Close()
	epub, err := NewEpub(zipReader)
	if err != nil {
		t.Fatal("init epub err: ", err)
	}
	// spew.Dump(epub)

	fmt.Printf("%-8s  %-6s %-6s %-15s %s\n", "", "本页字数", "章节总字数", "文件", "目录名称")
	for _, nav := range epub.Navigations {
		fmt.Printf("%-8s  %-10d %-10d %10s %s\n", nav.Tag, nav.CharactorCountSelf, nav.CharactorCountTotal, nav.Url, strings.Repeat("-", nav.Level*4)+nav.Title)
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
