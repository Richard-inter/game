package service

import (
	"errors"
	"fmt"

	"github.com/1nterdigital/game/internal/domain"
	"github.com/google/uuid"
)

type gameService struct {
	gameRepo domain.GameRepository
}

func NewGameService(gameRepo domain.GameRepository) domain.GameService {
	return &gameService{gameRepo: gameRepo}
}

func (s *gameService) CreateGame(name, description string, maxPlayers int) (*domain.Game, error) {
	if name == "" {
		return nil, errors.New("game name is required")
	}
	if maxPlayers <= 0 {
		maxPlayers = 10 // Default
	}

	game := &domain.Game{
		ID:             uuid.New().String(),
		Name:           name,
		Description:    description,
		Status:         domain.GameStatusInactive,
		MaxPlayers:     maxPlayers,
		CurrentPlayers: 0,
	}

	if err := s.gameRepo.Create(game); err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	return game, nil
}

func (s *gameService) GetGame(id string) (*domain.Game, error) {
	if id == "" {
		return nil, errors.New("game ID is required")
	}

	game, err := s.gameRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return game, nil
}

func (s *gameService) UpdateGame(id string, updates map[string]interface{}) (*domain.Game, error) {
	if id == "" {
		return nil, errors.New("game ID is required")
	}

	game, err := s.gameRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		game.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		game.Description = description
	}
	if status, ok := updates["status"].(domain.GameStatus); ok {
		game.Status = status
	}
	if maxPlayers, ok := updates["max_players"].(int); ok {
		game.MaxPlayers = maxPlayers
	}

	if err := s.gameRepo.Update(game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return game, nil
}

func (s *gameService) DeleteGame(id string) error {
	if id == "" {
		return errors.New("game ID is required")
	}

	if err := s.gameRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}

	return nil
}

func (s *gameService) ListGames(page, limit int, status domain.GameStatus) ([]*domain.Game, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	games, total, err := s.gameRepo.List(page, limit, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list games: %w", err)
	}

	return games, total, nil
}

func (s *gameService) JoinGame(gameID, playerID string) (*domain.Game, error) {
	if gameID == "" {
		return nil, errors.New("game ID is required")
	}
	if playerID == "" {
		return nil, errors.New("player ID is required")
	}

	game, err := s.gameRepo.GetByID(gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	if game.Status != domain.GameStatusActive {
		return nil, errors.New("game is not active")
	}

	if game.CurrentPlayers >= game.MaxPlayers {
		return nil, errors.New("game is full")
	}

	// Join the game
	if err := s.gameRepo.JoinGame(gameID, playerID); err != nil {
		return nil, fmt.Errorf("failed to join game: %w", err)
	}

	// Update current players count
	game.CurrentPlayers++
	if err := s.gameRepo.Update(game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return game, nil
}

func (s *gameService) LeaveGame(gameID, playerID string) error {
	if gameID == "" {
		return errors.New("game ID is required")
	}
	if playerID == "" {
		return errors.New("player ID is required")
	}

	game, err := s.gameRepo.GetByID(gameID)
	if err != nil {
		return fmt.Errorf("failed to get game: %w", err)
	}

	// Leave the game
	if err := s.gameRepo.LeaveGame(gameID, playerID); err != nil {
		return fmt.Errorf("failed to leave game: %w", err)
	}

	// Update current players count
	if game.CurrentPlayers > 0 {
		game.CurrentPlayers--
		if err := s.gameRepo.Update(game); err != nil {
			return fmt.Errorf("failed to update game: %w", err)
		}
	}

	return nil
}
