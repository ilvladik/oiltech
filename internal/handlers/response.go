package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"oiltech/internal/domain"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func WriteAPIError(w http.ResponseWriter, status int, code string, message string) {
	WriteJSON(w, status, ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}

func WriteError(w http.ResponseWriter, err error) {
	var domainError *domain.DomainError
	if errors.As(err, &domainError) {
		status := http.StatusBadRequest
		switch domainError.Code() {
		case domain.ErrNotFoundCode:
			status = http.StatusNotFound
		case domain.ErrAlreadyExistsCode:
			status = http.StatusConflict
		}
		WriteAPIError(w, status, string(domainError.Code()), domainError.Message())
		return
	}
	WriteAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
}
