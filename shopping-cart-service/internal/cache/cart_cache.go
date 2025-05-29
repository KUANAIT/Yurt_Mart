package cache

import (
	"log"
	"sync"
	"time"

	"shopping-cart-service/internal/model"
)

type CartCache struct {
	mu    sync.RWMutex
	items map[string][]model.CartItem
}

func NewCartCache() *CartCache {
	c := &CartCache{
		items: make(map[string][]model.CartItem),
	}
	// автоматическое обновление каждые 12 часов
	go func() {
		for {
			time.Sleep(12 * time.Hour)
			c.Clear()
			log.Println("🌀 Cart cache refreshed after 12 hours.")
		}
	}()
	return c
}

func (c *CartCache) Get(userID string) ([]model.CartItem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	items, found := c.items[userID]
	return items, found
}

func (c *CartCache) Set(userID string, items []model.CartItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[userID] = items
}

func (c *CartCache) Invalidate(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, userID)
}

func (c *CartCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string][]model.CartItem)
}
