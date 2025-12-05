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

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req domain.ShortenRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create shortened URL
	url, err := h.service.ShortenURL(req.URL, req.Alias)
	if err != nil {
		// Check for specific errors
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

func (h *URLHandler) GetURLInfo(c *gin.Context) {
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

	// Build response
	response := domain.URLInfoResponse{
		Alias:       url.Alias,
		OriginalURL: url.OriginalURL,
		ClickCount:  url.ClickCount,
		CreatedAt:   url.CreatedAt,
		UpdatedAt:   url.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

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
