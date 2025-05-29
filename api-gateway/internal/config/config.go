package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPPort       int
	UserServiceURL string
	RedisAddress   string
	RedisPassword  string
	RedisDB        int
	NATSAddress    string
}

func Load() *Config {
	return &Config{
		HTTPPort:       getEnvAsInt("HTTP_PORT", 8080),
		UserServiceURL: getEnv("USER_SERVICE_URL", "localhost:50051"),
		RedisAddress:   getEnv("REDIS_ADDRESS", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getEnvAsInt("REDIS_DB", 0),
		NATSAddress:    getEnv("NATS_ADDRESS", "nats://localhost:4222"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, err := strconv.Atoi(getEnv(key, "")); err == nil {
		return value
	}
	return defaultValue
}
