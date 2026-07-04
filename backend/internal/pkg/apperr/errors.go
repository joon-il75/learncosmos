package apperr

import "net/http"

type AppError struct {
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Code + ": " + e.Message
}

func New(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(err error) *AppError {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*AppError); ok {
		return ae
	}
	return &AppError{Code: "INTERNAL", Message: err.Error()}
}

var (
	ErrUnauthorized       = New("UNAUTHORIZED", "authentication required")
	ErrForbidden          = New("FORBIDDEN", "insufficient permissions")
	ErrInvalidState       = New("INVALID_STATE", "invalid or expired oauth state")
	ErrInsufficientPoints = New("INSUFFICIENT_POINTS", "not enough AI points")
	ErrNotFound           = New("NOT_FOUND", "resource not found")
)

func HTTPStatus(err error) int {
	ae, ok := err.(*AppError)
	if !ok {
		return http.StatusInternalServerError
	}
	switch ae.Code {
	case "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "FORBIDDEN":
		return http.StatusForbidden
	case "INVALID_STATE":
		return http.StatusBadRequest
	case "NOT_FOUND":
		return http.StatusNotFound
	case "INSUFFICIENT_POINTS":
		return http.StatusPaymentRequired
	default:
		return http.StatusInternalServerError
	}
}
