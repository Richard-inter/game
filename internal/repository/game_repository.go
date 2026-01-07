package repository

import (
	"github.com/Richard-inter/game/internal/domain"
	"gorm.io/gorm"
)

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) domain.GameRepository {
	return &gameRepository{db: db}
}

func (r *gameRepository) Create(game *domain.Game) error {
	return r.db.Create(game).Error
}

func (r *gameRepository) GetByID(id string) (*domain.Game, error) {
	var game domain.Game
	err := r.db.Where("id = ?", id).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *gameRepository) Update(game *domain.Game) error {
	return r.db.Save(game).Error
}

func (r *gameRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Game{}).Error
}

func (r *gameRepository) List(page, limit int, status domain.GameStatus) ([]*domain.Game, int, error) {
	var games []*domain.Game
	var total int64

	query := r.db.Model(&domain.Game{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&games).Error
	if err != nil {
		return nil, 0, err
	}

	return games, int(total), nil
}

func (r *gameRepository) JoinGame(gameID, playerID string) error {
	// Check if player already joined
	var count int64
	r.db.Table("game_players").
		Where("game_id = ? AND player_id = ?", gameID, playerID).
		Count(&count)

	if count > 0 {
		return nil // Already joined
	}

	// Add player to game
	return r.db.Exec("INSERT INTO game_players (id, game_id, player_id) VALUES (?, ?, ?)",
		generateUUID(), gameID, playerID).Error
}

func (r *gameRepository) LeaveGame(gameID, playerID string) error {
	return r.db.Where("game_id = ? AND player_id = ?", gameID, playerID).
		Delete(&struct{}{}).Error
}

// Helper function to generate UUID (placeholder)
func generateUUID() string {
	return "generated-uuid"
}
