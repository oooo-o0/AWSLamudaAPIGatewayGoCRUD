package config

import (
	"os"
)

type Config struct {
	Env            string
	DynamoRegion   string
	DynamoEndpoint string
	LogLevel       string
	ClaimTable     string
}

func Load() (*Config, error) {
	cfg := &Config{
		Env:            getEnv("ENV", "dev"),
		DynamoRegion:   getEnv("DYNAMO_REGION", "ap-northeast-1"),
		DynamoEndpoint: getEnv("DYNAMO_ENDPOINT", ""), // For LocalStack
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		ClaimTable:     getEnv("CLAIM_TABLE", "Claim"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

