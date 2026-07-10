package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bookshelf/monolith/internal/domain"
	"github.com/bookshelf/monolith/internal/service"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10

	var err error

	if p := r.URL.Query().Get("page"); p != "" {
		page, err = strconv.Atoi(p)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid page")
			return
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		limit, err = strconv.Atoi(l)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid limit")
			return
		}
	}

	filter := domain.BookFilter{
		Search: r.URL.Query().Get("search"),
		Sort:   r.URL.Query().Get("sort"),
		Order:  r.URL.Query().Get("order"),
		Page:   page,
		Limit:  limit,
	}

	resp, err := h.services.Book.List(r.Context(), filter)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")

	bookResponse, err := h.services.Book.GetByID(r.Context(), bookID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookNotFound):
			writeError(w, r, http.StatusNotFound, "BOOK_NOT_FOUND", "book not found")
		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, bookResponse)
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())

	var req domain.CreateBookRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid json body")
		return
	}

	var details []domain.ErrorDetail

	if req.Title == "" {
		details = append(details, domain.ErrorDetail{
			Field:   "title",
			Message: "title is required",
		})
	}

	if req.Author == "" {
		details = append(details, domain.ErrorDetail{
			Field:   "author",
			Message: "author is required",
		})
	}

	if len(details) > 0 {
		writeValidationError(w, r, details)
		return
	}

	book, err := h.services.Book.Create(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookTitleEmpty),
			errors.Is(err, service.ErrBookAuthorEmpty):
			writeError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())

		case errors.Is(err, service.ErrUserNotFound):
			writeError(w, r, http.StatusNotFound, "USER_NOT_FOUND", "user not found")

		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, book)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	bookID := chi.URLParam(r, "bookId")

	var req domain.UpdateBookRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid json body")
		return
	}

	book, err := h.services.Book.Update(r.Context(), userID, bookID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookNotFound):
			writeError(w, r, http.StatusNotFound, "BOOK_NOT_FOUND", "book not found")

		case errors.Is(err, service.ErrNotBookOwner):
			writeError(w, r, http.StatusForbidden, "FORBIDDEN", "you are not the owner of this book")

		case errors.Is(err, service.ErrBookTitleEmpty),
			errors.Is(err, service.ErrBookAuthorEmpty):
			writeError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())

		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, book)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	bookID := chi.URLParam(r, "bookId")

	if err := h.services.Book.Delete(r.Context(), userID, bookID); err != nil {
		switch {
		case errors.Is(err, service.ErrBookNotFound):
			writeError(w, r, http.StatusNotFound, "BOOK_NOT_FOUND", "book not found")

		case errors.Is(err, service.ErrNotBookOwner):
			writeError(w, r, http.StatusForbidden, "FORBIDDEN", "you are not the owner of this book")

		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
