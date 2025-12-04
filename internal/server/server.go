package server

import (
	"Url-Shortener-Service/internal/config"
	"net/http"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Run() {
	r := NewRouter()

	http.ListenAndServe(":"+s.cfg.Server.Port, r)
}
