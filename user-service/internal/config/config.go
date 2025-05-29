package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort   int
	MongoURI     string
	DBName       string
	NATSAddress  string
	RedisAddress string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	GRPCPort     int
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:   50053,
		MongoURI:     "mongodb://localhost:27017",
		DBName:       "user_service",
		NATSAddress:  getEnv("NATS_ADDRESS", "nats://localhost:4222"),
		RedisAddress: getEnv("REDIS_ADDRESS", "localhost:6379"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.mail.ru"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 465),
		SMTPUsername: getEnv("SMTP_USERNAME", "aytzhanovk@internet.ru"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "Q1mfUaA5BcuryV4Fo3Gq"),
		GRPCPort:     getEnvAsInt("GRPC_PORT", 50053),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
