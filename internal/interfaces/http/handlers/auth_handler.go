package handlers

import (
	"net/http"
	"strconv"

	service "github.com/nhutphat1203/hestia-backend/internal/services"
	"github.com/nhutphat1203/hestia-backend/pkg/errorf"
	"github.com/nhutphat1203/hestia-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// LoginRequest represents the request body for the login endpoint.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

// Login handles user login and token generation.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, errorf.HttpStatus(errorf.InvalidJSON), errorf.Message(errorf.InvalidJSON), errorf.InvalidJSON)
		return
	}

	if req.Username == "" {
		response.SendError(c, errorf.HttpStatus(errorf.Validation), "Username is required", errorf.Validation)
		return
	}

	tokens, err := h.authService.Login(req.Username, req.Password)

	if err != nil {
		response.SendError(c, http.StatusUnauthorized, "Invalid username or password", err.Error())
		return
	}

	response.SendSuccess(c, http.StatusOK, "Login successful", tokens)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var body LogoutRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.SendError(c, errorf.HttpStatus(errorf.InvalidJSON), errorf.Message(errorf.InvalidJSON), errorf.InvalidJSON)
		return
	}

	if body.RefreshToken == "" {
		response.SendError(c, errorf.HttpStatus(errorf.Validation), "Refresh token is required", errorf.Validation)
		return
	}

	err := h.authService.Logout(body.RefreshToken)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "Failed to logout", err.Error())
		return
	}

	response.SendSuccess(c, http.StatusOK, "Logout successful", nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var body RefreshTokenRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.SendError(c, errorf.HttpStatus(errorf.InvalidJSON), errorf.Message(errorf.InvalidJSON), errorf.InvalidJSON)
		return
	}

	if body.UserID == "" || body.RefreshToken == "" {
		response.SendError(c, errorf.HttpStatus(errorf.Validation), "User ID and refresh token are required", errorf.Validation)
		return
	}
	atoi, err := strconv.Atoi(body.UserID)
	if err != nil {
		response.SendError(c, errorf.HttpStatus(errorf.Validation), "User ID must be a valid integer", errorf.Validation)
		return
	}

	tokens, err := h.authService.RefreshToken(uint(atoi), body.RefreshToken)
	if err != nil {
		response.SendError(c, http.StatusUnauthorized, "Invalid refresh token", err.Error())
		return
	}

	response.SendSuccess(c, http.StatusOK, "Token refreshed successfully", tokens)
}
