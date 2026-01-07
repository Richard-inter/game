package domain

import (
	"time"
)

type Player struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Score     int       `json:"score" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PlayerRepository interface {
	Create(player *Player) error
	GetByID(id string) (*Player, error)
	GetByUsername(username string) (*Player, error)
	GetByEmail(email string) (*Player, error)
	Update(player *Player) error
	Delete(id string) error
	List(page, limit int) ([]*Player, int, error)
	UpdateScore(playerID string, score int) error
}

type PlayerService interface {
	CreatePlayer(username, email string) (*Player, error)
	GetPlayer(id string) (*Player, error)
	GetPlayerByUsername(username string) (*Player, error)
	UpdatePlayer(id string, updates map[string]interface{}) (*Player, error)
	DeletePlayer(id string) error
	ListPlayers(page, limit int) ([]*Player, int, error)
	AddScore(playerID string, points int) error
	SubtractScore(playerID string, points int) error
}
