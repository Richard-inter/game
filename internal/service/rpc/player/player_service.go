package service

import (
	"context"
	"fmt"

	"github.com/Richard-inter/game/internal/domain"
	"github.com/Richard-inter/game/internal/repository"
	pb "github.com/Richard-inter/game/pkg/protocol/player"
)

// PlayerGRPCService implements the PlayerService gRPC service
type PlayerGRPCService struct {
	pb.UnimplementedPlayerServiceServer
	repo repository.PlayerRepository
}

// NewPlayerGRPCService creates a new PlayerGRPCService
func NewPlayerGRPCService(repo repository.PlayerRepository) *PlayerGRPCService {
	return &PlayerGRPCService{
		repo: repo,
	}
}

func (s *PlayerGRPCService) GetPlayerInfo(ctx context.Context, req *pb.GetPlayerInfoReq) (*pb.GetPlayerInfoResp, error) {
	fmt.Println("masuk rpc")
	p := &pb.Player{
		PlayerID: req.PlayerID,
		UserName: "player" + fmt.Sprint(req.PlayerID),
	}
	resp, err := s.repo.GetPlayerinfo(req.PlayerID)
	if err != nil {
		return nil, err
	}
	p.PlayerID = resp.ID
	p.UserName = resp.UserName

	return &pb.GetPlayerInfoResp{
		Player: p,
	}, nil
}

func (s *PlayerGRPCService) CreatePlayer(ctx context.Context, req *pb.CreatePlayerReq) (*pb.CreatePlayerResp, error) {
	fmt.Println("masuk rpc")
	player := &domain.Player{
		UserName: req.UserName,
	}
	resp, err := s.repo.CreatePlayer(player)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePlayerResp{
		Player: &pb.Player{
			PlayerID: resp.ID,
			UserName: resp.UserName,
		},
	}, nil
}
