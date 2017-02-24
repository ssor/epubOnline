package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ssor/epubOnline/epub"

	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

var (
	configFile      = flag.String("config", "conf/config.json", "config file for system")
	listeningPort   = flag.String("port", "8092", "listeningPort")
	book_dir        = "books_raw"
	book_online_dir = "books"
	books           = epub.EpubArray{}

	app_dir = []string{"books"}

	default_coverage = "/images/bb.png"
)

func main() {

	flag.Parse()
	if flag.Parsed() == false {
		flag.PrintDefaults()
		return
	}
	os.RemoveAll(app_dir[0])
	initAppDir(app_dir)
	initBooks()

	router := gin.Default()
	router.Static("/epub", book_dir)
	router.Static("/"+book_online_dir, book_online_dir)
	router.Static("/js", "static/js")
	router.Static("/images", "static/img")
	router.Static("/css", "static/css")
	router.Static("/bootstrap", "bootstrap")

	router.LoadHTMLGlob("views/*.html")

	router.GET("/books", func(c *gin.Context) {
		c.JSON(http.StatusOK, books)
	})

	router.GET("/book", func(c *gin.Context) {
		id := c.Query("id")
		c.JSON(http.StatusOK, books.Find(func(e *epub.Epub) bool {
			return e.ID == id
		}))
	})

	router.GET("/readBookIndex", func(c *gin.Context) {
		id := c.Query("id")
		c.HTML(http.StatusOK, "book_nav.html", gin.H{"ID": id})
	})
	router.GET("/BookNavIndex", func(c *gin.Context) {
		id := c.Query("id")
		c.HTML(http.StatusOK, "book_nav.html", gin.H{"ID": id})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.Run(":" + *listeningPort)
}

func initBooks() {
	files, err := getEpubFiles(book_dir)
	if err != nil {
		panic(err)
	}
	books, err = InitEpub(files)
	if err != nil {
		panic(err)
	}
}

func getEpubFiles(dir string) ([]string, error) {
	files := []string{}

	file_info_list, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("[ERR] read dir err: ", err)
		return files, err
	}

	for _, file_info := range file_info_list {
		// spew.Dump(file_info)
		name := file_info.Name()
		if isEpubFormatFile(name) {
			// fmt.Println(" - ", name)
			files = append(files, path.Join(book_dir, name))
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

func InitEpub(files []string) (epub.EpubArray, error) {
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
		err = epub.MoveEpub(book_online_dir, epub_book)
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

//创建程序运行的必需基础目录
func initAppDir(dirList []string) {
	for _, dir := range dirList {
		if b := isFileExist(dir); b == true {
			continue
		}

		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

// exists returns whether the given file or directory exists or not
func isFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// /static_server/html/resource/books/src
// rmcbs_1487300029118047270.epub
//  rmcbs_1487222550599695691.epub
// rmcbs_1487222638914027665.epub
