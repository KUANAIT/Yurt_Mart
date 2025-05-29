package main

import (
	"log"
	"net"

	"meaningfullname/Yurt_Mart/common/config"
	"meaningfullname/Yurt_Mart/common/database"
	"meaningfullname/Yurt_Mart/orderhistory/internal/repository"
	"meaningfullname/Yurt_Mart/orderhistory/internal/service"
	"meaningfullname/Yurt_Mart/orderhistory/proto"

	"google.golang.org/grpc"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.InitMongo(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Ensure database is connected
	if database.GetDB() == nil {
		log.Fatalf("Database connection is not initialized")
	}

	// Create repository and service
	repo := repository.NewOrderRepository()
	if repo == nil {
		log.Fatalf("Failed to create repository")
	}
	orderService := service.NewOrderHistoryService(repo)

	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterOrderHistoryServiceServer(grpcServer, orderService)

	log.Printf("Starting Order History Service on :8086")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
