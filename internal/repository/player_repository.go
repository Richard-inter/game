package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Richard-inter/game/internal/domain"
)

type playerRepository struct {
	db *gorm.DB
}

type PlayerRepository interface {
	GetPlayerinfo(id int64) (*domain.Player, error)
	CreatePlayer(player *domain.Player) (*domain.Player, error)
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{db: db}
}

func validateUsernameUnique(db *gorm.DB, username string) error {
	var count int64
	err := db.Model(&domain.Player{}).Where("user_name = ?", username).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check username uniqueness: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("username '%s' already exists", username)
	}
	return nil
}

func (r *playerRepository) GetPlayerinfo(id int64) (*domain.Player, error) {
	var player domain.Player
	err := r.db.Where("id = ?", id).First(&player).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) CreatePlayer(player *domain.Player) (*domain.Player, error) {
	if err := validateUsernameUnique(r.db, player.UserName); err != nil {
		return nil, err
	}

	err := r.db.Create(player).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}
	return player, nil
}
