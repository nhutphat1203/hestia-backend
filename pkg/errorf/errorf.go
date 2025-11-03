package errorf

import "net/http"

const (
	// Standard errors
	Default      = "error.default"
	InvalidJSON  = "error.invalid_json"
	NotFound     = "error.not_found"
	Unauthorized = "error.unauthorized"
	Forbidden    = "error.forbidden"

	// Auth errors
	InvalidToken = "error.invalid_token"
	ExpiredToken = "error.expired_token"

	// Validation errors
	Validation = "error.validation"

	// Internal errors
	Internal = "error.internal"
)

var message = map[string]string{
	Default:      "An unexpected error occurred. Please try again later.",
	InvalidJSON:  "Invalid JSON format.",
	NotFound:     "The requested resource was not found.",
	Unauthorized: "Authentication is required and has failed or has not yet been provided.",
	Forbidden:    "You do not have permission to access this resource.",
	InvalidToken: "The provided token is invalid.",
	ExpiredToken: "The provided token has expired.",
	Validation:   "The given data was invalid.",
	Internal:     "An internal server error occurred.",
}

var httpStatus = map[string]int{
	Default:      http.StatusInternalServerError,
	InvalidJSON:  http.StatusBadRequest,
	NotFound:     http.StatusNotFound,
	Unauthorized: http.StatusUnauthorized,
	Forbidden:    http.StatusForbidden,
	InvalidToken: http.StatusUnauthorized,
	ExpiredToken: http.StatusUnauthorized,
	Validation:   http.StatusBadRequest,
	Internal:     http.StatusInternalServerError,
}

// Message returns the error message for a given error code.
func Message(code string) string {
	if msg, ok := message[code]; ok {
		return msg
	}
	return message[Default]
}

// HttpStatus returns the HTTP status code for a given error code.
func HttpStatus(code string) int {
	if status, ok := httpStatus[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}
