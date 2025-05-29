package ports

import (
	"context"
	//"time"
)

type UserServicePort interface {
	RegisterUser(ctx context.Context, email, password, name string) (string, error)
	GetUser(ctx context.Context, userID string) (*UserResponse, error)
}

type CachePort interface {
	Get(ctx context.Context, key string, dest interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
}

type EventPublisherPort interface {
	PublishUserEvent(event *UserEvent) error
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserEvent struct {
	Type    string `json:"type"`
	UserID  string `json:"user_id"`
	Payload []byte `json:"payload,omitempty"`
}
