package domain

import (
	"context"
)

type User struct {
	ID       string `bson:"_id"`
	Email    string `bson:"email"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
