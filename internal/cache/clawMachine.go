package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// StoreGameResults stores the game results in Redis with expiration
func (r *RedisClient) StoreGameResults(ctx context.Context, gameID int64, results any) error {
	key := fmt.Sprintf("%s:%d", GameResultsKeyPrefix, gameID)

	data, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal game results: %w", err)
	}

	// Store with 5 minutes expiration
	return r.client.Set(ctx, key, data, time.Minute*5).Err()
}

// GetGameResults retrieves the game results from Redis
func (r *RedisClient) GetGameResults(ctx context.Context, gameID int64, dest any) error {
	key := fmt.Sprintf("%s:%d", GameResultsKeyPrefix, gameID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("game results not found for game ID: %d", gameID)
		}
		return fmt.Errorf("failed to get game results: %w", err)
	}

	return json.Unmarshal([]byte(data), dest)
}

// DeleteGameResults removes the game results from Redis
func (r *RedisClient) DeleteGameResults(ctx context.Context, gameID int64) error {
	key := fmt.Sprintf("%s:%d", GameResultsKeyPrefix, gameID)
	return r.client.Del(ctx, key).Err()
}
