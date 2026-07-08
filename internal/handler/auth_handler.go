package handler

import (
	"errors"
	"net/http"

	"github.com/bookshelf/monolith/internal/domain"
	"github.com/bookshelf/monolith/internal/service"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid json body")
		return
	}

	var details []domain.ErrorDetail

	if req.Username == "" {
		details = append(details, domain.ErrorDetail{
			Field:   "username",
			Message: "username is required",
		})
	}

	if req.Email == "" {
		details = append(details, domain.ErrorDetail{
			Field:   "email",
			Message: "username is required",
		})
	}

	if req.Password == "" {
		details = append(details, domain.ErrorDetail{
			Field:   "password",
			Message: "password is required",
		})
	}

	if len(details) > 0 {
		writeValidationError(w, r, details)
		return
	}

	response, err := h.services.User.Register(r.Context(), req)
	if err != nil {
		switch err {
		case service.ErrUserExists:
			writeError(w, r, http.StatusConflict, "USER_EXISTS", "user already exists")
		case service.ErrUsernameExists:
			writeError(w, r, http.StatusConflict, "USERNAME_EXISTS", "username already exists")
		case service.ErrInvalidUsername, service.ErrInvalidPassword, service.ErrInvalidEmail:
			writeError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid json body")
		return
	}

	var details []domain.ErrorDetail
	if req.Email == "" {
		details = append(details, domain.ErrorDetail{
			Field:   "email",
			Message: "username is required",
		})
	}

	if req.Password == "" {
		details = append(details, domain.ErrorDetail{
			Field:   "password",
			Message: "password is required",
		})
	}

	if len(details) > 0 {
		writeValidationError(w, r, details)
		return
	}

	resp, err := h.services.User.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			writeError(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())

	user, err := h.services.User.GetByID(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			writeError(w, r, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, user.ToPublic())
}

func (h *Handler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())

	var req domain.UpdateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid json body")
		return
	}

	user, err := h.services.User.Update(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			writeError(w, r, http.StatusNotFound, "USER_NOT_FOUND", "user not found")

		case errors.Is(err, service.ErrUsernameExists):
			writeError(w, r, http.StatusConflict, "USERNAME_EXISTS", "username already exists")

		case errors.Is(err, service.ErrInvalidUsername):
			writeError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())

		default:
			writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, user.ToPublic())
}
