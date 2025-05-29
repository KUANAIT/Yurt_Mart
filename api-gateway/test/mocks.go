package test

import (
	"api-gateway/internal/core/ports"
	"context"
)

type MockUserService struct {
	users map[string]*ports.UserResponse
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		users: make(map[string]*ports.UserResponse),
	}
}

func (m *MockUserService) RegisterUser(ctx context.Context, email, password, name string) (string, error) {
	userID := "test-user-" + email
	m.users[userID] = &ports.UserResponse{
		ID:    userID,
		Email: email,
		Name:  name,
	}
	return userID, nil
}

func (m *MockUserService) GetUser(ctx context.Context, userID string) (*ports.UserResponse, error) {
	if user, ok := m.users[userID]; ok {
		return user, nil
	}
	return nil, ports.ErrUserNotFound
}

func (m *MockUserService) Close() error {
	return nil
}

type MockCache struct {
	store map[string]interface{}
}

func NewMockCache() *MockCache {
	return &MockCache{
		store: make(map[string]interface{}),
	}
}

func (m *MockCache) Get(ctx context.Context, key string, value interface{}) (bool, error) {
	if _, ok := m.store[key]; ok {
		return true, nil
	}
	return false, nil
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}) error {
	m.store[key] = value
	return nil
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	delete(m.store, key)
	return nil
}

type MockEventPublisher struct{}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{}
}

func (m *MockEventPublisher) PublishUserEvent(event *ports.UserEvent) error {
	return nil
}
