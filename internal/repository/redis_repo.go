package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	Client *redis.Client
}

func NewRedisRepo(client *redis.Client) *RedisRepository {
	return &RedisRepository{Client: client}
}

func (r *RedisRepository) GetLongURL(shortCode string) (string, error) {
	longURL, err := r.Client.Get(context.Background(), shortCode).Result()
	
	if err == redis.Nil {
		return "", nil 
	}
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func (r *RedisRepository) SetLongURL(shortCode string, longURL string) error {
	ttl := 24 * time.Hour 
	
	return r.Client.SetEX(context.Background(), shortCode, longURL, ttl).Err()
}