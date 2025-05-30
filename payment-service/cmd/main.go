package main

import (
	"log"

	"github.com/hsibAD/payment-service/internal/config"
	"github.com/hsibAD/payment-service/internal/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create and start server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
} 