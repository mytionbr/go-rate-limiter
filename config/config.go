package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RateLimitIP      int
	RateLimitToken   int
	BlockDurationSec int
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
}

func Load() *Config {
	_ = godotenv.Load()
	toInt := func(key string, def int) int {
		if v := os.Getenv(key); v != "" {
			i, err := strconv.Atoi(v)
			if err == nil {
				return i
			}
			log.Fatalf("Invalid %s: %v", key, err)
		}
		return def
	}

	return &Config{
		RateLimitIP:      toInt("RATE_LIMIT_IP", 100),
		RateLimitToken:   toInt("RATE_LIMIT_TOKEN", 100),
		BlockDurationSec: toInt("BLOCK_DURATION_SECONDS", 60),
		RedisAddr:        os.Getenv("REDIS_ADDR"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		RedisDB:          toInt("REDIS_DB", 0),
	}
}
