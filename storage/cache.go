package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis connects to the Redis server
func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Println("Connected to Redis")
	return nil
}

// SetCache stores a key-value pair with TTL in Redis
func SetCache(key string, value string, ttl time.Duration) error {
	return RedisClient.Set(ctx, key, value, ttl).Err()
}

// GetCache retrieves the value for a given key from Redis
func GetCache(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}
