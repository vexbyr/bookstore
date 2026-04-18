package main

import (
	"log"

	"bookstore_gin/handlers"
	"bookstore_gin/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	authorHandler := handlers.NewAuthorHandler()
	categoryHandler := handlers.NewCategoryHandler()
	bookHandler := handlers.NewBookHandler(authorHandler, categoryHandler)
	favoriteHandler := handlers.NewFavoriteHandler(bookHandler)

	r := gin.Default()

	// ── Public routes ──────────────────────────────────────────────────────────
	authors := r.Group("/authors")
	{
		authors.GET("", authorHandler.GetAll)
		authors.POST("", authorHandler.Create)
	}

	categories := r.Group("/categories")
	{
		categories.GET("", categoryHandler.GetAll)
		categories.POST("", categoryHandler.Create)
	}

	// All /books routes live in one group so the wildcard name is consistent.
	// Gin resolves static segments ("favorites") before wildcards (":id"),
	// so GET /books/favorites will never be caught by GET /books/:id.
	books := r.Group("/books")
	{
		// Public book CRUD
		books.GET("", bookHandler.GetAll)
		books.POST("", bookHandler.Create)
		books.GET("/:id", bookHandler.GetByID)
		books.PUT("/:id", bookHandler.Update)
		books.DELETE("/:id", bookHandler.Delete)

		// Protected favorites — same :id wildcard, JWT middleware per-route
		auth := middleware.RequireAuth()
		books.GET("/favorites", auth, favoriteHandler.GetFavorites)
		books.PUT("/:id/favorites", auth, favoriteHandler.AddFavorite)
		books.DELETE("/:id/favorites", auth, favoriteHandler.RemoveFavorite)
	}

	// Print a dev token for testing
	if token, err := middleware.GenerateToken(1); err == nil {
		log.Println("----------------------------------------------------")
		log.Println("Dev JWT for user_id=1 (valid 24 h):")
		log.Println("  " + token)
		log.Println("  Usage: Authorization: Bearer <token>")
		log.Println("----------------------------------------------------")
	}

	log.Println("Bookstore API running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
