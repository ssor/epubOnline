package main

import (
	"github.com/gin-gonic/gin"

	"flag"
)

var (
	configFile    = flag.String("config", "conf/config.json", "config file for system")
	listeningPort = flag.String("port", "8092", "listeningPort")
)

func main() {

	flag.Parse()
	if flag.Parsed() == false {
		flag.PrintDefaults()
		return
	}

	router := gin.Default()
	router.Static("/javascripts", "static/js")
	// router.Static("/images", "static/img")
	// router.Static("/stylesheets", "static/css")

	router.LoadHTMLGlob("views/*.tpl")

	router.Run(":" + *listeningPort)
}

// /static_server/html/resource/books/src
// rmcbs_1487300029118047270.epub
//  rmcbs_1487222550599695691.epub
// rmcbs_1487222638914027665.epub
