package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bookshelf/monolith/internal/domain"
	"github.com/bookshelf/monolith/internal/service"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListBookReviews(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid page")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid limit")
		return
	}

	reviewList, err := h.services.Review.ListByBookID(r.Context(), bookID, page, limit)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, reviewList)
}

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID := chi.URLParam(r, "reviewId")

	review, err := h.services.Review.GetByID(r.Context(), reviewID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrReviewNotFound):
			writeError(w, r, http.StatusNotFound, "USER_NOT_FOUND", "review not found")
		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, review)
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	bookID := chi.URLParam(r, "bookId")

	var req domain.CreateReviewRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid json body")
		return
	}

	var details []domain.ErrorDetail

	if req.Rating < 1 || req.Rating > 5 {
		details = append(details, domain.ErrorDetail{
			Field:   "rating",
			Message: "rating must be between 1 and 5",
		})
	}

	if len(req.Content) < 10 {
		details = append(details, domain.ErrorDetail{
			Field:   "content",
			Message: "content must be at least 10 characters",
		})
	}

	if len(details) > 0 {
		writeValidationError(w, r, details)
		return
	}

	review, err := h.services.Review.Create(r.Context(), userID, bookID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookNotFound):
			writeError(w, r, http.StatusNotFound, "BOOK_NOT_FOUND", "book not found")

		case errors.Is(err, service.ErrAlreadyReviewed):
			writeError(w, r, http.StatusConflict, "ALREADY_REVIEWED", "you have already reviewed this book")

		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, review)
}

func (h *Handler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	reviewID := chi.URLParam(r, "reviewId")

	var req domain.UpdateReviewRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid json body")
		return
	}

	review, err := h.services.Review.Update(r.Context(), userID, reviewID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrReviewNotFound):
			writeError(w, r, http.StatusNotFound, "REVIEW_NOT_FOUND", "review not found")

		case errors.Is(err, service.ErrNotReviewOwner):
			writeError(w, r, http.StatusForbidden, "FORBIDDEN", "you are not the owner of this review")

		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, review)
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	reviewID := chi.URLParam(r, "reviewId")

	if err := h.services.Review.Delete(r.Context(), userID, reviewID); err != nil {
		switch {
		case errors.Is(err, service.ErrReviewNotFound):
			writeError(w, r, http.StatusNotFound, "REVIEW_NOT_FOUND", "review not found")

		case errors.Is(err, service.ErrNotReviewOwner):
			writeError(w, r, http.StatusForbidden, "FORBIDDEN", "you are not the owner of this review")

		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
