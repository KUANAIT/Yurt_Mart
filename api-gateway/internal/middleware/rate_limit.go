package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"api-gateway/internal/config/order"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	redis  *redis.Client
	config *order.RateLimitConfig
}

func NewRateLimiter(redisConfig *order.RedisConfig, rateLimitConfig *order.RateLimitConfig) *RateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.URL,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	return &RateLimiter{
		redis:  client,
		config: rateLimitConfig,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()
		key := "rate_limit:" + clientIP

		ctx := context.Background()
		now := time.Now().Unix()
		windowStart := now - 60 // 1-minute window

		// Clean up old requests
		rl.redis.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

		// Count requests in the current window
		recentRequests, err := rl.redis.ZCount(ctx, key, strconv.FormatInt(windowStart, 10), strconv.FormatInt(now, 10)).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "rate limit check failed"})
			return
		}

		// Check if rate limit is exceeded
		if recentRequests >= int64(rl.config.RequestsPerMinute) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"limit": rl.config.RequestsPerMinute,
				"reset": 60 - (now % 60), // Seconds until the next window
			})
			return
		}

		// Record this request
		member := redis.Z{
			Score:  float64(now),
			Member: now,
		}
		err = rl.redis.ZAdd(ctx, key, &member).Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "rate limit recording failed"})
			return
		}

		// Set key expiration
		rl.redis.Expire(ctx, key, time.Minute*2)

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(int64(rl.config.RequestsPerMinute)-recentRequests-1, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(60-(now%60), 10))

		c.Next()
	}
}

func (rl *RateLimiter) Close() error {
	return rl.redis.Close()
}
