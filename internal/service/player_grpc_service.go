package service

import (
	"context"

	"github.com/1nterdigital/game/pkg/protocol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PlayerGRPCService implements the PlayerService gRPC service
type PlayerGRPCService struct {
	protocol.UnimplementedPlayerServiceServer
}

// NewPlayerGRPCService creates a new PlayerGRPCService
func NewPlayerGRPCService() *PlayerGRPCService {
	return &PlayerGRPCService{}
}

// CreatePlayer implements the CreatePlayer RPC method
func (_ *PlayerGRPCService) CreatePlayer(_ context.Context, req *protocol.CreatePlayerRequest) (*protocol.CreatePlayerResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// TODO: Implement actual player creation logic
	player := &protocol.Player{
		Id:       "generated-player-id",
		Username: req.Username,
		Email:    req.Email,
		Score:    0,
	}

	return &protocol.CreatePlayerResponse{
		Success: true,
		Player:  player,
	}, nil
}

// GetPlayer implements the GetPlayer RPC method
func (_ *PlayerGRPCService) GetPlayer(_ context.Context, req *protocol.GetPlayerRequest) (*protocol.GetPlayerResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "player ID is required")
	}

	// TODO: Implement actual player retrieval logic
	player := &protocol.Player{
		Id:       req.Id,
		Username: "sample_player",
		Email:    "player@example.com",
		Score:    100,
	}

	return &protocol.GetPlayerResponse{
		Success: true,
		Player:  player,
	}, nil
}

// ListPlayers implements the ListPlayers RPC method
func (_ *PlayerGRPCService) ListPlayers(_ context.Context, req *protocol.ListPlayersRequest) (*protocol.ListPlayersResponse, error) {
	// TODO: Implement actual player listing logic
	players := []*protocol.Player{
		{
			Id:       "1",
			Username: "player1",
			Email:    "player1@example.com",
			Score:    100,
		},
		{
			Id:       "2",
			Username: "player2",
			Email:    "player2@example.com",
			Score:    200,
		},
	}

	return &protocol.ListPlayersResponse{
		Success: true,
		Players: players,
		Total:   int32(len(players)),
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}
