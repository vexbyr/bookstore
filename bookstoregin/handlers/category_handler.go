package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"bookstoregin/models"
)

var Categories = []models.Category{}
var CategoryID = 1

func SeedCategories() {
	Categories = append(Categories,
		models.Category{ID: 1, Name: "Fantasy"},
		models.Category{ID: 2, Name: "Drama"},
	)
	CategoryID = 3
}

func GetCategories(c *gin.Context) {
	c.JSON(http.StatusOK, Categories)
}

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category.ID = CategoryID
	CategoryID++
	Categories = append(Categories, category)
	c.JSON(http.StatusCreated, category)
}