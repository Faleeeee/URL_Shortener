package server

import (
	"Url-Shortener-Service/internal/database"
	"Url-Shortener-Service/internal/handler"
	"Url-Shortener-Service/internal/middleware"
	"Url-Shortener-Service/internal/repository"
	"Url-Shortener-Service/internal/service"
	"Url-Shortener-Service/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(db *database.DB, baseURL, jwtSecret string, jwtExpiration time.Duration) *gin.Engine {
	r := gin.Default()

	// Initialize JWT Manager
	jwtManager := utils.NewJWTManager(jwtSecret, jwtExpiration)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(jwtManager)

	// Initialize URL layers
	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo, baseURL)
	urlHandler := handler.NewURLHandler(urlService, baseURL)

	// Initialize Auth layers
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	// Authentication routes
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	// Public URL shortener routes
	r.GET("/:alias", urlHandler.RedirectURL)

	// Protected URL shortener routes (require authentication)
	r.POST("/url/shorten", authMiddleware, urlHandler.ShortenURL)
	r.GET("/url/links/:alias", authMiddleware, urlHandler.GetURLInfo)
	r.GET("/url/my-links", authMiddleware, urlHandler.GetUserURLs)

	// Admin routes (require authentication)
	r.GET("/admin/url", authMiddleware, urlHandler.ListURLs)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
