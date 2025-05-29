package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"shopping-cart-service/internal/model"
	"shopping-cart-service/internal/service"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) AddToCart(ctx context.Context, item model.CartItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}
func (m *mockRepo) GetCart(ctx context.Context, userID string) ([]model.CartItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.CartItem), args.Error(1)
}
func (m *mockRepo) RemoveFromCart(ctx context.Context, userID string, productID string) error {
	args := m.Called(ctx, userID, productID)
	return args.Error(0)
}

// === ТЕСТЫ ===

func TestAddToCart(t *testing.T) {
	repo := new(mockRepo)
	svc := service.NewCartService(repo)

	item := model.CartItem{UserID: "u1", ProductID: "p1", Quantity: 1}
	repo.On("AddToCart", mock.Anything, item).Return(nil)

	err := svc.AddToCart(context.Background(), item)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetCart_CacheMiss(t *testing.T) {
	repo := new(mockRepo)
	svc := service.NewCartService(repo)

	userID := "u2"
	expectedItems := []model.CartItem{{UserID: userID, ProductID: "p1", Quantity: 2}}
	repo.On("GetCart", mock.Anything, userID).Return(expectedItems, nil)

	items, err := svc.GetCart(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
	repo.AssertExpectations(t)
}

func TestRemoveFromCart(t *testing.T) {
	repo := new(mockRepo)
	svc := service.NewCartService(repo)

	userID := "u3"
	productID := "p2"
	repo.On("RemoveFromCart", mock.Anything, userID, productID).Return(nil)

	err := svc.RemoveFromCart(context.Background(), userID, productID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
