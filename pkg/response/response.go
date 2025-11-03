package response

import (
	"github.com/gin-gonic/gin"
)

// SuccessResponse represents a standard success API response.
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents a standard error API response.
type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"error_code,omitempty"`
}

// SendSuccess sends a standardized success response.
func SendSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// SendError sends a standardized error response and aborts the request.
func SendError(c *gin.Context, status int, message string, errorCode string) {
	c.AbortWithStatusJSON(status, ErrorResponse{
		Message:   message,
		ErrorCode: errorCode,
	})
}
