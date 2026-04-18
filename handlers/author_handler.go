package handlers

import (
	"net/http"
	"strings"
	"sync"

	"bookstore_gin/models"

	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	mu      sync.RWMutex
	authors map[int]models.Author
	nextID  int
}

func NewAuthorHandler() *AuthorHandler {
	return &AuthorHandler{
		authors: map[int]models.Author{
			1: {ID: 1, Name: "George Orwell"},
			2: {ID: 2, Name: "J.K. Rowling"},
			3: {ID: 3, Name: "Frank Herbert"},
		},
		nextID: 4,
	}
}

// GetAll returns all authors.
func (h *AuthorHandler) GetAll(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	list := make([]models.Author, 0, len(h.authors))
	for _, a := range h.authors {
		list = append(list, a)
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "total": len(list)})
}

// Create adds a new author.
func (h *AuthorHandler) Create(c *gin.Context) {
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

	h.mu.Lock()
	defer h.mu.Unlock()

	author := models.Author{ID: h.nextID, Name: name}
	h.authors[h.nextID] = author
	h.nextID++

	c.JSON(http.StatusCreated, author)
}

// GetByID returns an author by ID (used by book handler for validation).
func (h *AuthorHandler) GetByID(id int) (models.Author, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	a, ok := h.authors[id]
	return a, ok
}
