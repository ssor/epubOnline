package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func ReadBookIndex(c *gin.Context) {
	id := c.Query("id")
	c.HTML(http.StatusOK, "book_nav.html", gin.H{"ID": id})
}
func BookNavIndex(c *gin.Context) {
	id := c.Query("id")
	c.HTML(http.StatusOK, "book_nav.html", gin.H{"ID": id})
}
