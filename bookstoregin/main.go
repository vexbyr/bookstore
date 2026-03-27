package main

import (
	"bookstoregin/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	handlers.SeedAuthors()
	handlers.SeedCategories()
	handlers.SeedBooks()

	r.GET("/books", handlers.GetBooks)
	r.POST("/books", handlers.CreateBook)
	r.GET("/books/:id", handlers.GetBookByID)
	r.PUT("/books/:id", handlers.UpdateBook)
	r.DELETE("/books/:id", handlers.DeleteBook)

	r.GET("/authors", handlers.GetAuthors)
	r.POST("/authors", handlers.CreateAuthor)

	r.GET("/categories", handlers.GetCategories)
	r.POST("/categories", handlers.CreateCategory)

	r.Run(":8080")
}