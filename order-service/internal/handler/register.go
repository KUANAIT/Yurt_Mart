package handler

import (
	"github.com/hsibAD/order-service/internal/config"
	"github.com/hsibAD/order-service/internal/domain"
	"github.com/hsibAD/order-service/internal/events"
	pb "github.com/hsibAD/order-service/proto"
	"google.golang.org/grpc"
)

func RegisterServices(server *grpc.Server, cfg *config.Config) error {
	// TODO: Initialize repositories
	var orderRepo domain.OrderRepository
	
	// Initialize cart subscriber
	cartSub, err := events.NewCartSubscriber(cfg.NatsURL)
	if err != nil {
		return err
	}

	// Create and register order handler
	orderHandler := NewOrderHandler(orderRepo, cartSub)
	pb.RegisterOrderServiceServer(server, orderHandler)
	
	return nil
} 