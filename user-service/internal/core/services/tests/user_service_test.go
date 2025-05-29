package tests

//
//import (
//	"context"
//	//"errors"
//	"testing"
//
//	"user-service/internal/core/domain"
//	"user-service/internal/core/services"
//)
//
//type MockUserRepository struct {
//	users map[string]*domain.User
//}
//
//func NewMockUserRepository() *MockUserRepository {
//	return &MockUserRepository{
//		users: make(map[string]*domain.User),
//	}
//}
//
//func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
//	m.users[user.ID] = user
//	return nil
//}
//
//func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
//	if user, exists := m.users[id]; exists {
//		return user, nil
//	}
//	return nil, nil
//}
//
//func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
//	for _, user := range m.users {
//		if user.Email == email {
//			return user, nil
//		}
//	}
//	return nil, nil
//}
//
//func TestUserService_Register(t *testing.T) {
//	tests := []struct {
//		name          string
//		email         string
//		password      string
//		userName      string
//		existingUser  *domain.User
//		expectedError bool
//	}{
//		{
//			name:     "successful registration",
//			email:    "test@example.com",
//			password: "password123",
//			userName: "Test User",
//		},
//		{
//			name:     "duplicate email",
//			email:    "existing@example.com",
//			password: "password123",
//			userName: "Test User",
//			existingUser: &domain.User{
//				ID:       "1",
//				Email:    "existing@example.com",
//				Name:     "Existing User",
//				Password: "hashedpass",
//			},
//			expectedError: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo := NewMockUserRepository()
//			if tt.existingUser != nil {
//				repo.Create(context.Background(), tt.existingUser)
//			}
//
//			service := services.NewUserService(repo)
//			user, err := service.Register(context.Background(), tt.email, tt.password, tt.userName)
//
//			if tt.expectedError {
//				if err == nil {
//					t.Error("expected error but got none")
//				}
//				return
//			}
//
//			if err != nil {
//				t.Errorf("unexpected error: %v", err)
//				return
//			}
//
//			if user.Email != tt.email {
//				t.Errorf("expected email %s, got %s", tt.email, user.Email)
//			}
//
//			if user.Name != tt.userName {
//				t.Errorf("expected name %s, got %s", tt.userName, user.Name)
//			}
//
//			if user.Password != tt.password {
//				t.Errorf("expected password %s, got %s", tt.password, user.Password)
//			}
//		})
//	}
//}
//
//func TestUserService_GetByID(t *testing.T) {
//	repo := NewMockUserRepository()
//	service := services.NewUserService(repo)
//
//	// Create a test user
//	testUser := &domain.User{
//		ID:       "1",
//		Email:    "test@example.com",
//		Name:     "Test User",
//		Password: "password123",
//	}
//	repo.Create(context.Background(), testUser)
//
//	tests := []struct {
//		name          string
//		id            string
//		expectedError bool
//		errorMessage  string
//	}{
//		{
//			name: "existing user",
//			id:   "1",
//		},
//		{
//			name:          "non-existing user",
//			id:            "999",
//			expectedError: true,
//			errorMessage:  "user not found",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			user, err := service.GetByID(context.Background(), tt.id)
//
//			if tt.expectedError {
//				if err == nil {
//					t.Error("expected error but got none")
//					return
//				}
//				if err.Error() != tt.errorMessage {
//					t.Errorf("expected error message %q, got %q", tt.errorMessage, err.Error())
//				}
//				return
//			}
//
//			if err != nil {
//				t.Errorf("unexpected error: %v", err)
//				return
//			}
//
//			if user.ID != tt.id {
//				t.Errorf("expected user ID %s, got %s", tt.id, user.ID)
//			}
//		})
//	}
//}
//
//func TestUserService_GetByEmail(t *testing.T) {
//	repo := NewMockUserRepository()
//	service := services.NewUserService(repo)
//
//
//	testUser := &domain.User{
//		ID:       "1",
//		Email:    "test@example.com",
//		Name:     "Test User",
//		Password: "password123",
//	}
//	repo.Create(context.Background(), testUser)
//
//	tests := []struct {
//		name          string
//		email         string
//		expectedError bool
//		errorMessage  string
//	}{
//		{
//			name:  "existing user",
//			email: "test@example.com",
//		},
//		{
//			name:          "non-existing user",
//			email:         "nonexistent@example.com",
//			expectedError: true,
//			errorMessage:  "user not found",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			user, err := service.GetByEmail(context.Background(), tt.email)
//
//			if tt.expectedError {
//				if err == nil {
//					t.Error("expected error but got none")
//					return
//				}
//				if err.Error() != tt.errorMessage {
//					t.Errorf("expected error message %q, got %q", tt.errorMessage, err.Error())
//				}
//				return
//			}
//
//			if err != nil {
//				t.Errorf("unexpected error: %v", err)
//				return
//			}
//
//			if user.Email != tt.email {
//				t.Errorf("expected email %s, got %s", tt.email, user.Email)
//			}
//		})
//	}
//}
