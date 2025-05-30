package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hsibAD/payment-service/internal/domain"
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

// Payment-specific cache methods
func (c *RedisCache) GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error) {
	key := "payment:" + paymentID
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var payment domain.Payment
	if err := json.Unmarshal(data, &payment); err != nil {
		return nil, err
	}

	return &payment, nil
}

func (c *RedisCache) SetPayment(ctx context.Context, payment *domain.Payment, ttl int) error {
	key := "payment:" + payment.ID
	return c.Set(ctx, key, payment, ttl)
}

func (c *RedisCache) DeletePayment(ctx context.Context, paymentID string) error {
	key := "payment:" + paymentID
	return c.Delete(ctx, key)
}

// Order payments cache methods
func (c *RedisCache) GetOrderPayments(ctx context.Context, orderID string) ([]*domain.Payment, error) {
	key := "order_payments:" + orderID
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var payments []*domain.Payment
	if err := json.Unmarshal(data, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

func (c *RedisCache) SetOrderPayments(ctx context.Context, orderID string, payments []*domain.Payment, ttl int) error {
	key := "order_payments:" + orderID
	return c.Set(ctx, key, payments, ttl)
}

func (c *RedisCache) DeleteOrderPayments(ctx context.Context, orderID string) error {
	key := "order_payments:" + orderID
	return c.Delete(ctx, key)
}

// User payments cache methods
func (c *RedisCache) GetUserPayments(ctx context.Context, userID string, page, limit int) ([]*domain.Payment, error) {
	key := "user_payments:" + userID + ":" + string(page) + ":" + string(limit)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var payments []*domain.Payment
	if err := json.Unmarshal(data, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

func (c *RedisCache) SetUserPayments(ctx context.Context, userID string, page, limit int, payments []*domain.Payment, ttl int) error {
	key := "user_payments:" + userID + ":" + string(page) + ":" + string(limit)
	return c.Set(ctx, key, payments, ttl)
}

func (c *RedisCache) DeleteUserPayments(ctx context.Context, userID string) error {
	pattern := "user_payments:" + userID + ":*"
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

// MetaMask transaction cache methods
func (c *RedisCache) GetTransactionStatus(ctx context.Context, txHash string) (string, error) {
	key := "tx_status:" + txHash
	return c.client.Get(ctx, key).Result()
}

func (c *RedisCache) SetTransactionStatus(ctx context.Context, txHash string, status string, ttl int) error {
	key := "tx_status:" + txHash
	return c.client.Set(ctx, key, status, time.Duration(ttl)*time.Second).Err()
}

func (c *RedisCache) DeleteTransactionStatus(ctx context.Context, txHash string) error {
	key := "tx_status:" + txHash
	return c.Delete(ctx, key)
} 