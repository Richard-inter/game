package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/cache"
	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/domain"
	"github.com/Richard-inter/game/internal/repository"
	g "github.com/Richard-inter/game/internal/service/rpc/gachaMachine"
	"github.com/Richard-inter/game/pkg/logger"
)

type GachaStreamConsumer struct {
	repo    repository.GachaMachineRepository
	redis   *cache.RedisClient
	config  *config.StreamConsumerConfig
	service *g.GachaMachineGRPCService // Add service reference
	log     *zap.SugaredLogger
}

func NewGachaStreamConsumer(repo repository.GachaMachineRepository, redis *cache.RedisClient, cfg *config.StreamConsumerConfig, service *g.GachaMachineGRPCService) *GachaStreamConsumer {
	return &GachaStreamConsumer{
		repo:    repo,
		redis:   redis,
		config:  cfg,
		service: service,
		log:     logger.GetSugar(),
	}
}

func (g *GachaStreamConsumer) Start(ctx context.Context) {
	go func() {
		if err := g.StartEventConsumer(ctx); err != nil {
			g.log.Errorf("Event consumer error: %v", err)
		}
	}()

	g.log.Info("Gacha stream consumer started")
}

func (g *GachaStreamConsumer) StartEventConsumer(ctx context.Context) error {
	blockTimeout, err := time.ParseDuration(g.config.BlockTimeout)
	if err != nil {
		blockTimeout = 5 * time.Second
	}

	if err := g.redis.GetClient().XGroupCreateMkStream(ctx, g.config.StreamKey, g.config.ConsumerGroup, "0").Err(); err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			return fmt.Errorf("failed to create consumer group: %w", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			g.log.Info("Event consumer stopped")
			return ctx.Err()
		default:
			messages, err := g.redis.GetClient().XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    g.config.ConsumerGroup,
				Consumer: g.config.ConsumerName,
				Streams:  []string{g.config.StreamKey, ">"},
				Count:    int64(g.config.BatchSize),
				Block:    blockTimeout,
			}).Result()

			if err != nil && err != redis.Nil {
				g.log.Errorf("Failed to read from stream: %v", err)
				time.Sleep(time.Second)
				continue
			}

			for _, stream := range messages {
				for _, message := range stream.Messages {
					if _, err := g.parseHistoryMessage(ctx, message); err != nil {
						g.log.Errorf("Failed to process event message: %v", err)
						continue
					}
					g.redis.GetClient().XAck(ctx, g.config.StreamKey, g.config.ConsumerGroup, message.ID)
				}
			}
		}
	}
}

func (g *GachaStreamConsumer) parseHistoryMessage(ctx context.Context, message redis.XMessage) (*domain.GachaPullHistory, error) {
	data, ok := message.Values["data"]
	if !ok {
		return nil, fmt.Errorf("message missing data field")
	}

	dataStr, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("data field is not a string")
	}

	var streamMessage cache.GachaStreamMessage
	if err := json.Unmarshal([]byte(dataStr), &streamMessage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stream message: %w", err)
	}

	if streamMessage.Type != cache.GachaStreamEventType {
		return nil, fmt.Errorf("expected gacha event message type, got: %s", streamMessage.Type)
	}

	var session *domain.GachaPullSession
	var itemIDs []int64

	if streamMessage.Session != nil && len(streamMessage.ItemIDs) > 0 {
		session = streamMessage.Session
		itemIDs = streamMessage.ItemIDs
	} else {
		return nil, fmt.Errorf("history message missing required session and item data")
	}

	if err := g.service.AddGameToHistory(ctx, *session, itemIDs); err != nil {
		return nil, fmt.Errorf("failed to add game to history for item %v: %w", itemIDs, err)
	}

	history := domain.GachaPullHistory{
		GachaPullSessionID: session.ID,
		ItemID:             itemIDs[len(itemIDs)-1],
	}

	return &history, nil
}
