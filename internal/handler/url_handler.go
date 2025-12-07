package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Faleeeee/URL_Shortener/internal/domain"
	"github.com/Faleeeee/URL_Shortener/internal/repository"
	"github.com/Faleeeee/URL_Shortener/internal/service"
	"github.com/Faleeeee/URL_Shortener/internal/utils"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service service.URLService
	baseURL string
}

// NewURLHandler creates a new URL handler
func NewURLHandler(service service.URLService, baseURL string) *URLHandler {
	return &URLHandler{
		service: service,
		baseURL: baseURL,
	}
}

// ShortenURL godoc
// @Summary Create a shortened URL
// @Description Create a short URL from a long URL with optional custom alias
// @Tags URL Shortener
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.ShortenRequest true "URL to shorten and optional alias"
// @Success 200 {object} domain.APIResponse{data=domain.ShortenResponse} "Successfully created short URL"
// @Failure 400 {object} domain.APIResponse "Invalid request or validation error"
// @Failure 401 {object} domain.APIResponse "Unauthorized"
// @Failure 409 {object} domain.APIResponse "Alias already exists"
// @Failure 500 {object} domain.APIResponse "Internal server error"
// @Router /url/shorten [post]
func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req domain.ShortenRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", "UNAUTHORIZED", "Missing or invalid authentication token")
		return
	}

	url, err := h.service.ShortenURL(req.URL, req.Alias, userID.(int64))
	if err != nil {
		if errors.Is(err, domain.ErrInvalidURL) ||
			errors.Is(err, domain.ErrURLTooLong) ||
			errors.Is(err, domain.ErrInvalidAlias) ||
			errors.Is(err, domain.ErrAliasTooLong) ||
			errors.Is(err, domain.ErrPrivateURL) {
			utils.SendError(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR", err.Error())
			return
		}

		// Duplicate alias error
		if errors.Is(err, repository.ErrDuplicateAlias) ||
			err.Error() == "alias '"+req.URL+"' is already taken" {
			utils.SendError(c, http.StatusConflict, "Alias already exists", "ALIAS_EXISTS", "The provided alias is already in use")
			return
		}

		// Internal server error
		utils.SendError(c, http.StatusInternalServerError, "Failed to create short URL", "INTERNAL_ERROR", "An unexpected error occurred")
		return
	}

	// Build response
	response := domain.ShortenResponse{
		Alias:       url.Alias,
		ShortURL:    h.baseURL + "/" + url.Alias,
		OriginalURL: url.OriginalURL,
	}

	utils.SendSuccess(c, "Short URL created successfully", response, nil)
}

// RedirectURL godoc
// @Summary Redirect to original URL
// @Description Redirect to the original URL using the short alias
// @Tags URL Shortener
// @Param alias path string true "Short URL alias"
// @Success 302 "Redirects to original URL"
// @Failure 404 {object} domain.APIResponse "Short URL not found"
// @Failure 500 {object} domain.APIResponse "Internal server error"
// @Router /{alias} [get]
func (h *URLHandler) RedirectURL(c *gin.Context) {
	alias := c.Param("alias")

	// Get URL by alias
	url, err := h.service.GetURLByAlias(alias)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.SendError(c, http.StatusNotFound, "Short URL not found", "URL_NOT_FOUND", "The requested alias does not exist")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve URL", "INTERNAL_ERROR", "An unexpected error occurred")
		return
	}

	// Increment click counter (fire and forget - don't block redirect)
	go h.service.IncrementClickCount(alias)

	// Redirect to original URL (302 Found - temporary redirect)
	c.Redirect(http.StatusFound, url.OriginalURL)
}

// GetURLInfo godoc
// @Summary Get URL information
// @Description Get detailed information about a shortened URL including click count (only owner can view)
// @Tags URL Shortener
// @Produce json
// @Security BearerAuth
// @Param alias path string true "Short URL alias"
// @Success 200 {object} domain.APIResponse{data=domain.URLInfoResponse} "URL information"
// @Failure 401 {object} domain.APIResponse "Unauthorized"
// @Failure 403 {object} domain.APIResponse "Forbidden - not owner"
// @Failure 404 {object} domain.APIResponse "Short URL not found"
// @Failure 500 {object} domain.APIResponse "Internal server error"
// @Router /url/links/{alias} [get]
func (h *URLHandler) GetURLInfo(c *gin.Context) {
	alias := c.Param("alias")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", "UNAUTHORIZED", "Missing or invalid authentication token")
		return
	}

	// Get URL by alias
	url, err := h.service.GetURLByAlias(alias)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.SendError(c, http.StatusNotFound, "Short URL not found", "URL_NOT_FOUND", "The requested alias does not exist")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve URL", "INTERNAL_ERROR", "An unexpected error occurred")
		return
	}

	// Check if user is the owner of this URL
	if url.UserID != userID.(int64) {
		utils.SendError(c, http.StatusForbidden, "You don't have permission to view this URL", "FORBIDDEN", "You are not the owner of this URL")
		return
	}

	// Build response
	response := domain.URLInfoResponse{
		Alias:       url.Alias,
		OriginalURL: url.OriginalURL,
		UserID:      url.UserID,
		ClickCount:  url.ClickCount,
		CreatedAt:   url.CreatedAt,
		UpdatedAt:   url.UpdatedAt,
	}

	utils.SendSuccess(c, "URL information retrieved successfully", response, nil)
}

// ListURLs godoc
// @Summary List all shortened URLs (Admin only)
// @Description Get a paginated list of all shortened URLs in the system
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of results to return" default(50)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} domain.APIResponse{data=[]domain.URL} "List of URLs with pagination metadata"
// @Failure 401 {object} domain.APIResponse "Unauthorized"
// @Failure 500 {object} domain.APIResponse "Internal server error"
// @Router /admin/url [get]
func (h *URLHandler) ListURLs(c *gin.Context) {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get URLs
	urls, err := h.service.ListURLs(limit, offset)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve URLs", "INTERNAL_ERROR", "An unexpected error occurred")
		return
	}

	// If no URLs found, return empty array
	if urls == nil {
		urls = []*domain.URL{}
	}

	// Build response with metadata
	meta := &domain.Meta{
		Page:  offset/limit + 1,
		Limit: limit,
		Total: int64(len(urls)), // Note: This is just the count of returned items, ideally we should have total count from DB
	}

	utils.SendSuccess(c, "URLs retrieved successfully", urls, meta)
}

// GetUserURLs godoc
// @Summary Get URLs created by authenticated user
// @Description Get a paginated list of all URLs created by the authenticated user
// @Tags URL Shortener
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of results to return" default(50)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} domain.APIResponse{data=[]domain.URL} "List of user's URLs with pagination metadata"
// @Failure 401 {object} domain.APIResponse "Unauthorized"
// @Failure 500 {object} domain.APIResponse "Internal server error"
// @Router /url/my-links [get]
func (h *URLHandler) GetUserURLs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", "UNAUTHORIZED", "Missing or invalid authentication token")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	urls, err := h.service.GetURLsByUserID(userID.(int64), limit, offset)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve URLs", "INTERNAL_ERROR", "An unexpected error occurred")
		return
	}

	if urls == nil {
		urls = []*domain.URL{}
	}

	meta := &domain.Meta{
		Page:  offset/limit + 1,
		Limit: limit,
		Total: int64(len(urls)), // Note: This is just the count of returned items, ideally we should have total count from DB
	}

	utils.SendSuccess(c, "User URLs retrieved successfully", urls, meta)
}
