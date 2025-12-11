package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisRepository handles caching operations with the Redis client.
type RedisRepository struct {
	Client *redis.Client
}

// NewRedisRepo creates a new Redis repository instance.
func NewRedisRepo(client *redis.Client) *RedisRepository {
	return &RedisRepository{Client: client}
}

// GetLongURL retrieves the long URL from Redis cache.
func (r *RedisRepository) GetLongURL(shortCode string) (string, error) {
	// The key in Redis will be the short code (e.g., "XyzAbC1")
	longURL, err := r.Client.Get(context.Background(), shortCode).Result()
	
	if err == redis.Nil {
		// This error means the key was not found (Cache Miss)
		return "", nil 
	}
	if err != nil {
		// This handles connection errors or other Redis issues
		return "", err
	}
	return longURL, nil
}

// SetLongURL saves the long URL to Redis cache with an expiration time (TTL).
func (r *RedisRepository) SetLongURL(shortCode string, longURL string) error {
	// Set a TTL (Time To Live) of 24 hours for the cached URL
	ttl := 24 * time.Hour 
	
	// SETEX command sets the key and the expiration time simultaneously.
	return r.Client.SetEX(context.Background(), shortCode, longURL, ttl).Err()
}