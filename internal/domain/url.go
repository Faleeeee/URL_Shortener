package domain

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

// URL represents a shortened URL entity
type URL struct {
	ID          int64     `json:"id"`
	Alias       string    `json:"alias"`
	OriginalURL string    `json:"original_url"`
	ClickCount  int64     `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ShortenRequest represents the request to create a short URL
type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

// ShortenResponse represents the response after creating a short URL
type ShortenResponse struct {
	Alias       string `json:"alias"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// URLInfoResponse represents detailed URL information
type URLInfoResponse struct {
	Alias       string    `json:"alias"`
	OriginalURL string    `json:"original_url"`
	ClickCount  int64     `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validation errors
var (
	ErrInvalidURL   = errors.New("invalid URL format")
	ErrURLTooLong   = errors.New("URL exceeds maximum length of 2048 characters")
	ErrInvalidAlias = errors.New("alias must contain only alphanumeric characters, hyphens, and underscores")
	ErrAliasTooLong = errors.New("alias must not exceed 16 characters")
	ErrPrivateURL   = errors.New("private IP addresses and localhost are not allowed")
)

const (
	MaxURLLength   = 2048
	MaxAliasLength = 16
)

// ValidateURL validates the original URL
func ValidateURL(rawURL string) error {
	// Check length
	if len(rawURL) > MaxURLLength {
		return ErrURLTooLong
	}

	// Parse URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ErrInvalidURL
	}

	// Check scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrInvalidURL
	}

	// Check if host is present
	if parsedURL.Host == "" {
		return ErrInvalidURL
	}

	// Prevent localhost and private IPs
	host := strings.ToLower(parsedURL.Host)
	if strings.Contains(host, "localhost") ||
		strings.HasPrefix(host, "127.") ||
		strings.HasPrefix(host, "192.168.") ||
		strings.HasPrefix(host, "10.") ||
		strings.HasPrefix(host, "172.16.") {
		return ErrPrivateURL
	}

	return nil
}

// ValidateAlias validates a custom alias
func ValidateAlias(alias string) error {
	if alias == "" {
		return nil // Empty alias is okay, will be auto-generated
	}

	// Check length
	if len(alias) > MaxAliasLength {
		return ErrAliasTooLong
	}

	// Check characters (alphanumeric, hyphen, underscore only)
	for _, char := range alias {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return ErrInvalidAlias
		}
	}

	return nil
}
