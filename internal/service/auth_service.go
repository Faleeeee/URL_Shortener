package service

import (
	"github.com/Faleeeee/URL_Shortener/internal/domain"
	"github.com/Faleeeee/URL_Shortener/internal/repository"
	"github.com/Faleeeee/URL_Shortener/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *utils.JWTManager
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, jwtManager *utils.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register creates a new user account
func (s *AuthService) Register(username, password string) (*domain.AuthResponse, error) {
	// Validate username
	if err := domain.ValidateUsername(username); err != nil {
		return nil, err
	}

	// Validate password
	if err := domain.ValidatePassword(password); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user, err := s.userRepo.CreateUser(username, string(hashedPassword))
	if err != nil {
		if err == repository.ErrDuplicateUsername {
			return nil, domain.ErrUsernameAlreadyExists
		}
		return nil, err
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:    token,
		Username: user.Username,
		UserID:   user.ID,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(username, password string) (*domain.AuthResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:    token,
		Username: user.Username,
		UserID:   user.ID,
	}, nil
}
