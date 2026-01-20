package gachaMachine

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/cache"
	"github.com/Richard-inter/game/internal/domain"
	"github.com/Richard-inter/game/internal/repository"
	"github.com/Richard-inter/game/pkg/logger"
	pb "github.com/Richard-inter/game/pkg/protocol/gachaMachine"
	"github.com/Richard-inter/game/pkg/protocol/player"
)

type GachaMachineGRPCService struct {
	pb.UnimplementedGachaMachineServiceServer
	repo      repository.GachaMachineRepository
	redis     *cache.RedisClient
	streamKey string
	log       *zap.SugaredLogger
}

func NewGachaMachineGRPCService(repo repository.GachaMachineRepository, redis *cache.RedisClient, streamKey string) *GachaMachineGRPCService {
	return &GachaMachineGRPCService{
		repo:      repo,
		redis:     redis,
		streamKey: streamKey,
		log:       logger.GetSugar(),
	}
}

func (s *GachaMachineGRPCService) CreateGachaItems(ctx context.Context, req *pb.CreateGachaItemsReq) (*pb.CreateGachaItemsResp, error) {
	items := make([]domain.GachaItem, 0, len(req.GachaItems))
	for _, item := range req.GachaItems {
		items = append(items, domain.GachaItem{
			Name:       item.Name,
			Rarity:     item.Rarity,
			PullWeight: item.PullWeight,
		})
	}

	resp, err := s.repo.CreateGachaItems(ctx, &items)
	if err != nil {
		return nil, err
	}
	createdItems := make([]*pb.Item, 0, len(*resp))
	for _, it := range *resp {
		createdItems = append(createdItems, &pb.Item{
			ItemID:     it.ID,
			Name:       it.Name,
			Rarity:     it.Rarity,
			PullWeight: it.PullWeight,
		})
	}

	return &pb.CreateGachaItemsResp{
		GachaItems: createdItems,
	}, nil
}

func (s *GachaMachineGRPCService) CreateGachaPlayer(ctx context.Context, req *pb.CreateGachaPlayerReq) (*pb.CreateGachaPlayerResp, error) {
	gp := &domain.GachaPlayer{
		Player: domain.Player{
			ID:       req.Player.BasePlayer.PlayerID,
			UserName: req.Player.BasePlayer.UserName,
		},
		Coin:    req.Player.Coin,
		Diamond: req.Player.Diamond,
	}

	created, err := s.repo.CreateGachaPlayer(ctx, gp)
	if err != nil {
		s.log.Errorf("Failed to create gacha player: %v", err)
		return nil, fmt.Errorf("failed to create gacha player: %w", err)
	}

	gachaPlayer := &pb.GachaPlayer{
		BasePlayer: &player.Player{
			PlayerID: created.Player.ID,
			UserName: created.Player.UserName,
		},
		Coin:    created.Coin,
		Diamond: created.Diamond,
	}

	return &pb.CreateGachaPlayerResp{
		Player: gachaPlayer,
	}, nil
}

func (s *GachaMachineGRPCService) GetGachaPlayerInfo(ctx context.Context, req *pb.GetGachaPlayerInfoReq) (*pb.GetGachaPlayerInfoResp, error) {
	domainPlayer, err := s.repo.GetGachaPlayerInfo(ctx, req.PlayerID)
	if err != nil {
		s.log.Errorf("Failed to get gacha player info: %v", err)
		return nil, fmt.Errorf("failed to get gacha player info: %w", err)
	}

	gachaPlayer := &pb.GachaPlayer{
		BasePlayer: &player.Player{
			PlayerID: domainPlayer.Player.ID,
			UserName: domainPlayer.Player.UserName,
		},
		Coin:    domainPlayer.Coin,
		Diamond: domainPlayer.Diamond,
	}

	return &pb.GetGachaPlayerInfoResp{
		Player: gachaPlayer,
	}, nil
}

func (s *GachaMachineGRPCService) AdjustPlayerCoin(ctx context.Context, req *pb.AdjustPlayerCoinReq) (*pb.AdjustPlayerCoinResp, error) {
	updated, err := s.repo.AdjustPlayerCoin(ctx, req.PlayerID, req.Amount, req.Type)
	if err != nil {
		s.log.Errorf("Failed to adjust player coins: %v", err)
		return nil, fmt.Errorf("failed to adjust player coins: %w", err)
	}

	return &pb.AdjustPlayerCoinResp{
		PlayerID:       req.PlayerID,
		AdjustedAmount: updated.Coin,
	}, nil
}

func (s *GachaMachineGRPCService) AdjustPlayerDiamond(ctx context.Context, req *pb.AdjustPlayerDiamondReq) (*pb.AdjustPlayerDiamondResp, error) {
	updated, err := s.repo.AdjustPlayerDiamond(ctx, req.PlayerID, req.Amount, req.Type)
	if err != nil {
		s.log.Errorf("Failed to adjust player diamonds: %v", err)
		return nil, fmt.Errorf("failed to adjust player diamonds: %w", err)
	}

	return &pb.AdjustPlayerDiamondResp{
		PlayerID:       req.PlayerID,
		AdjustedAmount: updated.Diamond,
	}, nil
}

func (s *GachaMachineGRPCService) CreateGachaMachine(ctx context.Context, req *pb.CreateGachaMachineReq) (*pb.CreateGachaMachineResp, error) {
	g := &domain.GachaMachine{
		Name:          req.Name,
		Price:         req.Price,
		PriceTimesTen: req.PriceTimesTen,
		SuperRarePity: req.SuperRarePity,
		UltraRarePity: req.UltraRarePity,
		Items:         make([]domain.GachaMachineItem, 0, len(req.Items)),
	}

	for _, item := range req.Items {
		g.Items = append(g.Items, domain.GachaMachineItem{
			ItemID: item.ItemID,
		})
	}

	created, err := s.repo.CreateGachaMachine(ctx, g)
	if err != nil {
		return nil, err
	}

	fmt.Println(created.Items)
	items := make([]*pb.Item, 0, len(created.Items))
	for _, it := range created.Items {
		items = append(items, &pb.Item{
			ItemID:     it.Item.ID,
			Name:       it.Item.Name,
			Rarity:     it.Item.Rarity,
			PullWeight: it.Item.PullWeight,
		})
	}

	return &pb.CreateGachaMachineResp{
		Machine: &pb.GachaMachine{
			MachineID:     created.ID,
			Name:          created.Name,
			Price:         created.Price,
			PriceTimesTen: created.PriceTimesTen,
			SuperRarePity: created.SuperRarePity,
			UltraRarePity: created.UltraRarePity,
			Items:         items,
		},
	}, nil
}

func (s *GachaMachineGRPCService) GetGachaMachineInfo(ctx context.Context, req *pb.GetGachaMachineInfoReq) (*pb.GetGachaMachineInfoResp, error) {
	var machines []*pb.GachaMachine
	if req.MachineID == 0 {
		machineDomainList, err := s.repo.GetAllGachaMachines(ctx)
		if err != nil {
			return nil, err
		}

		for _, resp := range machineDomainList {
			items := make([]*pb.Item, 0, len(resp.Items))
			for _, it := range resp.Items {
				items = append(items, &pb.Item{
					ItemID:     it.Item.ID,
					Name:       it.Item.Name,
					Rarity:     it.Item.Rarity,
					PullWeight: it.Item.PullWeight,
				})
			}

			machines = append(machines, &pb.GachaMachine{
				MachineID:     resp.ID,
				Name:          resp.Name,
				Price:         resp.Price,
				PriceTimesTen: resp.PriceTimesTen,
				SuperRarePity: resp.SuperRarePity,
				UltraRarePity: resp.UltraRarePity,
				Items:         items,
			})
		}
	} else {
		resp, err := s.repo.GetGachaMachineInfo(ctx, req.MachineID)
		if err != nil {
			return nil, err
		}

		items := make([]*pb.Item, 0, len(resp.Items))
		for _, it := range resp.Items {
			items = append(items, &pb.Item{
				ItemID:     it.Item.ID,
				Name:       it.Item.Name,
				Rarity:     it.Item.Rarity,
				PullWeight: it.Item.PullWeight,
			})
		}

		machines = append(machines, &pb.GachaMachine{
			MachineID:     resp.ID,
			Name:          resp.Name,
			Price:         resp.Price,
			PriceTimesTen: resp.PriceTimesTen,
			SuperRarePity: resp.SuperRarePity,
			UltraRarePity: resp.UltraRarePity,
			Items:         items,
		})
	}

	return &pb.GetGachaMachineInfoResp{
		Machine: machines,
	}, nil
}

func (s *GachaMachineGRPCService) GetPullResult(ctx context.Context, req *pb.GetPullResultReq) (*pb.GetPullResultResp, error) {
	session := &domain.GachaPullSession{
		GachaMachineID: req.MachineID,
		PlayerID:       req.PlayerID,
		PullCount:      req.PullCount,
	}

	if req.PullCount == 1 {
		itemID, err := s.PullGachaSingle(ctx, req.MachineID, req.PlayerID)
		if err != nil {
			return nil, err
		}

		if err := s.redis.PublishGachaEvent(ctx, s.streamKey, session, itemID); err != nil {
			s.log.Errorf("Failed to publish gacha pull history to stream: %v", err)
			// Continue even if stream publish fails
		}

		return &pb.GetPullResultResp{
			ItemIDs: []int64{itemID},
		}, nil
	}

	itemIDs, err := s.PullGachaByMachineIDMulti(ctx, req.MachineID, req.PlayerID, int(req.PullCount))
	if err != nil {
		return nil, err
	}

	// Publish all pull histories to stream with session data in a single message
	if err := s.redis.PublishGachaEvent(ctx, s.streamKey, session, itemIDs); err != nil {
		s.log.Errorf("Failed to publish gacha pull history to stream: %v", err)
		// Continue even if stream publish fails
	}

	return &pb.GetPullResultResp{
		ItemIDs: itemIDs,
	}, nil
}

func (s *GachaMachineGRPCService) AddGameToHistory(ctx context.Context, session domain.GachaPullSession, itemIDs []int64) error {
	createdSession, err := s.repo.CreateGachaPullSession(ctx, &session)
	if err != nil {
		s.log.Errorf("Failed to create gacha pull session: %v", err)
		return err
	}

	for _, itemID := range itemIDs {
		history := &domain.GachaPullHistory{
			GachaPullSessionID: createdSession.ID,
			ItemID:             itemID,
		}

		_, err = s.repo.CreateGachaPullHistories(ctx, &[]domain.GachaPullHistory{*history})
		if err != nil {
			s.log.Errorf("Failed to create gacha pull history: %v", err)
			return err
		}
	}
	return nil
}
