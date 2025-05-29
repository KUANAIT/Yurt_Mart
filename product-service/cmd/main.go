package main

import (
	"log"
	"net"

	grpcServer "google.golang.org/grpc"
	productGrpc "product-service/internal/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpcServer.NewServer()
	productGrpc.RegisterProductServiceServer(s)

	log.Println("Product service is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
