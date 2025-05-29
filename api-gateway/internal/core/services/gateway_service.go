package services

import (
	"context"
	"encoding/json"
	//"errors"
	//"time"

	"api-gateway/internal/core/ports"
)

type GatewayService struct {
	userClient     ports.UserServicePort
	cache          ports.CachePort
	eventPublisher ports.EventPublisherPort
}

func NewGatewayService(
	userClient ports.UserServicePort,
	cache ports.CachePort,
	publisher ports.EventPublisherPort,
) *GatewayService {
	return &GatewayService{
		userClient:     userClient,
		cache:          cache,
		eventPublisher: publisher,
	}
}

func (s *GatewayService) RegisterUser(ctx context.Context, email, password, name string) (string, error) {
	// Check cache first
	var cachedID string
	if found, _ := s.cache.Get(ctx, "user:email:"+email, &cachedID); found {
		return cachedID, nil
	}

	userID, err := s.userClient.RegisterUser(ctx, email, password, name)
	if err != nil {
		return "", err
	}

	// Cache the response
	s.cache.Set(ctx, "user:email:"+email, userID)
	s.cache.Set(ctx, "user:"+userID, map[string]interface{}{
		"email": email,
		"name":  name,
	})

	// Publish event
	event := &ports.UserEvent{
		Type:   "user.registered",
		UserID: userID,
	}
	if payload, err := json.Marshal(map[string]string{"email": email}); err == nil {
		event.Payload = payload
	}
	s.eventPublisher.PublishUserEvent(event)

	return userID, nil
}

func (s *GatewayService) GetUser(ctx context.Context, userID string) (*ports.UserResponse, error) {
	// Check cache first
	var cachedUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if found, _ := s.cache.Get(ctx, "user:"+userID, &cachedUser); found {
		return &ports.UserResponse{
			ID:    userID,
			Email: cachedUser.Email,
			Name:  cachedUser.Name,
		}, nil
	}

	user, err := s.userClient.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache the response
	s.cache.Set(ctx, "user:"+userID, map[string]interface{}{
		"email": user.Email,
		"name":  user.Name,
	})
	s.cache.Set(ctx, "user:email:"+user.Email, userID)

	return user, nil
}

func (s *GatewayService) UpdateUser(ctx context.Context, userID, email, name string) error {
	// Invalidate cache
	s.cache.Delete(ctx, "user:"+userID)
	s.cache.Delete(ctx, "user:email:"+email)
	return nil
}
