package api

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/ssor/epubOnline/epub"
)

var (
	books = epub.EpubArray{}
)

func InitBooks(epub_src_dir, dest_dir, default_coverage string) {
	files, err := listEpubNames(epub_src_dir)
	if err != nil {
		panic(err)
	}

	files_full_path := []string{}
	for _, name := range files {
		files_full_path = append(files_full_path, path.Join(epub_src_dir, name))
	}

	books, err = InitEpub(files_full_path, dest_dir, default_coverage)
	if err != nil {
		panic(err)
	}
}

func listEpubNames(epub_src_dir string) ([]string, error) {
	files := []string{}

	file_info_list, err := ioutil.ReadDir(epub_src_dir)
	if err != nil {
		fmt.Println("[ERR] read dir err: ", err)
		return files, err
	}

	for _, file_info := range file_info_list {
		// spew.Dump(file_info)
		name := file_info.Name()
		if isEpubFormatFile(name) {
			// fmt.Println(" - ", name)
			// files = append(files, path.Join(epub_content_file_dir, name))
			files = append(files, name)
		}
	}
	return files, nil
}

func isEpubFormatFile(file string) bool {
	if strings.HasPrefix(file, ".") {
		return false
	}
	return path.Ext(file) == ".epub"
}

func InitEpub(files []string, dest_dir, default_coverage string) (epub.EpubArray, error) {
	books := epub.EpubArray{}
	for _, name := range files {
		fmt.Println(" - ", name)
		// epub_book, err := epub.LoadEpub(name, book_online_dir)
		epub_book, err := epub.LoadEpub(name)
		if err != nil {
			fmt.Println("[ERR] load book ", name, " err: ", err)
			return nil, err
		}
		epub_book.Url = path.Join("epub", path.Base(name))
		err = epub.MoveEpub(dest_dir, epub_book)
		if err != nil {
			fmt.Println("[ERR] move book ", name, " err: ", err)
			return nil, err
		}
		epub_book.SetCoverageIfEmpty(default_coverage)
		books = append(books, epub_book)
	}

	fmt.Println(len(books), " books loaded:")
	for _, book := range books {
		fmt.Println("-- ", book.Meta("title"))
	}
	return books, nil

}
