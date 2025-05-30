package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hsibAD/order-service/internal/domain"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	return value, nil
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Order-specific cache methods
func (c *RedisCache) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	key := "order:" + orderID
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var order domain.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (c *RedisCache) SetOrder(ctx context.Context, order *domain.Order, ttl int) error {
	key := "order:" + order.ID
	return c.Set(ctx, key, order, ttl)
}

func (c *RedisCache) DeleteOrder(ctx context.Context, orderID string) error {
	key := "order:" + orderID
	return c.Delete(ctx, key)
}

// Delivery address cache methods
func (c *RedisCache) GetDeliveryAddresses(ctx context.Context, userID string) ([]*domain.DeliveryAddress, error) {
	key := "addresses:" + userID
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var addresses []*domain.DeliveryAddress
	if err := json.Unmarshal(data, &addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (c *RedisCache) SetDeliveryAddresses(ctx context.Context, userID string, addresses []*domain.DeliveryAddress, ttl int) error {
	key := "addresses:" + userID
	return c.Set(ctx, key, addresses, ttl)
}

func (c *RedisCache) DeleteDeliveryAddresses(ctx context.Context, userID string) error {
	key := "addresses:" + userID
	return c.Delete(ctx, key)
}

// Delivery slots cache methods
func (c *RedisCache) GetDeliverySlots(ctx context.Context, date string) ([]*domain.DeliverySlot, error) {
	key := "slots:" + date
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var slots []*domain.DeliverySlot
	if err := json.Unmarshal(data, &slots); err != nil {
		return nil, err
	}

	return slots, nil
}

func (c *RedisCache) SetDeliverySlots(ctx context.Context, date string, slots []*domain.DeliverySlot, ttl int) error {
	key := "slots:" + date
	return c.Set(ctx, key, slots, ttl)
}

func (c *RedisCache) DeleteDeliverySlots(ctx context.Context, date string) error {
	key := "slots:" + date
	return c.Delete(ctx, key)
} 