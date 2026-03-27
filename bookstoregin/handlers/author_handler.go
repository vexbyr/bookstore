package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"bookstoregin/models"
)

var Authors = []models.Author{}
var AuthorID = 1

func SeedAuthors() {
	Authors = append(Authors,
		models.Author{ID: 1, Name: "J.K. Rowling"},
		models.Author{ID: 2, Name: "George R.R. Martin"},
	)
	AuthorID = 3
}

func GetAuthors(c *gin.Context) {
	c.JSON(http.StatusOK, Authors)
}

func CreateAuthor(c *gin.Context) {
	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	author.ID = AuthorID
	AuthorID++
	Authors = append(Authors, author)
	c.JSON(http.StatusCreated, author)
}