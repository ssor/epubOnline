package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssor/epubOnline/epub"
)

func Books(c *gin.Context) {

	c.JSON(http.StatusOK, books)
}
func Book(c *gin.Context) {
	id := c.Query("id")
	c.JSON(http.StatusOK, books.Find(func(e *epub.Epub) bool {
		return e.ID == id
	}))
}
