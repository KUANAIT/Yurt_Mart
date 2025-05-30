package server

import (
	"fmt"
	"net"

	"github.com/hsibAD/order-service/internal/config"
	"github.com/hsibAD/order-service/internal/handler"
	"google.golang.org/grpc"
)

type Server struct {
	cfg    *config.Config
	server *grpc.Server
}

func NewServer(cfg *config.Config) (*Server, error) {
	server := grpc.NewServer()

	// Register services
	if err := handler.RegisterServices(server, cfg); err != nil {
		return nil, fmt.Errorf("failed to register services: %v", err)
	}

	return &Server{
		cfg:    cfg,
		server: server,
	}, nil
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return s.server.Serve(lis)
}
