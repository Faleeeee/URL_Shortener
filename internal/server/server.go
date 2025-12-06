package server

import (
	"Url-Shortener-Service/internal/config"
	"Url-Shortener-Service/internal/database"
	"log"
	"net/http"
	"time"
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
	baseURL := "http://localhost:" + s.cfg.Server.Port

	jwtExpiration, err := time.ParseDuration(s.cfg.JWT.Expiration)
	if err != nil {
		log.Printf("Invalid JWT expiration format, using default 24h: %v", err)
		jwtExpiration = 24 * time.Hour
	}

	r := NewRouter(s.db, baseURL, s.cfg.JWT.Secret, jwtExpiration)

	log.Printf("Server running on %s", baseURL)
	log.Printf("Swagger docs: %s/swagger/index.html", baseURL)

	http.ListenAndServe(":"+s.cfg.Server.Port, r)
}
