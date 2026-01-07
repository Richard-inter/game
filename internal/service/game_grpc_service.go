package service

import (
	"context"

	"github.com/1nterdigital/game/pkg/protocol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GameGRPCService implements the GameService gRPC service
type GameGRPCService struct {
	protocol.UnimplementedGameServiceServer
}

// NewGameGRPCService creates a new GameGRPCService
func NewGameGRPCService() *GameGRPCService {
	return &GameGRPCService{}
}

// CreateGame implements the CreateGame RPC method
func (_ *GameGRPCService) CreateGame(_ context.Context, req *protocol.CreateGameRequest) (*protocol.CreateGameResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "game name is required")
	}

	// TODO: Implement actual game creation logic
	game := &protocol.Game{
		Id:             "generated-game-id",
		Name:           req.Name,
		Description:    req.Description,
		Status:         "active",
		MaxPlayers:     req.MaxPlayers,
		CurrentPlayers: 0,
	}

	return &protocol.CreateGameResponse{
		Success: true,
		Game:    game,
	}, nil
}

// GetGame implements the GetGame RPC method
func (_ *GameGRPCService) GetGame(_ context.Context, req *protocol.GetGameRequest) (*protocol.GetGameResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "game ID is required")
	}

	// TODO: Implement actual game retrieval logic
	game := &protocol.Game{
		Id:             req.Id,
		Name:           "Sample Game",
		Description:    "A sample game",
		Status:         "active",
		MaxPlayers:     10,
		CurrentPlayers: 2,
	}

	return &protocol.GetGameResponse{
		Success: true,
		Game:    game,
	}, nil
}

// ListGames implements the ListGames RPC method
func (_ *GameGRPCService) ListGames(_ context.Context, req *protocol.ListGamesRequest) (*protocol.ListGamesResponse, error) {
	// TODO: Implement actual game listing logic
	games := []*protocol.Game{
		{
			Id:             "1",
			Name:           "Game 1",
			Description:    "First game",
			Status:         "active",
			MaxPlayers:     10,
			CurrentPlayers: 2,
		},
		{
			Id:             "2",
			Name:           "Game 2",
			Description:    "Second game",
			Status:         "active",
			MaxPlayers:     8,
			CurrentPlayers: 4,
		},
	}

	return &protocol.ListGamesResponse{
		Success: true,
		Games:   games,
		Total:   int32(len(games)),
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}

// JoinGame implements JoinGame RPC method
func (_ *GameGRPCService) JoinGame(_ context.Context, req *protocol.JoinGameRequest) (*protocol.JoinGameResponse, error) {
	if req.GameId == "" {
		return nil, status.Error(codes.InvalidArgument, "game ID is required")
	}
	if req.PlayerId == "" {
		return nil, status.Error(codes.InvalidArgument, "player ID is required")
	}

	// TODO: Implement actual game joining logic
	game := &protocol.Game{
		Id:             req.GameId,
		Name:           "Joined Game",
		Description:    "Game with new player",
		Status:         "active",
		MaxPlayers:     10,
		CurrentPlayers: 3, // Incremented
	}

	return &protocol.JoinGameResponse{
		Success: true,
		Game:    game,
	}, nil
}
