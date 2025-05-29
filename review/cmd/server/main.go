package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"meaningfullname/Yurt_Mart/review/internal/repository"
	"meaningfullname/Yurt_Mart/review/internal/service"
	"meaningfullname/Yurt_Mart/review/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := repository.NewReviewRepository()
	reviewService := service.NewReviewService(repo)

	srv := grpc.NewServer()
	proto.RegisterReviewServiceServer(srv, reviewService)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	srv.GracefulStop()
}
