package cache

import (
	"github.com/redis/go-redis/v9"
)

const (
	// GameResultsKeyPrefix is the prefix for game results keys in Redis
	GameResultsKeyPrefix = "game_results"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, password string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &RedisClient{
		client: rdb,
	}
}

// GetClient returns the underlying redis.Client for advanced operations
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}
