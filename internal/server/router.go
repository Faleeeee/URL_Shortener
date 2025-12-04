package server

import (
	"Url-Shortener-Service/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// user routes
	userHandler := handler.NewUserHandler()
	r.GET("/users", userHandler.GetUsers)

	return r
}
