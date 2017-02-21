package main

import (
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
	configFile    = flag.String("config", "conf/config.json", "config file for system")
	listeningPort = flag.String("port", "8092", "listeningPort")
	book_dir      = "books/"
	books         = epub.EpubArray{}
)

func main() {

	flag.Parse()
	if flag.Parsed() == false {
		flag.PrintDefaults()
		return
	}
	InitEpub()

	router := gin.Default()
	router.Static("/javascripts", "static/js")
	router.Static("/epub", "books")
	// router.Static("/images", "static/img")
	// router.Static("/stylesheets", "static/css")

	// router.LoadHTMLGlob("views/*.tpl")

	router.GET("/books", func(c *gin.Context) {
		c.JSON(http.StatusOK, books)
	})
	router.Run(":" + *listeningPort)
}

func InitEpub() {

	file_info_list, err := ioutil.ReadDir(book_dir)
	if err != nil {
		fmt.Println("[ERR] read dir err: ", err)
		return
	}
	for _, file_info := range file_info_list {
		// spew.Dump(file_info)
		name := file_info.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if path.Ext(name) == ".epub" {
			fmt.Println(" - ", name)
			epub_book, err := epub.LoadEpub(path.Join(book_dir, name))
			if err != nil {
				fmt.Println("[ERR] load book err: ", err)
				continue
			}
			epub_book.Url = path.Join("epub", name)
			books = append(books, epub_book)
		}
	}

	fmt.Println(len(books), " books loaded:")
	for _, book := range books {
		fmt.Println("-- ", book.Meta("title"))
	}

}

// /static_server/html/resource/books/src
// rmcbs_1487300029118047270.epub
//  rmcbs_1487222550599695691.epub
// rmcbs_1487222638914027665.epub
