package handlers

import (
	"encoding/json"
	"net/http"

	"bookstore/models"
)

var authors = make(map[int]models.Author)
var nextAuthorID = 1

func AuthorsInit(list []models.Author) {
	for _, a := range list {
		a.ID = nextAuthorID
		nextAuthorID++
		authors[a.ID] = a
	}
}
func AuthorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var list []models.Author
		for _, a := range authors {
			list = append(list, a)
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var author models.Author
		if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if author.Name == "" {
			http.Error(w, "Name required", http.StatusBadRequest)
			return
		}
		author.ID = nextAuthorID
		nextAuthorID++
		authors[author.ID] = author
		json.NewEncoder(w).Encode(author)
	}
}