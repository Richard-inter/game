package domain

import (
	"time"
)

type GameStatus string

const (
	GameStatusActive      GameStatus = "active"
	GameStatusInactive    GameStatus = "inactive"
	GameStatusMaintenance GameStatus = "maintenance"
)

type Game struct {
	ID             string     `json:"id" gorm:"primaryKey"`
	Name           string     `json:"name" gorm:"not null"`
	Description    string     `json:"description"`
	Status         GameStatus `json:"status" gorm:"default:'inactive'"`
	MaxPlayers     int        `json:"max_players" gorm:"default:10"`
	CurrentPlayers int        `json:"current_players" gorm:"default:0"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type GameRepository interface {
	Create(game *Game) error
	GetByID(id string) (*Game, error)
	Update(game *Game) error
	Delete(id string) error
	List(page, limit int, status GameStatus) ([]*Game, int, error)
	JoinGame(gameID, playerID string) error
	LeaveGame(gameID, playerID string) error
}

type GameService interface {
	CreateGame(name, description string, maxPlayers int) (*Game, error)
	GetGame(id string) (*Game, error)
	UpdateGame(id string, updates map[string]interface{}) (*Game, error)
	DeleteGame(id string) error
	ListGames(page, limit int, status GameStatus) ([]*Game, int, error)
	JoinGame(gameID, playerID string) (*Game, error)
	LeaveGame(gameID, playerID string) error
}
