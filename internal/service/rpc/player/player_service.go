package service

import (
	"context"
	"fmt"

	pb "github.com/Richard-inter/game/pkg/protocol/player"
)

// PlayerGRPCService implements the PlayerService gRPC service
type PlayerGRPCService struct {
	pb.UnimplementedPlayerServiceServer
}

// NewPlayerGRPCService creates a new PlayerGRPCService
func NewPlayerGRPCService() *PlayerGRPCService {
	return &PlayerGRPCService{}
}

func (s *PlayerGRPCService) GetPlayerInfo(ctx context.Context, req *pb.GetPlayerInfoReq) (*pb.GetPlayerInfoResp, error) {
	p := &pb.Player{
		PlayerID: req.PlayerID,
		UserName: "player" + fmt.Sprint(req.PlayerID),
	}
	// Implementation goes here
	fmt.Println("masuk rpc")
	return &pb.GetPlayerInfoResp{
		Player: p,
	}, nil
}
