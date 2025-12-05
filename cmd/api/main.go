package main

import (
	_ "Url-Shortener-Service/docs"
	"Url-Shortener-Service/internal/config"
	"Url-Shortener-Service/internal/database"
	"Url-Shortener-Service/internal/server"
	"log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// Start server
	s := server.NewServer(cfg, db)
	s.Run()
}
