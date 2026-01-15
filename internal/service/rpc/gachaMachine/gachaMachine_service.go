package gachaMachine

import (
	"context"

	"github.com/Richard-inter/game/internal/cache"
	"github.com/Richard-inter/game/internal/repository"
	"github.com/Richard-inter/game/pkg/logger"
	gachaMachine "github.com/Richard-inter/game/pkg/protocol/gachaMachine"
	"go.uber.org/zap"
)

type GachaMachineGRPCService struct {
	gachaMachine.UnimplementedGachaMachineServiceServer
	repo  repository.GachaMachineRepository
	redis *cache.RedisClient
	log   *zap.SugaredLogger
}

func NewGachaMachineGRPCService(repo repository.GachaMachineRepository, redis *cache.RedisClient) *GachaMachineGRPCService {
	return &GachaMachineGRPCService{
		repo:  repo,
		redis: redis,
		log:   logger.GetSugar(),
	}
}

func (s *GachaMachineGRPCService) CreateGachaPlayer(ctx context.Context, req *gachaMachine.CreateGachaPlayerReq) (*gachaMachine.CreateGachaPlayerResp, error) {
	// TODO: Implement gacha player creation
	s.log.Infow("CreateGachaPlayer called", "playerID", req.Player.BasePlayer.PlayerID)
	return &gachaMachine.CreateGachaPlayerResp{
		Player: req.Player,
	}, nil
}

func (s *GachaMachineGRPCService) GetGachaPlayerInfo(ctx context.Context, req *gachaMachine.GetGachaPlayerInfoReq) (*gachaMachine.GetGachaPlayerInfoResp, error) {
	// TODO: Implement get gacha player info
	s.log.Infow("GetGachaPlayerInfo called", "playerID", req.PlayerID)
	return nil, nil
}

func (s *GachaMachineGRPCService) AdjustPlayerGems(ctx context.Context, req *gachaMachine.AdjustPlayerGemsReq) (*gachaMachine.AdjustPlayerGemsResp, error) {
	// TODO: Implement adjust player gems
	s.log.Infow("AdjustPlayerGems called", "playerID", req.PlayerID, "amount", req.Amount, "type", req.Type)
	return &gachaMachine.AdjustPlayerGemsResp{
		PlayerID:   req.PlayerID,
		NewBalance: 0,
	}, nil
}

func (s *GachaMachineGRPCService) AdjustPlayerTickets(ctx context.Context, req *gachaMachine.AdjustPlayerTicketsReq) (*gachaMachine.AdjustPlayerTicketsResp, error) {
	// TODO: Implement adjust player tickets
	s.log.Infow("AdjustPlayerTickets called", "playerID", req.PlayerID, "amount", req.Amount, "type", req.Type)
	return &gachaMachine.AdjustPlayerTicketsResp{
		PlayerID:   req.PlayerID,
		NewBalance: 0,
	}, nil
}

func (s *GachaMachineGRPCService) AddItemToPlayer(ctx context.Context, req *gachaMachine.AddItemToPlayerReq) (*gachaMachine.AddItemToPlayerResp, error) {
	// TODO: Implement add item to player
	s.log.Infow("AddItemToPlayer called", "playerID", req.PlayerID, "itemID", req.ItemID)
	return &gachaMachine.AddItemToPlayerResp{
		PlayerID: req.PlayerID,
		ItemID:   req.ItemID,
		WasNew:   false,
	}, nil
}

func (s *GachaMachineGRPCService) CreateGachaPool(ctx context.Context, req *gachaMachine.CreateGachaPoolReq) (*gachaMachine.CreateGachaPoolResp, error) {
	// TODO: Implement create gacha pool
	s.log.Infow("CreateGachaPool called", "name", req.Name, "cost", req.Cost)
	return &gachaMachine.CreateGachaPoolResp{
		Pool: &gachaMachine.GachaPool{
			PoolID:   0,
			Name:     req.Name,
			Items:    req.Items,
			Cost:     req.Cost,
			Currency: req.Currency,
			MaxPulls: req.MaxPulls,
		},
	}, nil
}

func (s *GachaMachineGRPCService) GetGachaPoolInfo(ctx context.Context, req *gachaMachine.GetGachaPoolInfoReq) (*gachaMachine.GetGachaPoolInfoResp, error) {
	// TODO: Implement get gacha pool info
	s.log.Infow("GetGachaPoolInfo called", "poolID", req.PoolID)
	return nil, nil
}

func (s *GachaMachineGRPCService) GetAllGachaPools(ctx context.Context, req *gachaMachine.GetAllGachaPoolsReq) (*gachaMachine.GetAllGachaPoolsResp, error) {
	// TODO: Implement get all gacha pools
	s.log.Infow("GetAllGachaPools called")
	return &gachaMachine.GetAllGachaPoolsResp{
		Pools: []*gachaMachine.GachaPool{},
	}, nil
}

func (s *GachaMachineGRPCService) PullGacha(ctx context.Context, req *gachaMachine.PullGachaReq) (*gachaMachine.PullGachaResp, error) {
	// TODO: Implement pull gacha
	s.log.Infow("PullGacha called", "playerID", req.PlayerID, "poolID", req.PoolID, "pullCount", req.PullCount)
	return &gachaMachine.PullGachaResp{
		TotalPulls:        0,
		Results:           []*gachaMachine.GachaResult{},
		RemainingCurrency: 0,
	}, nil
}
