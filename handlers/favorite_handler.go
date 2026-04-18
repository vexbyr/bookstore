package handlers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"bookstore_gin/middleware"
	"bookstore_gin/models"

	"github.com/gin-gonic/gin"
)

// favoriteKey uniquely identifies a row in the favorite_books table.
type favoriteKey struct {
	UserID int
	BookID int
}

// FavoriteHandler owns the in-memory favorite_books table and all three
// favorite endpoints. It holds a reference to BookHandler so it can look up
// full book details when returning the list.
type FavoriteHandler struct {
	mu        sync.RWMutex
	favorites map[favoriteKey]models.FavoriteBook // keyed by (user_id, book_id)
	books     *BookHandler
}

func NewFavoriteHandler(books *BookHandler) *FavoriteHandler {
	return &FavoriteHandler{
		favorites: make(map[favoriteKey]models.FavoriteBook),
		books:     books,
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// GET /books/favorites
// Returns a paginated list of the authenticated user's favorite books.
//
// Query params:
//
//	page      int (default 1)
//	page_size int (default 10, max 100)
//
// ──────────────────────────────────────────────────────────────────────────────
func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	// Collect all favorites that belong to this user.
	result := make([]models.FavoriteBookResponse, 0)
	for key, fav := range h.favorites {
		if key.UserID != userID {
			continue
		}

		// Enrich with full book data.
		book, ok := h.books.getByID(fav.BookID)
		if !ok {
			// Book was deleted after being favorited — skip it silently.
			continue
		}

		result = append(result, models.FavoriteBookResponse{
			Book:        book,
			FavoritedAt: fav.CreatedAt,
		})
	}

	total := len(result)

	// Paginate.
	start := (page - 1) * pageSize
	if start >= total {
		c.JSON(http.StatusOK, gin.H{
			"data":      []models.FavoriteBookResponse{},
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		})
		return
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      result[start:end],
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ──────────────────────────────────────────────────────────────────────────────
// PUT /books/:bookId/favorites
// Adds a book to the authenticated user's favorites.
// Idempotent: adding the same book twice returns 200 with a clear message.
// ──────────────────────────────────────────────────────────────────────────────
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID := middleware.GetUserID(c)

	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil || bookID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookId"})
		return
	}

	// Verify the book exists before favoriting it.
	if _, ok := h.books.getByID(bookID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	key := favoriteKey{UserID: userID, BookID: bookID}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Idempotency check.
	if _, exists := h.favorites[key]; exists {
		c.JSON(http.StatusOK, gin.H{"message": "book is already in favorites", "book_id": bookID})
		return
	}

	// Insert new row into the favorite_books table.
	h.favorites[key] = models.FavoriteBook{
		UserID:    userID,
		BookID:    bookID,
		CreatedAt: time.Now().UTC(),
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "book added to favorites",
		"user_id":    userID,
		"book_id":    bookID,
		"created_at": h.favorites[key].CreatedAt,
	})
}

// ──────────────────────────────────────────────────────────────────────────────
// DELETE /books/:bookId/favorites
// Removes a book from the authenticated user's favorites.
// ──────────────────────────────────────────────────────────────────────────────
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID := middleware.GetUserID(c)

	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil || bookID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookId"})
		return
	}

	key := favoriteKey{UserID: userID, BookID: bookID}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.favorites[key]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "book is not in favorites"})
		return
	}

	delete(h.favorites, key)
	c.JSON(http.StatusOK, gin.H{"message": "book removed from favorites", "book_id": bookID})
}
