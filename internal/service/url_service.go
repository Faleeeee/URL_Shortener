package service

import (
	"Url-Shortener-Service/internal/domain"
	"Url-Shortener-Service/internal/repository"
	"errors"
	"fmt"
)

const (
	MaxRetries = 3 // Maximum retries for collision resolution
)

var (
	ErrMaxRetriesExceeded = errors.New("maximum retries exceeded for generating unique alias")
)

// URLService defines the interface for URL shortening business logic
type URLService interface {
	ShortenURL(originalURL string, alias string, userID int64) (*domain.URL, error)
	GetURLByAlias(alias string) (*domain.URL, error)
	IncrementClickCount(alias string) error
	ListURLs(limit, offset int) ([]*domain.URL, error)
	GetURLsByUserID(userID int64, limit, offset int) ([]*domain.URL, error)
}

type urlService struct {
	repo    repository.URLRepository
	baseURL string
}

// NewURLService creates a new URL service
func NewURLService(repo repository.URLRepository, baseURL string) URLService {
	return &urlService{
		repo:    repo,
		baseURL: baseURL,
	}
}

// ShortenURL creates a shortened URL with automatic collision handling
func (s *urlService) ShortenURL(originalURL string, alias string, userID int64) (*domain.URL, error) {
	// Validate original URL
	if err := domain.ValidateURL(originalURL); err != nil {
		return nil, err
	}

	url := &domain.URL{
		OriginalURL: originalURL,
		UserID:      userID,
		ClickCount:  0,
	}

	// If custom alias is provided, use it directly
	if alias != "" {
		url.Alias = alias
		if err := s.repo.Create(url); err != nil {
			return nil, err
		}
		return url, nil
	}

	// Generate random alias with collision retry
	var lastErr error
	for i := 0; i < MaxRetries; i++ {
		generatedAlias, err := GenerateShortCode(DefaultCodeLength)
		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}

		url.Alias = generatedAlias
		if err := s.repo.Create(url); err != nil {
			if errors.Is(err, repository.ErrDuplicateAlias) {
				// Collision detected, retry with new code
				lastErr = err
				continue
			}
			return nil, err
		}

		// Success!
		return url, nil
	}

	// Max retries exceeded
	return nil, fmt.Errorf("%w: %v", ErrMaxRetriesExceeded, lastErr)
}

// GetURLByAlias retrieves URL information by alias
func (s *urlService) GetURLByAlias(alias string) (*domain.URL, error) {
	url, err := s.repo.FindByAlias(alias)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}
	return url, nil
}

// IncrementClickCount atomically increments the click counter
func (s *urlService) IncrementClickCount(alias string) error {
	return s.repo.IncrementClickCount(alias)
}

// ListURLs retrieves all URLs with pagination
func (s *urlService) ListURLs(limit, offset int) ([]*domain.URL, error) {
	// Set default limit if not specified
	if limit <= 0 {
		limit = 50
	}

	// Prevent excessive limit
	if limit > 100 {
		limit = 100
	}

	return s.repo.FindAll(limit, offset)
}

// GetURLsByUserID retrieves all URLs created by a specific user with pagination
func (s *urlService) GetURLsByUserID(userID int64, limit, offset int) ([]*domain.URL, error) {
	// Set default limit if not specified
	if limit <= 0 {
		limit = 50
	}

	// Prevent excessive limit
	if limit > 100 {
		limit = 100
	}

	return s.repo.FindByUserID(userID, limit, offset)
}
