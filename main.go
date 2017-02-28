package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ssor/epubOnline/api"
	"github.com/ssor/epubOnline/controller"

	"flag"
)

var (
	configFile      = flag.String("config", "conf/config.json", "config file for system")
	listeningPort   = flag.String("port", "8092", "listeningPort")
	book_dir        = "books_raw"
	book_online_dir = "books"

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
	api.InitBooks(book_dir, book_online_dir, default_coverage)

	router := gin.Default()
	router.Static("/epub", book_dir)
	router.Static("/"+book_online_dir, book_online_dir)
	router.Static("/js", "static/js")
	router.Static("/images", "static/img")
	router.Static("/css", "static/css")
	router.Static("/bootstrap", "bootstrap")

	router.LoadHTMLGlob("views/*.html")

	router.GET("/books", api.Books)
	router.GET("/book", api.Book)

	router.GET("/", controller.Index)
	router.GET("/readBookIndex", controller.ReadBookIndex)
	router.GET("/BookNavIndex", controller.BookNavIndex)

	router.Run(":" + *listeningPort)
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
