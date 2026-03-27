package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"bookstore/models"
)

var books = make(map[int]models.Book)
var nextBookID = 1

func BooksInit(list []models.Book) {
	for _, b := range list {
		b.ID = nextBookID
		nextBookID++
		books[b.ID] = b
	}
}
func BooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// фильтр по категории
		categoryFilter := r.URL.Query().Get("category")
		var list []models.Book
		for _, b := range books {
			if categoryFilter != "" && strconv.Itoa(b.CategoryID) != categoryFilter {
				continue
			}
			list = append(list, b)
		}

		// пагинация
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page <= 0 {
			page = 1
		}
		limit := 5
		start := (page - 1) * limit
		end := start + limit
		if start > len(list) {
			start = len(list)
		}
		if end > len(list) {
			end = len(list)
		}
		json.NewEncoder(w).Encode(list[start:end])

	case http.MethodPost:
		var book models.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if book.Title == "" || book.Price <= 0 {
			http.Error(w, "Invalid book data", http.StatusBadRequest)
			return
		}
		book.ID = nextBookID
		nextBookID++
		books[book.ID] = book
		json.NewEncoder(w).Encode(book)
	}
}

func BookByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	book, exists := books[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(book)
	case http.MethodPut:
		var updated models.Book
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if updated.Title == "" || updated.Price <= 0 {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}
		updated.ID = id
		books[id] = updated
		json.NewEncoder(w).Encode(updated)
	case http.MethodDelete:
		delete(books, id)
		w.WriteHeader(http.StatusNoContent)
	}
}