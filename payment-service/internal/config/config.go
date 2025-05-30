package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            string
	RedisURL        string
	RedisPassword   string
	RedisDB         int
	MongoURI        string
	MongoDB         string
	NatsURL         string
	JWTSecret       string
	RateLimit       int
	RateLimitBurst  int
	StripeSecretKey string
	EthereumRPC     string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "50052"),
		RedisURL:        getEnv("REDIS_URL", "redis:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnvAsInt("REDIS_DB", 0),
		MongoURI:        getEnv("MONGO_URI", "mongodb://mongodb:27017"),
		MongoDB:         getEnv("MONGO_DB", "payments"),
		NatsURL:         getEnv("NATS_URL", "nats://nats:4222"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		RateLimit:       getEnvAsInt("RATE_LIMIT", 60),
		RateLimitBurst:  getEnvAsInt("RATE_LIMIT_BURST", 10),
		StripeSecretKey: getEnv("STRIPE_SECRET_KEY", ""),
		EthereumRPC:     getEnv("ETHEREUM_RPC", "https://mainnet.infura.io/v3/your-project-id"),
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