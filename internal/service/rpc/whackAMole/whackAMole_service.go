package whackAMole

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/domain"
	"github.com/Richard-inter/game/internal/repository"
	"github.com/Richard-inter/game/pkg/logger"
	whackAMole "github.com/Richard-inter/game/pkg/protocol/whackAMole"
)

type WhackAMoleGRPCService struct {
	whackAMole.UnimplementedWhackAMoleServiceServer
	repo repository.WhackAMoleRepository
	log  *zap.SugaredLogger
}

func NewWhackAMoleGRPCService(repo repository.WhackAMoleRepository) *WhackAMoleGRPCService {
	return &WhackAMoleGRPCService{
		repo: repo,
		log:  logger.GetSugar(),
	}
}

func (s *WhackAMoleGRPCService) CreateWhackAMolePlayer(ctx context.Context, req *whackAMole.CreateWhackAMolePlayerReq) (*whackAMole.CreateWhackAMolePlayerResp, error) {
	if req.PlayerId <= 0 || req.Username == "" {
		return nil, errors.New("invalid player ID or username")
	}

	player := &domain.WhackAMolePlayer{
		Player: domain.Player{
			ID:       req.PlayerId,
			UserName: req.Username,
		},
	}

	createdPlayer, err := s.repo.CreateWhackAMolePlayer(ctx, player)
	if err != nil {
		s.log.Errorw("Failed to create Whack-A-Mole player", "playerID", req.PlayerId, "error", err)
		return nil, err
	}

	return &whackAMole.CreateWhackAMolePlayerResp{
		Player: &whackAMole.WhackAMolePlayer{
			PlayerId: createdPlayer.Player.ID,
			Username: createdPlayer.Player.UserName,
		},
	}, nil
}

func (s *WhackAMoleGRPCService) GetPlayerInfo(ctx context.Context, req *whackAMole.GetPlayerInfoReq) (*whackAMole.GetPlayerInfoResp, error) {
	if req.PlayerId <= 0 {
		return nil, errors.New("invalid player ID")
	}

	player, err := s.repo.GetWhackAMolePlayerInfo(ctx, req.PlayerId)
	if err != nil {
		s.log.Errorw("Failed to get player info", "playerID", req.PlayerId, "error", err)
		return nil, err
	}

	return &whackAMole.GetPlayerInfoResp{
		Player: &whackAMole.WhackAMolePlayer{
			PlayerId: player.Player.ID,
			Username: player.Player.UserName,
		},
	}, nil
}

func (s *WhackAMoleGRPCService) GetLeaderboard(ctx context.Context, req *whackAMole.GetLeaderboardReq) (*whackAMole.GetLeaderboardResp, error) {
	leaderboard, err := s.repo.GetLeaderboard(ctx, req.Limit)
	if err != nil {
		s.log.Errorw("Failed to get leaderboard", "limit", req.Limit, "error", err)
		return nil, err
	}

	var respLeaderboard []*whackAMole.LeaderBoard
	for _, entry := range leaderboard {
		respLeaderboard = append(respLeaderboard, &whackAMole.LeaderBoard{
			PlayerId: entry.PlayerID,
			Username: entry.Username,
			Score:    entry.Score,
			Rank:     entry.Rank,
		})
	}

	return &whackAMole.GetLeaderboardResp{
		Leaderboard: respLeaderboard,
	}, nil
}

func (s *WhackAMoleGRPCService) GetMoleWeightConfig(ctx context.Context, req *whackAMole.GetMoleWeightConfigReq) (*whackAMole.GetMoleWeightConfigResp, error) {
	if req.Id < 0 {
		return nil, errors.New("invalid mole ID")
	}

	config, err := s.repo.GetMoleWeightConfig(ctx, req.Id)
	if err != nil {
		s.log.Errorw("Failed to get mole weight config", "moleID", req.Id, "error", err)
		return nil, err
	}

	var respConfigs []*whackAMole.MoleWeightConfig
	for _, c := range config {
		respConfigs = append(respConfigs, &whackAMole.MoleWeightConfig{
			Id:       c.ID,
			MoleType: c.MoleType,
			Weight:   c.Weight,
		})
	}

	return &whackAMole.GetMoleWeightConfigResp{
		Config: respConfigs,
	}, nil
}

func (s *WhackAMoleGRPCService) UpdateScore(ctx context.Context, req *whackAMole.UpdateScoreReq) (*whackAMole.UpdateScoreResp, error) {
	// Get current score first
	currentRank, err := s.repo.GetPlayerRank(ctx, req.PlayerId)
	if err != nil {
		// Player might not exist in leaderboard yet, start from 0
		currentRank = &domain.LeaderBoard{Score: 0}
	}

	if currentRank.Score > req.Score {
		return nil, errors.New("new score must be greater than current high score")
	}

	err = s.repo.UpdatePlayerScore(ctx, req.PlayerId, req.Score)
	if err != nil {
		s.log.Errorw("Failed to update player score", "playerID", req.PlayerId, "score", req.Score, "error", err)
		return nil, err
	}

	return &whackAMole.UpdateScoreResp{
		Success: true,
	}, nil
}

func (s *WhackAMoleGRPCService) CreateMoleWeightConfig(ctx context.Context, req *whackAMole.CreateMoleWeightConfigReq) (*whackAMole.CreateMoleWeightConfigResp, error) {
	config := &domain.MoleWeightConfig{
		MoleType: req.MoleType,
		Weight:   req.Weight,
	}

	createdConfig, err := s.repo.CreateMoleWeightConfig(ctx, config)
	if err != nil {
		s.log.Errorw("Failed to create mole weight config", "moleType", req.MoleType, "weight", req.Weight, "error", err)
		return nil, err
	}

	return &whackAMole.CreateMoleWeightConfigResp{
		Config: &whackAMole.MoleWeightConfig{
			Id:       createdConfig.ID,
			MoleType: createdConfig.MoleType,
			Weight:   createdConfig.Weight,
		},
	}, nil
}

func (s *WhackAMoleGRPCService) UpdateMoleWeightConfig(ctx context.Context, req *whackAMole.UpdateMoleWeightConfigReq) (*whackAMole.UpdateMoleWeightConfigResp, error) {
	config := &domain.MoleWeightConfig{
		ID:       req.Id,
		MoleType: req.MoleType,
		Weight:   req.Weight,
	}

	updatedConfig, err := s.repo.UpdateMoleWeightConfig(ctx, config)
	if err != nil {
		s.log.Errorw("Failed to update mole weight config", "moleID", req.Id, "moleType", req.MoleType, "weight", req.Weight, "error", err)
		return nil, err
	}

	return &whackAMole.UpdateMoleWeightConfigResp{
		Config: &whackAMole.MoleWeightConfig{
			Id:       updatedConfig.ID,
			MoleType: updatedConfig.MoleType,
			Weight:   updatedConfig.Weight,
		},
	}, nil
}
