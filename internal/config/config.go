package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	BaseURL       string
	DefaultTTL    int
	Port          string
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func LoadConfig() Config {
	ttlStr := GetEnv("DEFAULT_TTL_SECONDS", "86400")
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		ttl = 86400
	}

	return Config{
		RedisAddr:     GetEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: GetEnv("REDIS_PASSWORD", ""),
		BaseURL:       GetEnv("BASE_URL", "http://localhost:8080"),
		DefaultTTL:    ttl,
		Port:          GetEnv("PORT", "8080"),
	}
}

func InitRedis(cfg Config) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	var Ctx = context.Background()
	pong, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected:", pong)

	return rdb
}
