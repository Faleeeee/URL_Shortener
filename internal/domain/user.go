package domain

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
}

var (
	ErrInvalidUsername       = errors.New("username must be 3-64 characters and contain only alphanumeric characters, hyphens, and underscores")
	ErrInvalidPassword       = errors.New("password must be at least 8 characters")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

const (
	MinUsernameLength = 3
	MaxUsernameLength = 64
	MinPasswordLength = 8
)

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)

	if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
		return ErrInvalidUsername
	}

	for _, char := range username {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return ErrInvalidUsername
		}
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ErrInvalidPassword
	}
	return nil
}
