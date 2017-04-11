package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssor/epub_online/epub"
)

// Books returns all books
func Books(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

// Book returns specified book of id
func Book(c *gin.Context) {
	id := c.Query("id")
	c.JSON(http.StatusOK, books.Find(func(e *epub.Epub) bool {
		return e.ID == id
	}))
}
