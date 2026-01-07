package service

import (
	"context"
	"fmt"

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

func (s *PlayerGRPCService) GetPlayerInfo(ctx context.Context, req *player.GetPlayerInfoReq) (*player.GetPlayerInfoResp, error) {
	// Implementation goes here
	fmt.Println("masuk rpc")
	return &player.GetPlayerInfoResp{}, nil
}
