package service

import (
	"Url-Shortener-Service/internal/domain"
	"Url-Shortener-Service/internal/repository"
)

type UserService interface {
	GetUsers() ([]domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUsers() ([]domain.User, error) {
	return s.repo.GetAll()
}
