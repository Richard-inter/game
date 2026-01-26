package repository

import (
	"context"
	"fmt"

	"github.com/Richard-inter/game/internal/domain"
	"gorm.io/gorm"
)

type whackAMoleRepository struct {
	db *gorm.DB
}

type WhackAMoleRepository interface {
	// Player operations
	CreateWhackAMolePlayer(ctx context.Context, player *domain.WhackAMolePlayer) (*domain.WhackAMolePlayer, error)
	GetWhackAMolePlayerInfo(ctx context.Context, playerID int64) (*domain.WhackAMolePlayer, error)

	// Leaderboard operations
	GetLeaderboard(ctx context.Context, limit int32) ([]*domain.LeaderBoard, error)
	UpdatePlayerScore(ctx context.Context, playerID int64, score int64) error
	GetPlayerRank(ctx context.Context, playerID int64) (*domain.LeaderBoard, error)
	RecalculateLeaderboard(ctx context.Context) error

	// Mole weight config operations
	GetMoleWeightConfig(ctx context.Context, id int64) ([]domain.MoleWeightConfig, error)
	CreateMoleWeightConfig(ctx context.Context, config *domain.MoleWeightConfig) (*domain.MoleWeightConfig, error)
	UpdateMoleWeightConfig(ctx context.Context, config *domain.MoleWeightConfig) (*domain.MoleWeightConfig, error)
}

func NewWhackAMoleRepository(db *gorm.DB) WhackAMoleRepository {
	return &whackAMoleRepository{
		db: db,
	}
}

func (r *whackAMoleRepository) CreateWhackAMolePlayer(ctx context.Context, player *domain.WhackAMolePlayer) (*domain.WhackAMolePlayer, error) {
	err := r.db.WithContext(ctx).Create(player).Error
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (r *whackAMoleRepository) GetWhackAMolePlayerInfo(ctx context.Context, playerID int64) (*domain.WhackAMolePlayer, error) {
	var player domain.WhackAMolePlayer
	err := r.db.WithContext(ctx).Where("player_id = ?", playerID).First(&player).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *whackAMoleRepository) GetLeaderboard(
	ctx context.Context,
	limit int32,
) ([]*domain.LeaderBoard, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than zero")
	}
	if limit > 100 {
		limit = 100
	}

	var leaderboard []*domain.LeaderBoard
	err := r.db.WithContext(ctx).
		Where("`rank` > 0").
		Order("`rank` ASC").
		Limit(int(limit)).
		Find(&leaderboard).Error

	if err != nil {
		return nil, err
	}

	return leaderboard, nil
}

func (r *whackAMoleRepository) UpdatePlayerScore(ctx context.Context, playerID int64, score int64) error {
	// First try to update existing player
	result := r.db.WithContext(ctx).Model(&domain.LeaderBoard{}).
		Where("player_id = ?", playerID).
		Update("score", score)

	if result.Error != nil {
		return result.Error
	}

	// If no rows were affected, create new entry
	if result.RowsAffected == 0 {
		// Get player info first
		player, err := r.GetWhackAMolePlayerInfo(ctx, playerID)
		if err != nil {
			return fmt.Errorf("player not found: %w", err)
		}

		leaderboardEntry := &domain.LeaderBoard{
			PlayerID: playerID,
			Username: player.Player.UserName,
			Score:    score,
			Rank:     0, // Will be calculated when fetching leaderboard
		}

		return r.db.WithContext(ctx).Create(leaderboardEntry).Error
	}

	return nil
}

func (r *whackAMoleRepository) GetPlayerRank(ctx context.Context, playerID int64) (*domain.LeaderBoard, error) {
	var leaderboard domain.LeaderBoard
	err := r.db.WithContext(ctx).Where("player_id = ?", playerID).First(&leaderboard).Error
	if err != nil {
		return nil, err
	}

	// Calculate current rank
	var rank int64
	err = r.db.WithContext(ctx).Model(&domain.LeaderBoard{}).
		Where("score > ? OR (score = ? AND player_id < ?)", leaderboard.Score, leaderboard.Score, playerID).
		Count(&rank).Error
	if err != nil {
		return nil, err
	}

	leaderboard.Rank = int32(rank + 1)
	return &leaderboard, nil
}

func (r *whackAMoleRepository) GetMoleWeightConfig(
	ctx context.Context,
	id int64,
) ([]domain.MoleWeightConfig, error) {

	var configs []domain.MoleWeightConfig

	query := r.db.WithContext(ctx)

	if id != 0 {
		query = query.Where("id = ?", id)
	}

	if err := query.Find(&configs).Error; err != nil {
		return nil, err
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("no mole weight configs found")
	}

	return configs, nil
}

func (r *whackAMoleRepository) CreateMoleWeightConfig(ctx context.Context, config *domain.MoleWeightConfig) (*domain.MoleWeightConfig, error) {
	err := r.db.WithContext(ctx).Create(config).Error
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (r *whackAMoleRepository) UpdateMoleWeightConfig(ctx context.Context, config *domain.MoleWeightConfig) (*domain.MoleWeightConfig, error) {
	err := r.db.WithContext(ctx).Model(&domain.MoleWeightConfig{}).
		Where("id = ?", config.ID).
		Updates(map[string]interface{}{
			"mole_type": config.MoleType,
			"weight":    config.Weight,
		}).Error
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (r *whackAMoleRepository) RecalculateLeaderboard(ctx context.Context) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Reset all ranks first
		if err := tx.Exec(`UPDATE whackAMole_leaderboard SET ` + "`rank`" + ` = 0`).Error; err != nil {
			return err
		}

		// 2. Set ranks for top 100 players
		if err := tx.Exec(`
			UPDATE whackAMole_leaderboard lb
			INNER JOIN (
				SELECT player_id, 
				       ROW_NUMBER() OVER (ORDER BY score DESC, player_id ASC) AS new_rank
				FROM whackAMole_leaderboard
			) ranked ON lb.player_id = ranked.player_id
			SET lb.` + "`rank`" + ` = ranked.new_rank
			WHERE ranked.new_rank <= 100
		`).Error; err != nil {
			return err
		}

		return nil
	})
}
