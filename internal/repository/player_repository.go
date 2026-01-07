package repository

import (
	"github.com/1nterdigital/game/internal/domain"
	"gorm.io/gorm"
)

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) domain.PlayerRepository {
	return &playerRepository{db: db}
}

func (r *playerRepository) Create(player *domain.Player) error {
	return r.db.Create(player).Error
}

func (r *playerRepository) GetByID(id string) (*domain.Player, error) {
	var player domain.Player
	err := r.db.Where("id = ?", id).First(&player).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) GetByUsername(username string) (*domain.Player, error) {
	var player domain.Player
	err := r.db.Where("username = ?", username).First(&player).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) GetByEmail(email string) (*domain.Player, error) {
	var player domain.Player
	err := r.db.Where("email = ?", email).First(&player).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) Update(player *domain.Player) error {
	return r.db.Save(player).Error
}

func (r *playerRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Player{}).Error
}

func (r *playerRepository) List(page, limit int) ([]*domain.Player, int, error) {
	var players []*domain.Player
	var total int64

	// Get total count
	if err := r.db.Model(&domain.Player{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Find(&players).Error
	if err != nil {
		return nil, 0, err
	}

	return players, int(total), nil
}

func (r *playerRepository) UpdateScore(playerID string, score int) error {
	return r.db.Model(&domain.Player{}).
		Where("id = ?", playerID).
		Update("score", score).Error
}
