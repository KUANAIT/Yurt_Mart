package repository

import (
	"context"
	"shopping-cart-service/internal/model"
)

type CartRepositoryInterface interface {
	AddToCart(ctx context.Context, item model.CartItem) error
	GetCart(ctx context.Context, userID string) ([]model.CartItem, error)
	RemoveFromCart(ctx context.Context, userID string, productID string) error
}
