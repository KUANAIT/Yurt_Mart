package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"user-service/internal/core/domain"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(client *redis.Client, ttl time.Duration) *RedisCache {
	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *RedisCache) GetUser(ctx context.Context, id string) (*domain.User, error) {
	data, err := c.client.Get(ctx, "user:"+id).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *RedisCache) SetUser(ctx context.Context, user *domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, "user:"+user.ID, data, c.ttl).Err()
}

func (c *RedisCache) DeleteUser(ctx context.Context, id string) error {
	return c.client.Del(ctx, "user:"+id).Err()
}
