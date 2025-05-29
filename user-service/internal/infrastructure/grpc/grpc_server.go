package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"user-service/internal/auth"
	"user-service/internal/core/services"
	authpb "user-service/proto/auth"
	"user-service/proto/user"
)

type Server struct {
	user.UnimplementedUserServiceServer
	userService services.UserService
	authService *auth.Service
}

func NewGRPCServer(userService services.UserService, authService *auth.Service) *Server {
	return &Server{
		userService: userService,
		authService: authService,
	}
}

func (s *Server) RegisterUser(ctx context.Context, req *user.RegisterUserRequest) (*user.RegisterUserResponse, error) {
	newUser, err := s.userService.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, err
	}

	return &user.RegisterUserResponse{
		UserId: newUser.ID,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	fetchedUser, err := s.userService.GetByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &user.GetUserResponse{
		UserId: fetchedUser.ID,
		Email:  fetchedUser.Email,
		Name:   fetchedUser.Name,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	updatedUser, err := s.userService.UpdateUser(ctx, req.UserId, req.Email, req.Name)
	if err != nil {
		return nil, err
	}

	return &user.UpdateUserResponse{
		UserId: updatedUser.ID,
		Email:  updatedUser.Email,
		Name:   updatedUser.Name,
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	err := s.userService.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &user.DeleteUserResponse{
		Success: true,
	}, nil
}

func StartGRPCServer(port int, userService services.UserService) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	authService := auth.NewService("your-secret-key-here", userService)

	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, NewGRPCServer(userService, authService))
	authpb.RegisterAuthServiceServer(grpcServer, authService)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
