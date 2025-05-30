package order

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server       ServerConfig
	Services     ServicesConfig
	Auth         AuthConfig
	RateLimiting RateLimitConfig
	Redis        RedisConfig
}

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type ServicesConfig struct {
	OrderServiceURL   string
	PaymentServiceURL string
}

type AuthConfig struct {
	JWTSecret       string
	TokenExpiration time.Duration
}

type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ReadTimeout:     time.Second * 5,
			WriteTimeout:    time.Second * 10,
			ShutdownTimeout: time.Second * 30,
		},
		Services: ServicesConfig{
			OrderServiceURL:   getEnv("ORDER_SERVICE_URL", "localhost:50051"),
			PaymentServiceURL: getEnv("PAYMENT_SERVICE_URL", "localhost:50052"),
		},
		Auth: AuthConfig{
			JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
			TokenExpiration: time.Hour * 24,
		},
		RateLimiting: RateLimitConfig{
			RequestsPerMinute: getEnvAsInt("RATE_LIMIT", 60),
			BurstSize:         getEnvAsInt("RATE_LIMIT_BURST", 10),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
