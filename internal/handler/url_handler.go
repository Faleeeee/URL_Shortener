package handler

import (
	"Url-Shortener-Service/internal/domain"
	"Url-Shortener-Service/internal/repository"
	"Url-Shortener-Service/internal/service"
	"errors"
	"net/http"
	"strconv"

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
// @Success 200 {object} domain.ShortenResponse "Successfully created short URL"
// @Failure 400 {object} map[string]string "Invalid request or validation error"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 409 {object} map[string]string "Alias already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /url/shorten [post]
func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req domain.ShortenRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	url, err := h.service.ShortenURL(req.URL, req.Alias, userID.(int64))
	if err != nil {
		if errors.Is(err, domain.ErrInvalidURL) ||
			errors.Is(err, domain.ErrURLTooLong) ||
			errors.Is(err, domain.ErrInvalidAlias) ||
			errors.Is(err, domain.ErrAliasTooLong) ||
			errors.Is(err, domain.ErrPrivateURL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Duplicate alias error
		if errors.Is(err, repository.ErrDuplicateAlias) ||
			err.Error() == "alias '"+req.URL+"' is already taken" {
			c.JSON(http.StatusConflict, gin.H{"error": "Alias already exists"})
			return
		}

		// Internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	// Build response
	response := domain.ShortenResponse{
		Alias:       url.Alias,
		ShortURL:    h.baseURL + "/" + url.Alias,
		OriginalURL: url.OriginalURL,
	}

	c.JSON(http.StatusOK, response)
}

// RedirectURL godoc
// @Summary Redirect to original URL
// @Description Redirect to the original URL using the short alias
// @Tags URL Shortener
// @Param alias path string true "Short URL alias"
// @Success 302 "Redirects to original URL"
// @Failure 404 {object} map[string]string "Short URL not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /{alias} [get]
func (h *URLHandler) RedirectURL(c *gin.Context) {
	alias := c.Param("alias")

	// Get URL by alias
	url, err := h.service.GetURLByAlias(alias)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
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
// @Success 200 {object} domain.URLInfoResponse "URL information"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - not owner"
// @Failure 404 {object} map[string]string "Short URL not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /url/links/{alias} [get]
func (h *URLHandler) GetURLInfo(c *gin.Context) {
	alias := c.Param("alias")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get URL by alias
	url, err := h.service.GetURLByAlias(alias)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		return
	}

	// Check if user is the owner of this URL
	if url.UserID != userID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this URL"})
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

	c.JSON(http.StatusOK, response)
}

// ListURLs godoc
// @Summary List all shortened URLs (Admin only)
// @Description Get a paginated list of all shortened URLs in the system
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of results to return" default(50)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} map[string]interface{} "List of URLs with pagination metadata"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/url [get]
func (h *URLHandler) ListURLs(c *gin.Context) {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get URLs
	urls, err := h.service.ListURLs(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}

	// If no URLs found, return empty array
	if urls == nil {
		urls = []*domain.URL{}
	}

	// Build response with metadata
	c.JSON(http.StatusOK, gin.H{
		"urls":   urls,
		"count":  len(urls),
		"limit":  limit,
		"offset": offset,
	})
}

// GetUserURLs godoc
// @Summary Get URLs created by authenticated user
// @Description Get a paginated list of all URLs created by the authenticated user
// @Tags URL Shortener
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of results to return" default(50)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} map[string]interface{} "List of user's URLs with pagination metadata"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /url/my-links [get]
func (h *URLHandler) GetUserURLs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	urls, err := h.service.GetURLsByUserID(userID.(int64), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}

	if urls == nil {
		urls = []*domain.URL{}
	}

	c.JSON(http.StatusOK, gin.H{
		"urls":   urls,
		"count":  len(urls),
		"limit":  limit,
		"offset": offset,
	})
}
