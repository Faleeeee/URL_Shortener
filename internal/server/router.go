package server

import (
	"time"

	"github.com/Faleeeee/URL_Shortener/internal/database"
	"github.com/Faleeeee/URL_Shortener/internal/handler"
	"github.com/Faleeeee/URL_Shortener/internal/middleware"
	"github.com/Faleeeee/URL_Shortener/internal/repository"
	"github.com/Faleeeee/URL_Shortener/internal/service"
	"github.com/Faleeeee/URL_Shortener/internal/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(db *database.DB, baseURL, jwtSecret string, jwtExpiration time.Duration, base62Chars string) *gin.Engine {
	r := gin.Default()

	// Initialize JWT Manager
	jwtManager := utils.NewJWTManager(jwtSecret, jwtExpiration)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(jwtManager)

	// Initialize URL layers
	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo, baseURL, base62Chars)
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
	r.GET("/admin/url", urlHandler.ListURLs)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
