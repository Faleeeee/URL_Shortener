package server

import (
	"Url-Shortener-Service/internal/database"
	"Url-Shortener-Service/internal/handler"
	"Url-Shortener-Service/internal/repository"
	"Url-Shortener-Service/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(db *database.DB, baseURL string) *gin.Engine {
	r := gin.Default()

	// Initialize layers
	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo, baseURL)
	urlHandler := handler.NewURLHandler(urlService, baseURL)

	// URL shortener routes
	r.POST("/url/shorten", urlHandler.ShortenURL)
	r.GET("/:alias", urlHandler.RedirectURL)
	r.GET("/url/links/:alias", urlHandler.GetURLInfo)
	r.GET("/url/links", urlHandler.ListURLs)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
