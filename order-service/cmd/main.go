package main

import (
	"log"
	"time"

	"github.com/hsibAD/order-service/internal/config"
	"github.com/hsibAD/order-service/internal/server"
)

func main() {
	log.Println("Order-service starting...")
	// Load configuration
	cfg := config.Load()

	// Create and start server with retries
	var srv *server.Server
	var err error
	for i := 0; i < 5; i++ {
		srv, err = server.NewServer(cfg)
		if err == nil {
			log.Println("Order-service: gRPC server created successfully.")
			break
		}
		log.Printf("Failed to create server, attempt %d: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		log.Fatalf("Failed to create server after retries: %v", err)
	}

	log.Println("Order-service: running gRPC server...")
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
