package repository

import "Url-Shortener-Service/internal/domain"

type UserRepository interface {
	GetAll() ([]domain.User, error)
}

type userRepo struct{}

func NewUserRepository() UserRepository {
	return &userRepo{}
}

func (r *userRepo) GetAll() ([]domain.User, error) {
	users := []domain.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	return users, nil
}
