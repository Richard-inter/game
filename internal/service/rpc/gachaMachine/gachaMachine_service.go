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

func (s *GachaMachineGRPCService) CreateGachaItems(ctx context.Context, req *pb.CreateGachaItemsReq) (*pb.CreateGachaItemsResp, error) {
	items := make([]domain.GachaItem, 0, len(req.GachaItems))
	for _, item := range req.GachaItems {
		items = append(items, domain.GachaItem{
			Name:           item.Name,
			Rarity:         item.Rarity,
			PullPercentage: item.PullPercentage,
		})
	}

	resp, err := s.repo.CreateGachaItems(&items)
	if err != nil {
		return nil, err
	}
	createdItems := make([]*pb.Item, 0, len(*resp))
	for _, it := range *resp {
		createdItems = append(createdItems, &pb.Item{
			ItemID:         it.ID,
			Name:           it.Name,
			Rarity:         it.Rarity,
			PullPercentage: it.PullPercentage,
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

	created, err := s.repo.CreateGachaPlayer(gp)
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
	domainPlayer, err := s.repo.GetGachaPlayerInfo(req.PlayerID)
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
	updated, err := s.repo.AdjustPlayerCoin(req.PlayerID, req.Amount, req.Type)
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
	updated, err := s.repo.AdjustPlayerDiamond(req.PlayerID, req.Amount, req.Type)
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

	created, err := s.repo.CreateGachaMachine(g)
	if err != nil {
		return nil, err
	}

	fmt.Println(created.Items)
	items := make([]*pb.Item, 0, len(created.Items))
	for _, it := range created.Items {
		items = append(items, &pb.Item{
			ItemID:         it.Item.ID,
			Name:           it.Item.Name,
			Rarity:         it.Item.Rarity,
			PullPercentage: it.Item.PullPercentage,
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

	// get all machines
	if req.MachineID == 0 {
		machineDomainList, err := s.repo.GetAllGachaMachines()
		if err != nil {
			return nil, err
		}

		for _, resp := range machineDomainList {
			items := make([]*pb.Item, 0, len(resp.Items))
			for _, it := range resp.Items {
				items = append(items, &pb.Item{
					ItemID:         it.Item.ID,
					Name:           it.Item.Name,
					Rarity:         it.Item.Rarity,
					PullPercentage: it.Item.PullPercentage,
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
		resp, err := s.repo.GetGachaMachineInfo(req.MachineID)
		if err != nil {
			return nil, err
		}

		items := make([]*pb.Item, 0, len(resp.Items))
		for _, it := range resp.Items {
			items = append(items, &pb.Item{
				ItemID:         it.Item.ID,
				Name:           it.Item.Name,
				Rarity:         it.Item.Rarity,
				PullPercentage: it.Item.PullPercentage,
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
