package handlers

import (
	"encoding/json"
	"net/http"

	"bookstore/models"
)

var categories = make(map[int]models.Category)
var nextCategoryID = 1

func CategoriesInit(list []models.Category) {
	for _, c := range list {
		c.ID = nextCategoryID
		nextCategoryID++
		categories[c.ID] = c
	}
}
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var list []models.Category
		for _, c := range categories {
			list = append(list, c)
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var category models.Category
		if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if category.Name == "" {
			http.Error(w, "Name required", http.StatusBadRequest)
			return
		}
		category.ID = nextCategoryID
		nextCategoryID++
		categories[category.ID] = category
		json.NewEncoder(w).Encode(category)
	}
}