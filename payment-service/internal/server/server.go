package server

import (
	"fmt"
	"net"

	"github.com/hsibAD/payment-service/internal/config"
	"github.com/hsibAD/payment-service/internal/handler"
	"google.golang.org/grpc"
)

type Server struct {
	cfg    *config.Config
	server *grpc.Server
}

func NewServer(cfg *config.Config) (*Server, error) {
	server := grpc.NewServer()
	
	// Register services
	handler.RegisterServices(server, cfg)

	return &Server{
		cfg:    cfg,
		server: server,
	}, nil
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return s.server.Serve(lis)
} 