package models

import "time"

// FavoriteBook represents a row in the favorite_books table.
type FavoriteBook struct {
	UserID    int       `json:"user_id"`
	BookID    int       `json:"book_id"`
	CreatedAt time.Time `json:"created_at"`
}

// FavoriteBookResponse is what the API returns: the full book enriched with
// the timestamp at which the user added it to favorites.
type FavoriteBookResponse struct {
	Book
	FavoritedAt time.Time `json:"favorited_at"`
}
