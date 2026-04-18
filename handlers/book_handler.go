package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"bookstore_gin/models"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	mu         sync.RWMutex
	books      map[int]models.Book
	nextID     int
	authors    *AuthorHandler
	categories *CategoryHandler
}

func NewBookHandler(authors *AuthorHandler, categories *CategoryHandler) *BookHandler {
	return &BookHandler{
		books: map[int]models.Book{
			1: {ID: 1, Title: "1984", AuthorID: 1, CategoryID: 1, Price: 9.99},
			2: {ID: 2, Title: "Harry Potter and the Sorcerer's Stone", AuthorID: 2, CategoryID: 3, Price: 14.99},
			3: {ID: 3, Title: "Dune", AuthorID: 3, CategoryID: 2, Price: 12.99},
			4: {ID: 4, Title: "Animal Farm", AuthorID: 1, CategoryID: 1, Price: 7.99},
		},
		nextID:     5,
		authors:    authors,
		categories: categories,
	}
}

// bookInput is the validated request body for create/update.
type bookInput struct {
	Title      string  `json:"title"      binding:"required"`
	AuthorID   int     `json:"author_id"  binding:"required,min=1"`
	CategoryID int     `json:"category_id" binding:"required,min=1"`
	Price      float64 `json:"price"      binding:"required,min=0.01"`
}

// validate checks cross-field business rules after binding.
func (h *BookHandler) validate(input bookInput) string {
	if strings.TrimSpace(input.Title) == "" {
		return "title cannot be blank"
	}
	if _, ok := h.authors.GetByID(input.AuthorID); !ok {
		return "author_id does not exist"
	}
	if _, ok := h.categories.GetByID(input.CategoryID); !ok {
		return "category_id does not exist"
	}
	return ""
}

// GetAll handles GET /books with pagination and optional filters.
//
//	Query params:
//	  page      int    (default 1)
//	  page_size int    (default 10, max 100)
//	  category  string (filter by category name, case-insensitive)
//	  author_id int    (filter by author ID)
//	  min_price float  (filter books with price >= min_price)
//	  max_price float  (filter books with price <= max_price)
func (h *BookHandler) GetAll(c *gin.Context) {
	// --- parse query params ---
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	categoryFilter := strings.TrimSpace(c.Query("category"))
	authorIDFilter, _ := strconv.Atoi(c.Query("author_id"))
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	// Resolve category name → ID
	var filterCategoryID int
	if categoryFilter != "" {
		cat, ok := h.categories.GetByName(categoryFilter)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"data": []models.Book{}, "total": 0, "page": page, "page_size": pageSize})
			return
		}
		filterCategoryID = cat.ID
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	// --- filter ---
	filtered := make([]models.Book, 0)
	for _, b := range h.books {
		if filterCategoryID != 0 && b.CategoryID != filterCategoryID {
			continue
		}
		if authorIDFilter != 0 && b.AuthorID != authorIDFilter {
			continue
		}
		if minPrice > 0 && b.Price < minPrice {
			continue
		}
		if maxPrice > 0 && b.Price > maxPrice {
			continue
		}
		filtered = append(filtered, b)
	}

	total := len(filtered)

	// --- paginate ---
	start := (page - 1) * pageSize
	if start >= total {
		c.JSON(http.StatusOK, gin.H{"data": []models.Book{}, "total": total, "page": page, "page_size": pageSize})
		return
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      filtered[start:end],
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetByID handles GET /books/:id.
func (h *BookHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	book, ok := h.books[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// Create handles POST /books.
func (h *BookHandler) Create(c *gin.Context) {
	var input bookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if msg := h.validate(input); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	book := models.Book{
		ID:         h.nextID,
		Title:      strings.TrimSpace(input.Title),
		AuthorID:   input.AuthorID,
		CategoryID: input.CategoryID,
		Price:      input.Price,
	}
	h.books[h.nextID] = book
	h.nextID++

	c.JSON(http.StatusCreated, book)
}

// Update handles PUT /books/:id.
func (h *BookHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input bookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if msg := h.validate(input); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.books[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	updated := models.Book{
		ID:         id,
		Title:      strings.TrimSpace(input.Title),
		AuthorID:   input.AuthorID,
		CategoryID: input.CategoryID,
		Price:      input.Price,
	}
	h.books[id] = updated
	c.JSON(http.StatusOK, updated)
}

// Delete handles DELETE /books/:id.
func (h *BookHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.books[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	delete(h.books, id)
	c.JSON(http.StatusOK, gin.H{"message": "book deleted", "id": id})
}

// getByID is an internal lock-safe helper for FavoriteHandler.
// Callers must NOT already hold h.mu.
func (h *BookHandler) getByID(id int) (models.Book, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	b, ok := h.books[id]
	return b, ok
}
