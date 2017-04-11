package api

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/ssor/epub_online/epub"
)

var (
	books = epub.Array{}
)

// InitBooks check book list and extract content from zip
func InitBooks(epubSrcDir, destDir, defaultCoverage string) {
	files, err := listEpubNames(epubSrcDir)
	if err != nil {
		panic(err)
	}

	filesFullPath := []string{}
	for _, name := range files {
		filesFullPath = append(filesFullPath, path.Join(epubSrcDir, name))
	}

	books, err = InitEpub(filesFullPath, destDir, defaultCoverage)
	if err != nil {
		panic(err)
	}
}

func listEpubNames(epubSrcDir string) ([]string, error) {
	files := []string{}

	fileInfoList, err := ioutil.ReadDir(epubSrcDir)
	if err != nil {
		fmt.Println("[ERR] read dir err: ", err)
		return files, err
	}

	for _, fileInfo := range fileInfoList {
		// spew.Dump(fileInfo)
		name := fileInfo.Name()
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

// InitEpub init an epub book, includes load it's info, move it's extracted files to dest dir
func InitEpub(files []string, destDir, defaultCoverage string) (epub.Array, error) {
	books := epub.Array{}
	for _, name := range files {
		fmt.Println(" - ", name)
		// epub_book, err := epub.LoadEpub(name, book_online_dir)
		epubBook, err := epub.LoadEpub(name)
		if err != nil {
			fmt.Println("[ERR] load book ", name, " err: ", err)
			return nil, err
		}
		epubBook.URL = path.Join("epub", path.Base(name))
		err = epub.MoveEpub(destDir, epubBook)
		if err != nil {
			fmt.Println("[ERR] move book ", name, " err: ", err)
			return nil, err
		}
		epubBook.SetCoverageIfEmpty(defaultCoverage)
		books = append(books, epubBook)
	}

	fmt.Println(len(books), " books loaded:")
	for _, book := range books {
		fmt.Println("-- ", book.Meta("title"))
	}
	return books, nil

}
