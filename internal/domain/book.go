package domain

import (
	"database/sql"
	"time"
)

type Book struct {
	ID        string `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Author    string `json:"author" db:"author"`
	CreatedBy string `json:"created_by" db:"created_by"`

	Description   sql.NullString  `json:"description" db:"description"`
	ISBN          sql.NullString  `json:"isbn" db:"isbn"`
	PublishedYear sql.NullInt32   `json:"published_year" db:"published_year"`
	AverageRating sql.NullFloat64 `json:"average_rating" db:"average_rating"`

	ReviewsCount int `json:"reviews_count" db:"reviews_count"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type BookResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	CreatedBy string `json:"created_by"`

	Description   *string  `json:"description"`
	ISBN          *string  `json:"isbn"`
	PublishedYear *int     `json:"published_year"`
	AverageRating *float64 `json:"average_rating"`

	ReviewsCount int `json:"reviews_count"`

	Creator *UserSummary `json:"creator,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request типы
type CreateBookRequest struct {
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	Description   *string `json:"description"`
	ISBN          *string `json:"isbn"`
	PublishedYear *int    `json:"published_year"`
}

type UpdateBookRequest struct {
	Title         *string `json:"title"`
	Author        *string `json:"author"`
	Description   *string `json:"description"`
	ISBN          *string `json:"isbn"`
	PublishedYear *int    `json:"published_year"`
}

type BookFilter struct {
	Search string
	Sort   string
	Order  string
	Page   int
	Limit  int
}

type BookListResponse struct {
	Data       []BookResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

func (b *Book) ToResponse() BookResponse {
	response := BookResponse{
		ID:           b.ID,
		Title:        b.Title,
		Author:       b.Author,
		CreatedBy:    b.CreatedBy,
		ReviewsCount: b.ReviewsCount,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}

	if b.Description.Valid {
		description := b.Description.String
		response.Description = &description
	}

	if b.ISBN.Valid {
		isbn := b.ISBN.String
		response.ISBN = &isbn
	}

	if b.PublishedYear.Valid {
		year := int(b.PublishedYear.Int32)
		response.PublishedYear = &year
	}

	if b.AverageRating.Valid {
		rating := b.AverageRating.Float64
		response.AverageRating = &rating
	}

	return response
}

