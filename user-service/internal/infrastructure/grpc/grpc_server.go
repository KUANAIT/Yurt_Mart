package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
	"user-service/internal/auth"
	"user-service/internal/core/services"
	"user-service/internal/metrics"
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
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("RegisterUser").Observe(duration)
	}()

	newUser, err := s.userService.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		metrics.ErrorCount.WithLabelValues("RegisterUser", "registration_failed").Inc()
		return nil, err
	}

	metrics.RequestCount.WithLabelValues("RegisterUser", "success").Inc()
	metrics.ActiveUsers.Inc()

	return &user.RegisterUserResponse{
		UserId: newUser.ID,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("GetUser").Observe(duration)
	}()

	fetchedUser, err := s.userService.GetByID(ctx, req.UserId)
	if err != nil {
		metrics.ErrorCount.WithLabelValues("GetUser", "not_found").Inc()
		return nil, err
	}

	metrics.RequestCount.WithLabelValues("GetUser", "success").Inc()

	return &user.GetUserResponse{
		UserId: fetchedUser.ID,
		Email:  fetchedUser.Email,
		Name:   fetchedUser.Name,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("UpdateUser").Observe(duration)
	}()

	updatedUser, err := s.userService.UpdateUser(ctx, req.UserId, req.Email, req.Name)
	if err != nil {
		metrics.ErrorCount.WithLabelValues("UpdateUser", "update_failed").Inc()
		return nil, err
	}

	metrics.RequestCount.WithLabelValues("UpdateUser", "success").Inc()

	return &user.UpdateUserResponse{
		UserId: updatedUser.ID,
		Email:  updatedUser.Email,
		Name:   updatedUser.Name,
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("DeleteUser").Observe(duration)
	}()

	err := s.userService.DeleteUser(ctx, req.UserId)
	if err != nil {
		metrics.ErrorCount.WithLabelValues("DeleteUser", "delete_failed").Inc()
		return nil, err
	}

	metrics.RequestCount.WithLabelValues("DeleteUser", "success").Inc()
	metrics.ActiveUsers.Dec()

	return &user.DeleteUserResponse{
		Success: true,
	}, nil
}

func StartGRPCServer(port int, userService services.UserService) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Start metrics HTTP server
	go func() {
		metricsMux := http.NewServeMux()
		metricsMux.Handle("/metrics", promhttp.Handler())
		metricsServer := &http.Server{
			Addr:    ":9090",
			Handler: metricsMux,
		}
		if err := metricsServer.ListenAndServe(); err != nil {
			fmt.Printf("Failed to start metrics server: %v\n", err)
		}
	}()

	authService := auth.NewService("your-secret-key-here", userService)

	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, NewGRPCServer(userService, authService))
	authpb.RegisterAuthServiceServer(grpcServer, authService)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
