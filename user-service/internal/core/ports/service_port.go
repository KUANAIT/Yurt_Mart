package ports

import (
	"context"
	"user-service/internal/core/domain"
)

type User struct {
	ID       string
	Email    string
	Name     string
	Password string
}

type UserService interface {
	Register(ctx context.Context, email, password, name string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(email string) (*domain.User, error)
	UpdateUser(ctx context.Context, id, email, name string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
}

type EmailService interface {
	SendWelcomeEmail(email, name string) error
	SendPasswordResetEmail(email, token string) error
}
