package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/Richard-inter/game/internal/domain"
)

func (r *RedisClient) SetGachaPityStateToRedis(ctx context.Context, machineID, playerID int64, pityState *domain.GachaPityState) error {
	key := fmt.Sprintf("gacha:pity:%d:%d", machineID, playerID)

	return r.client.HSet(ctx, key, map[string]interface{}{
		"ultra_rare_pity_count": pityState.UltraRarePityCount,
		"super_rare_pity_count": pityState.SuperRarePityCount,
	}).Err()
}

func (r *RedisClient) DeleteGachaPityStateFromRedis(ctx context.Context, machineID, playerID int64) error {
	key := fmt.Sprintf("gacha:pity:%d:%d", machineID, playerID)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) GetGachaPityStateFromRedis(ctx context.Context, machineID, playerID int64) (*domain.GachaPityState, error) {
	key := fmt.Sprintf("gacha:pity:%d:%d", machineID, playerID)

	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	pityState := &domain.GachaPityState{
		PlayerID:           playerID,
		GachaMachineID:     machineID,
		UltraRarePityCount: 0,
		SuperRarePityCount: 0,
	}

	if ultraRareCount, ok := result["ultra_rare_pity_count"]; ok {
		if parsed, err := fmt.Sscanf(ultraRareCount, "%d", &pityState.UltraRarePityCount); err == nil && parsed == 1 {
		}
	}

	if superRareCount, ok := result["super_rare_pity_count"]; ok {
		if parsed, err := fmt.Sscanf(superRareCount, "%d", &pityState.SuperRarePityCount); err == nil && parsed == 1 {
		}
	}

	return pityState, nil
}

type GachaStreamMessageType string

const (
	GachaStreamEventType GachaStreamMessageType = "gacha_event"
)

type GachaStreamMessage struct {
	Type      GachaStreamMessageType   `json:"type"`
	SessionID *int64                   `json:"session_id,omitempty"`
	ItemIDs   []int64                  `json:"item_ids,omitempty"`
	Session   *domain.GachaPullSession `json:"session,omitempty"`
}

func (r *RedisClient) PublishGachaEvent(ctx context.Context, streamKey string, data ...interface{}) error {
	var message GachaStreamMessage
	if len(data) == 2 {
		if session, ok := data[0].(*domain.GachaPullSession); ok {
			switch v := data[1].(type) {
			case int64:
				message = GachaStreamMessage{
					Type:    GachaStreamEventType,
					Session: session,
					ItemIDs: []int64{v},
				}
			case []int64:
				message = GachaStreamMessage{
					Type:    GachaStreamEventType,
					Session: session,
					ItemIDs: v,
				}
			default:
				return fmt.Errorf("invalid itemID type for gacha event message, expected int64 or []int64")
			}
		} else {
			return fmt.Errorf("invalid session type for gacha event message")
		}
	} else {
		return fmt.Errorf("gacha event message requires 2 parameters: session and itemID(s)")
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal gacha stream message: %w", err)
	}

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		Values: map[string]interface{}{
			"data": string(messageJSON),
		},
	}).Err()
}
