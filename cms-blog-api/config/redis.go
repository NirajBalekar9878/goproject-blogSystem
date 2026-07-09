package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// ConnectRedis initializes the Redis client and verifies connection with Ping.
func ConnectRedis() *redis.Client {
	redisHost := getEnv("REDIS_HOST", "127.0.0.1")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPass := getEnv("REDIS_PASSWORD", "")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Could not connect to Redis: %v. Caching will be skipped if Redis is down.", err)
	} else {
		log.Println("Redis connected successfully")
	}

	return RedisClient
}
