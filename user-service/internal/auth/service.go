package auth

import (
	"context"
	//"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user-service/internal/core/services"
	pb "user-service/proto/auth"
)

type Service struct {
	pb.UnimplementedAuthServiceServer
	jwtSecret   []byte
	userService services.UserService
}

func NewService(jwtSecret string, userService services.UserService) *Service {
	return &Service{
		jwtSecret:   []byte(jwtSecret),
		userService: userService,
	}
}

func (s *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	user, err := s.userService.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	if !s.userService.VerifyPassword(user.Password, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.LoginResponse{
		Token:  tokenString,
		UserId: user.ID,
		Email:  user.Email,
	}, nil
}
