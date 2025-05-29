package main

import (
	"log"
	"net"

	"shopping-cart-service/database"
	"shopping-cart-service/events"
	"shopping-cart-service/internal/handler"
	"shopping-cart-service/internal/repository"
	"shopping-cart-service/internal/service"
	pb "shopping-cart-service/shopping-cart-service/proto/cartpb"

	"google.golang.org/grpc"
)

func main() {
	client := database.ConnectMongo("mongodb://localhost:27017")
	defer client.Disconnect(nil)

	events.InitNATS("nats://localhost:4222")
	events.SubscribeToCartEvents()

	db := client.Database("shop")
	repo := repository.NewCartRepository(db)
	svc := service.NewCartService(repo)
	h := handler.NewCartHandler(svc)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCartServiceServer(s, h)

	log.Println("Shopping Cart Service is running on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
