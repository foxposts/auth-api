package context

import "net/http"

type ApiError struct {
	Code    int    `json:",omitempty"`
	Message string `json:"message"`
}

func (e ApiError) AsMessage() *ApiError {
	return &ApiError{
		Message: e.Message,
	}
}

func AuthenticationError(message string) *ApiError {
	return &ApiError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func UnexpectedError(message string) *ApiError {
	return &ApiError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func NotFoundError(message string) *ApiError {
	return &ApiError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func StatusError(message string) *ApiError {
	return &ApiError{
		Message: message,
		Code:    http.StatusForbidden,
	}
}
