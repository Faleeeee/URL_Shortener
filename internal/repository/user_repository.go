package repository

import (
	"Url-Shortener-Service/internal/database"
	"Url-Shortener-Service/internal/domain"
	"database/sql"
	"errors"
)

var (
	ErrDuplicateUsername = errors.New("username already exists")
)

// UserRepository handles user data access
type UserRepository struct {
	db *database.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(username, hashedPassword string) (*domain.User, error) {
	query := `
		INSERT INTO users (username, password, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, username, created_at, updated_at
	`

	user := &domain.User{
		Password: hashedPassword,
	}

	err := r.db.QueryRow(
		query,
		username,
		hashedPassword,
	).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		// Check for duplicate username constraint violation
		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return nil, ErrDuplicateUsername
		}
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, username, password, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int64) (*domain.User, error) {
	query := `
		SELECT id, username, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}
