package handler

import (
	"Url-Shortener-Service/internal/domain"
	"Url-Shortener-Service/internal/repository"
	"Url-Shortener-Service/internal/service"
	"Url-Shortener-Service/internal/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with username and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Registration request with username and password"
// @Success 201 {object} domain.APIResponse{data=domain.AuthResponse} "Successfully registered, returns JWT token"
// @Failure 400 {object} domain.APIResponse "Invalid request or validation error"
// @Failure 409 {object} domain.APIResponse "Username already exists"
// @Failure 500 {object} domain.APIResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	response, err := h.authService.Register(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUsername) ||
			errors.Is(err, domain.ErrInvalidPassword) {
			utils.SendError(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR", err.Error())
			return
		}

		if errors.Is(err, domain.ErrUsernameAlreadyExists) ||
			errors.Is(err, repository.ErrDuplicateUsername) {
			utils.SendError(c, http.StatusConflict, "Username already exists", "USERNAME_EXISTS", "The provided username is already taken")
			return
		}

		utils.SendError(c, http.StatusInternalServerError, "Failed to register user", "INTERNAL_ERROR", "An unexpected error occurred")
		return
	}

	utils.SendSuccess(c, "User registered successfully", response, nil)
}

// Login godoc
// @Summary User login
// @Description Authenticate user with username and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login request with username and password"
// @Success 200 {object} domain.APIResponse{data=domain.AuthResponse} "Successfully authenticated, returns JWT token"
// @Failure 400 {object} domain.APIResponse "Invalid request body"
// @Failure 401 {object} domain.APIResponse "Invalid username or password"
// @Failure 500 {object} domain.APIResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}
	response, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) ||
			errors.Is(err, domain.ErrUserNotFound) {
			utils.SendError(c, http.StatusUnauthorized, "Invalid username or password", "AUTH_FAILED", "Invalid credentials provided")
			return
		}

		utils.SendError(c, http.StatusInternalServerError, "Failed to login", "INTERNAL_ERROR", "An unexpected error occurred")
		return
	}

	utils.SendSuccess(c, "Login successful", response, nil)
}
