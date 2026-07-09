package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cms-blog-api/config"
)

const (
	// CacheTTL defines the 5-minute cache expiration as required.
	CacheTTL = 5 * time.Minute

	KeyAllBlogs       = "blogs:all"
	KeyBlogByIDPrefix = "blogs:id:"
	KeyCategoryPrefix = "blogs:category:"
)

// SetCache serializes the value to JSON and stores it in Redis with the given TTL.
func SetCache(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if config.RedisClient == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(ctx, key, data, ttl).Err()
}

// GetCache retrieves JSON data from Redis and unmarshals it into dest.
func GetCache(ctx context.Context, key string, dest interface{}) error {
	if config.RedisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}

	val, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// InvalidateBlogCaches clears all blog-related cache entries in Redis.
// Called whenever a blog is created, updated, or deleted.
func InvalidateBlogCaches(ctx context.Context, blogID uint, category string) {
	if config.RedisClient == nil {
		return
	}

	keysToDelete := []string{
		KeyAllBlogs,
	}

	if blogID > 0 {
		keysToDelete = append(keysToDelete, fmt.Sprintf("%s%d", KeyBlogByIDPrefix, blogID))
	}
	if category != "" {
		keysToDelete = append(keysToDelete, fmt.Sprintf("%s%s", KeyCategoryPrefix, category))
	}

	// Delete explicit keys
	err := config.RedisClient.Del(ctx, keysToDelete...).Err()
	if err != nil {
		log.Printf("Error deleting explicit cache keys: %v", err)
	}

	// Also clear any matching category or list keys to be safe
	iter := config.RedisClient.Scan(ctx, 0, "blogs:*", 100).Iterator()
	for iter.Next(ctx) {
		config.RedisClient.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		log.Printf("Error scanning cache keys for invalidation: %v", err)
	} else {
		log.Println("Blog caches invalidated successfully")
	}
}
