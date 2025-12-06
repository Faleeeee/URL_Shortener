package repository

import (
	"Url-Shortener-Service/internal/database"
	"Url-Shortener-Service/internal/domain"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNotFound       = errors.New("URL not found")
	ErrDuplicateAlias = errors.New("alias already exists")
)

// URLRepository defines the interface for URL data access
type URLRepository interface {
	Create(url *domain.URL) error
	FindByAlias(alias string) (*domain.URL, error)
	IncrementClickCount(alias string) error
	FindAll(limit, offset int) ([]*domain.URL, error)
	FindByUserID(userID int64, limit, offset int) ([]*domain.URL, error)
	ExistsByAlias(alias string) (bool, error)
}

type urlRepository struct {
	db *database.DB
}

// NewURLRepository creates a new URL repository
func NewURLRepository(db *database.DB) URLRepository {
	return &urlRepository{db: db}
}

// Create inserts a new URL into the database
func (r *urlRepository) Create(url *domain.URL) error {
	query := `
		INSERT INTO urls (alias, original_url, create_id, click_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		url.Alias,
		url.OriginalURL,
		url.UserID,
		url.ClickCount,
	).Scan(&url.ID, &url.CreatedAt, &url.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "pq: duplicate key value violates unique constraint \"urls_alias_key\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"idx_alias\"" {
			return ErrDuplicateAlias
		}
		return fmt.Errorf("failed to create URL: %w", err)
	}

	return nil
}

// FindByAlias retrieves a URL by its alias
func (r *urlRepository) FindByAlias(alias string) (*domain.URL, error) {
	query := `
		SELECT id, alias, original_url, create_id, click_count, created_at, updated_at
		FROM urls
		WHERE alias = $1
	`

	url := &domain.URL{}
	err := r.db.QueryRow(query, alias).Scan(
		&url.ID,
		&url.Alias,
		&url.OriginalURL,
		&url.UserID,
		&url.ClickCount,
		&url.CreatedAt,
		&url.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find URL: %w", err)
	}

	return url, nil
}

// IncrementClickCount atomically increments the click counter for a URL
func (r *urlRepository) IncrementClickCount(alias string) error {
	query := `
		UPDATE urls
		SET click_count = click_count + 1,
		    updated_at = NOW()
		WHERE alias = $1
	`

	result, err := r.db.Exec(query, alias)
	if err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// FindAll retrieves all URLs with pagination
func (r *urlRepository) FindAll(limit, offset int) ([]*domain.URL, error) {
	query := `
		SELECT id, alias, original_url, create_id, click_count, created_at, updated_at
		FROM urls
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query URLs: %w", err)
	}
	defer rows.Close()

	var urls []*domain.URL
	for rows.Next() {
		url := &domain.URL{}
		err := rows.Scan(
			&url.ID,
			&url.Alias,
			&url.OriginalURL,
			&url.UserID,
			&url.ClickCount,
			&url.CreatedAt,
			&url.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan URL: %w", err)
		}
		urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return urls, nil
}

// FindByUserID retrieves all URLs created by a specific user with pagination
func (r *urlRepository) FindByUserID(userID int64, limit, offset int) ([]*domain.URL, error) {
	query := `
		SELECT id, alias, original_url, create_id, click_count, created_at, updated_at
		FROM urls
		WHERE create_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query URLs by user ID: %w", err)
	}
	defer rows.Close()

	var urls []*domain.URL
	for rows.Next() {
		url := &domain.URL{}
		err := rows.Scan(
			&url.ID,
			&url.Alias,
			&url.OriginalURL,
			&url.UserID,
			&url.ClickCount,
			&url.CreatedAt,
			&url.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan URL: %w", err)
		}
		urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return urls, nil
}

// ExistsByAlias checks if an alias already exists
func (r *urlRepository) ExistsByAlias(alias string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE alias = $1)`

	var exists bool
	err := r.db.QueryRow(query, alias).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check alias existence: %w", err)
	}

	return exists, nil
}
