package service

import (
	"github.com/Richard-inter/game/pkg/protocol/player"
)

// PlayerGRPCService implements the PlayerService gRPC service
type PlayerGRPCService struct {
	player.UnimplementedPlayerServiceServer
}

// NewPlayerGRPCService creates a new PlayerGRPCService
func NewPlayerGRPCService() *PlayerGRPCService {
	return &PlayerGRPCService{}
}
