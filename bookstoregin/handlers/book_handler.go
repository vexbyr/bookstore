package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"bookstoregin/models"
)

var Books = []models.Book{}
var BookID = 1

func SeedBooks() {
	Books = append(Books,
		models.Book{ID: 1, Title: "Harry Potter", AuthorID: 1, CategoryID: 1, Price: 15.5},
		models.Book{ID: 2, Title: "Game of Thrones", AuthorID: 2, CategoryID: 1, Price: 20},
	)
	BookID = 3
}

func GetBooks(c *gin.Context) {
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 5
	start := (page - 1) * limit
	end := start + limit

	filtered := []models.Book{}
	for _, b := range Books {
		if category == "" || strconv.Itoa(b.CategoryID) == category {
			filtered = append(filtered, b)
		}
	}

	if start > len(filtered) {
		c.JSON(http.StatusOK, []models.Book{})
		return
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	c.JSON(http.StatusOK, filtered[start:end])
}

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book.ID = BookID
	BookID++
	Books = append(Books, book)
	c.JSON(http.StatusCreated, book)
}

func GetBookByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, b := range Books {
		if b.ID == id {
			c.JSON(http.StatusOK, b)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func UpdateBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updated models.Book
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, b := range Books {
		if b.ID == id {
			updated.ID = id
			Books[i] = updated
			c.JSON(http.StatusOK, updated)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func DeleteBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, b := range Books {
		if b.ID == id {
			Books = append(Books[:i], Books[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}