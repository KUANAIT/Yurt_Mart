package service

import (
	"context"
	"log"
	"shopping-cart-service/internal/cache"
	"shopping-cart-service/internal/model"
	"shopping-cart-service/internal/repository"
)

type CartService struct {
	repo  repository.CartRepositoryInterface
	cache *cache.CartCache
}

func NewCartService(repo repository.CartRepositoryInterface) *CartService {
	c := cache.NewCartCache()
	return &CartService{repo: repo, cache: c}
}

func (s *CartService) AddToCart(ctx context.Context, item model.CartItem) error {
	err := s.repo.AddToCart(ctx, item)
	if err == nil {
		log.Printf("[CACHE] Invalidated → user_id: %s", item.UserID)
		s.cache.Invalidate(item.UserID)
	}
	return err
}

func (s *CartService) GetCart(ctx context.Context, userID string) ([]model.CartItem, error) {
	if items, found := s.cache.Get(userID); found {
		log.Println("[CACHE] Returning cart from cache for user:", userID)
		return items, nil
	}
	items, err := s.repo.GetCart(ctx, userID)
	if err == nil {
		s.cache.Set(userID, items)
		log.Printf("[CACHE] Set → user_id: %s, items_count: %d", userID, len(items))
	}
	return items, err
}

func (s *CartService) RemoveFromCart(ctx context.Context, userID, productID string) error {
	err := s.repo.RemoveFromCart(ctx, userID, productID)
	if err == nil {
		log.Println(" [CACHE] invalidated after RemoveFromCart for user:", userID)
		s.cache.Invalidate(userID)
	}
	return err
}
