package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	//"user-service/internal/adapters/cache"
	"user-service/internal/adapters/handlers"
	"user-service/internal/adapters/repositories"
	"user-service/internal/config"
	"user-service/internal/core/services"
	"user-service/internal/infrastructure/database"
	"user-service/internal/infrastructure/email"
	"user-service/internal/infrastructure/grpc"
	infrastructure_nats "user-service/internal/infrastructure/nats"
	"user-service/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	logger := logger.NewLogger()

	mongoClient, err := database.NewMongoClient(cfg.MongoURI)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := database.CloseMongoClient(mongoClient); err != nil {
			logger.Error("Failed to close MongoDB connection: %v", err)
		}
	}()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	nc, err := nats.Connect(cfg.NATSAddress)
	if err != nil {
		logger.Fatal("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	userRepo := repositories.NewMongoUserRepository(mongoClient, cfg.DBName)

	emailService := email.NewSMTPSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPUsername,
	)
	userService := services.NewUserService(userRepo)

	eventHandler := handlers.NewEventHandler(emailService)
	subscriber := infrastructure_nats.NewSubscriber(nc, eventHandler)
	if err := subscriber.Subscribe(); err != nil {
		logger.Fatal("Failed to subscribe to NATS: %v", err)
	}

	logger.Info("Starting gRPC server on port %d", cfg.GRPCPort)
	if err := grpc.StartGRPCServer(cfg.GRPCPort, userService); err != nil {
		logger.Fatal("Failed to start gRPC server: %v", err)
	}
}
