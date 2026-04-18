package handlers

import (
	"net/http"
	"strings"
	"sync"

	"bookstore_gin/models"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	mu         sync.RWMutex
	categories map[int]models.Category
	nextID     int
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categories: map[int]models.Category{
			1: {ID: 1, Name: "Fiction"},
			2: {ID: 2, Name: "Science Fiction"},
			3: {ID: 3, Name: "Fantasy"},
			4: {ID: 4, Name: "Non-Fiction"},
		},
		nextID: 5,
	}
}

// GetAll returns all categories.
func (h *CategoryHandler) GetAll(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	list := make([]models.Category, 0, len(h.categories))
	for _, cat := range h.categories {
		list = append(list, cat)
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "total": len(list)})
}

// Create adds a new category.
func (h *CategoryHandler) Create(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := strings.TrimSpace(input.Name)
	if len(name) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name must be at least 2 characters"})
		return
	}

	// Check for duplicate category name
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, cat := range h.categories {
		if strings.EqualFold(cat.Name, name) {
			c.JSON(http.StatusConflict, gin.H{"error": "category already exists"})
			return
		}
	}

	cat := models.Category{ID: h.nextID, Name: name}
	h.categories[h.nextID] = cat
	h.nextID++

	c.JSON(http.StatusCreated, cat)
}

// GetByID returns a category by ID (used by book handler for validation).
func (h *CategoryHandler) GetByID(id int) (models.Category, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	cat, ok := h.categories[id]
	return cat, ok
}

// GetByName looks up a category by name (case-insensitive).
func (h *CategoryHandler) GetByName(name string) (models.Category, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, cat := range h.categories {
		if strings.EqualFold(cat.Name, name) {
			return cat, true
		}
	}
	return models.Category{}, false
}
