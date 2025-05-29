package cache

import (
	"product-service/internal/domain"
	"sync"
	"time"
)

type ProductCache struct {
	products    []*domain.Product
	productByID map[string]*domain.Product
	mu          sync.RWMutex
	lastRefresh time.Time
}

func NewProductCache() *ProductCache {
	return &ProductCache{
		productByID: make(map[string]*domain.Product),
	}
}

func (c *ProductCache) SetProducts(products []*domain.Product) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.products = products
	c.productByID = make(map[string]*domain.Product)
	for _, p := range products {
		c.productByID[p.ID] = p
	}
	c.lastRefresh = time.Now()
}

func (c *ProductCache) AddProduct(p *domain.Product) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.products = append(c.products, p)
	c.productByID[p.ID] = p
}

func (c *ProductCache) GetAll() []*domain.Product {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.products
}

func (c *ProductCache) GetByID(id string) (*domain.Product, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, found := c.productByID[id]
	return p, found
}

func (c *ProductCache) ShouldRefresh() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Since(c.lastRefresh) > 12*time.Hour
}
