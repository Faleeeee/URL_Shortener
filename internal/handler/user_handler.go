package handler

import (
	"Url-Shortener-Service/internal/repository"
	"Url-Shortener-Service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler() *UserHandler {
	repo := repository.NewUserRepository()
	sv := service.NewUserService(repo)

	return &UserHandler{service: sv}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, _ := h.service.GetUsers()
	c.JSON(http.StatusOK, users)
}
