package server

import (
	"Url-Shortener-Service/internal/config"
	"Url-Shortener-Service/internal/database"
	"log"
	"net/http"
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
	r := NewRouter(s.db, baseURL)

	log.Printf("Server running on %s", baseURL)
	log.Printf("Swagger docs: %s/swagger/index.html", baseURL)

	http.ListenAndServe(":"+s.cfg.Server.Port, r)
}
