package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bookshelf/monolith/internal/domain"
	"github.com/bookshelf/monolith/internal/service"
)

type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
)

type Handler struct {
	services  *service.Service
	jwtSecret string
}

func New(services *service.Service, jwtSecret string) *Handler {
	return &Handler{
		services:  services,
		jwtSecret: jwtSecret,
	}
}

// helper functions

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	requestID, _ := r.Context().Value(requestIDKey).(string)

	response := domain.ErrorResponse{
		Code:      code,
		Message:   message,
		RequestID: requestID,
	}

	writeJSON(w, status, response)
}

func writeValidationError(w http.ResponseWriter, r *http.Request, details []domain.ErrorDetail) {
	requestID, _ := r.Context().Value(requestIDKey).(string)

	response := domain.ErrorResponse{
		Code:      "VALIDATION_ERROR",
		Message:   "validation failed",
		Details:   details,
		RequestID: requestID,
	}

	writeJSON(w, http.StatusBadRequest, response)
}

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func getUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC(),
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC(),
		"checks": map[string]string{
			"database": "ok",
		},
	}

	writeJSON(w, http.StatusOK, response)
}
