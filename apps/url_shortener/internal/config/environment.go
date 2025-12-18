package config

import (
	"fmt"
	"os"
	"strconv"
)

type Environment struct {
	ServerPort     string
	MongoURI       string
	RedisHost      string
	RedisPort      string
	RedisPass      string
	CacheTTL       int
	CounterOffset  int64
	Base62Alphabet string
}

func LoadEnvironment() (*Environment, error) {
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	serverPort := getEnv("SERVER_PORT", "8080")

	if os.Getenv("REDIS_PASS") == "" && os.Getenv("GO_ENV") == "production" {
		return nil, fmt.Errorf("REDIS_PASS is required in production")
	}

	cacheTTL, err := strconv.Atoi(getEnv("CACHE_TTL", "3600"))
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_TTL: %w", err)
	}

	counterOffset, err := strconv.ParseInt(getEnv("COUNTER_OFFSET", "0"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid COUNTER_OFFSET: %w", err)
	}

	alphabet := getEnv("BASE62_ALPHABET", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if len(alphabet) != 62 {
		return nil, fmt.Errorf("alphabet must have exactly 62 characters, got %d", len(alphabet))
	}

	return &Environment{
		ServerPort:     serverPort,
		MongoURI:       mongoURI,
		RedisHost:      redisHost,
		RedisPort:      redisPort,
		RedisPass:      os.Getenv("REDIS_PASS"),
		CacheTTL:       cacheTTL,
		CounterOffset:  counterOffset,
		Base62Alphabet: alphabet,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
