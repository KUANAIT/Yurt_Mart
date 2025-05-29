package main

import (
	"context"
	//"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/internal/adapters/cache"
	"api-gateway/internal/adapters/clients"
	"api-gateway/internal/adapters/handlers"
	"api-gateway/internal/config"
	"api-gateway/internal/core/services"
	"api-gateway/internal/infrastructure/nats"
	"api-gateway/internal/infrastructure/server"
	"api-gateway/pkg/logger"
	"github.com/go-redis/redis/v8"
	natsio "github.com/nats-io/nats.go"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.NewLogger()

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Redis
	log.Info("Connecting to Redis at %s", cfg.RedisAddress)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test Redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatal("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize NATS
	log.Info("Connecting to NATS at %s", cfg.NATSAddress)
	nc, err := natsio.Connect(cfg.NATSAddress,
		natsio.MaxReconnects(5),
		natsio.ReconnectWait(2*time.Second),
		natsio.DisconnectErrHandler(func(conn *natsio.Conn, err error) {
			log.Error("NATS disconnected: %v", err)
		}),
		natsio.ReconnectHandler(func(conn *natsio.Conn) {
			log.Info("Reconnected to NATS at %s", conn.ConnectedUrl())
		}),
	)
	if err != nil {
		log.Fatal("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Initialize User Service Client
	log.Info("Connecting to User Service at %s", cfg.UserServiceURL)
	userClient, err := clients.NewUserServiceClient(cfg.UserServiceURL)
	if err != nil {
		log.Fatal("Failed to create user service client: %v", err)
	}
	defer userClient.Close()

	// Initialize Adapters
	redisCache := cache.NewRedisCache(redisClient, 30*time.Minute)
	natsPublisher := nats.NewPublisher(nc)

	// Initialize Services
	gatewayService := services.NewGatewayService(
		userClient,
		redisCache,
		natsPublisher,
	)

	// Initialize HTTP Handler
	handler := handlers.NewHTTPHandler(gatewayService, log)

	// Initialize HTTP Server
	httpServer := server.NewServer(handler, log, cfg.HTTPPort)

	// Graceful shutdown channel
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server in a goroutine
	go func() {
		log.Info("Starting HTTP server on port %d", cfg.HTTPPort)
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-shutdown
	log.Info("Shutdown signal received")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	log.Info("Shutting down server gracefully")
	if err := httpServer.Shutdown(shutdownCtx); err != nil && err != http.ErrServerClosed {
		log.Error("Graceful shutdown failed: %v", err)
	} else {
		log.Info("Server stopped gracefully")
	}
}
