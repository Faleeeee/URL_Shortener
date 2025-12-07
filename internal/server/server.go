package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Faleeeee/URL_Shortener/internal/database"

	"github.com/Faleeeee/URL_Shortener/internal/config"
)

type Server struct {
	cfg *config.Config
	db  *database.DB
}

func NewServer(cfg *config.Config, db *database.DB) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) Run() {
	baseURL := s.cfg.Server.BaseURL

	jwtExpiration, err := time.ParseDuration(s.cfg.JWT.Expiration)
	if err != nil {
		log.Printf("Invalid JWT expiration format, using default 24h: %v", err)
		jwtExpiration = 24 * time.Hour
	}

	r := NewRouter(s.db, baseURL, s.cfg.JWT.Secret, jwtExpiration, s.cfg.Shortener.Base62Chars)

	log.Printf("Server running on %s", baseURL)
	log.Printf("Swagger docs: %s/swagger/index.html", baseURL)

	http.ListenAndServe(":"+s.cfg.Server.Port, r)
}
