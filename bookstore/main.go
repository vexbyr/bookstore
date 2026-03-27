package main

import (
	"log"
	"net/http"

	"bookstore/handlers"
	"bookstore/models"
)

func main() {
	// === Инициализация тестовых данных ===
	handlers.AuthorsInit([]models.Author{
		{Name: "Isaac Asimov"},
		{Name: "J.K. Rowling"},
		{Name: "George Orwell"},
	})

	handlers.CategoriesInit([]models.Category{
		{Name: "Science Fiction"},
		{Name: "Fantasy"},
		{Name: "Dystopia"},
	})

	handlers.BooksInit([]models.Book{
		{Title: "Foundation", AuthorID: 1, CategoryID: 1, Price: 12.99},
		{Title: "Harry Potter", AuthorID: 2, CategoryID: 2, Price: 9.99},
		{Title: "1984", AuthorID: 3, CategoryID: 3, Price: 11.50},
	})

	// === Роуты ===
	http.HandleFunc("/books", handlers.BooksHandler)
	http.HandleFunc("/books/", handlers.BookByIDHandler)

	http.HandleFunc("/authors", handlers.AuthorsHandler)
	http.HandleFunc("/categories", handlers.CategoriesHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}