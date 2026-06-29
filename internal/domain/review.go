package domain

import (
	"database/sql"
	"time"
)

// ========== entities ==========

// Review основная доменная модель рецензии для хранения в БД
type Review struct {
	ID      string         `json:"id" db:"id"`
	BookID  string         `json:"book_id" db:"book_id"`
	UserID  string         `json:"user_id" db:"user_id"`
	Rating  int            `json:"rating" db:"rating"`
	Title   sql.NullString `json:"title" db:"title"`
	Content string         `json:"content" db:"content"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ReviewResponse модель рецензии для API ответа
type ReviewResponse struct {
	ID      string `json:"id"`
	BookID  string `json:"book_id"`
	UserID  string `json:"user_id"`
	Rating  int    `json:"rating"`
	Title   *string `json:"title"`
	Content string `json:"content"`

	User UserSummary `json:"user"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ========== DTO ==========

// CreateReviewRequest данные для создания рецензии
type CreateReviewRequest struct {
	Rating  int     `json:"rating"`
	Title   *string `json:"title"`
	Content string  `json:"content"`
}

// UpdateReviewRequest данные для обновления рецензии (все поля опциональные)
type UpdateReviewRequest struct {
	Rating  *int    `json:"rating"`
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

// ReviewListResponse ответ со списком рецензий и пагинацией
type ReviewListResponse struct {
	Data       []ReviewResponse `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

// ========== methods ==========

// ToResponse конвертирует Review (DB модель) в ReviewResponse (API модель)
func (r *Review) ToResponse(user *User) ReviewResponse {
	response := ReviewResponse{
		ID:        r.ID,
		BookID:    r.BookID,
		UserID:    r.UserID,
		Rating:    r.Rating,
		Content:   r.Content,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	if r.Title.Valid {
		title := r.Title.String
		response.Title = &title
	}

	if user != nil {
		response.User = user.ToSummary()
	}

	return response
}