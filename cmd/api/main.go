package main

import (
	"Url-Shortener-Service/internal/config"
	"Url-Shortener-Service/internal/server"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	s := server.NewServer(cfg)
	log.Println("Server running on port", cfg.Server.Port)
	s.Run()
}
