package auth

import (
	"context"
	"sync"
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
	blacklist   map[string]time.Time
	mu          sync.RWMutex
}

func NewService(jwtSecret string, userService services.UserService) *Service {
	return &Service{
		jwtSecret:   []byte(jwtSecret),
		userService: userService,
		blacklist:   make(map[string]time.Time),
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

func (s *Service) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Parse the token to get expiration time
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.InvalidArgument, "invalid token signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid token claims")
	}

	// Add token to blacklist with its expiration time
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid token expiration")
	}

	s.mu.Lock()
	s.blacklist[req.Token] = time.Unix(int64(exp), 0)
	s.mu.Unlock()

	return &pb.LogoutResponse{
		Success: true,
	}, nil
}

// Helper method to check if a token is blacklisted
func (s *Service) isTokenBlacklisted(token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if exp, exists := s.blacklist[token]; exists {
		if time.Now().After(exp) {
			// Token has expired, remove from blacklist
			s.mu.RUnlock()
			s.mu.Lock()
			delete(s.blacklist, token)
			s.mu.Unlock()
			s.mu.RLock()
			return false
		}
		return true
	}
	return false
}
