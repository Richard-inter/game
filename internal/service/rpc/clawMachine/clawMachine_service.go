package clawmachine

import (
	"context"
	"fmt"

	"github.com/Richard-inter/game/internal/cache"
	"github.com/Richard-inter/game/internal/domain"
	"github.com/Richard-inter/game/internal/repository"
	"github.com/Richard-inter/game/pkg/logger"
	pb "github.com/Richard-inter/game/pkg/protocol/clawMachine"
	"github.com/Richard-inter/game/pkg/protocol/player"
)

// ClawMachineGRPCService implements the ClawMachineService gRPC service
type ClawMachineGRPCServices struct {
	pb.UnimplementedClawMachineServiceServer
	repo  repository.ClawMachineRepository
	redis *cache.RedisClient
}

// NewClawMachineGRPCService creates a new ClawMachineGRPCService
func NewClawMachineGRPCService(repo repository.ClawMachineRepository, redis *cache.RedisClient) *ClawMachineGRPCServices {
	return &ClawMachineGRPCServices{
		repo:  repo,
		redis: redis,
	}
}

func (s *ClawMachineGRPCServices) GetClawPlayerInfo(
	ctx context.Context,
	req *pb.GetClawPlayerInfoReq,
) (*pb.GetClawPlayerInfoResp, error) {
	domainPlayer, err := s.repo.GetClawPlayerInfo(ctx, req.PlayerID)
	if err != nil {
		return nil, err
	}

	clawPlayer := &pb.ClawPlayer{
		BasePlayer: &player.Player{
			PlayerID: domainPlayer.Player.ID,
			UserName: domainPlayer.Player.UserName,
		},
		Coin:    domainPlayer.Coin,
		Diamond: domainPlayer.Diamond,
	}

	return &pb.GetClawPlayerInfoResp{
		Player: clawPlayer,
	}, nil
}

func (s *ClawMachineGRPCServices) StartClawGame(ctx context.Context, req *pb.StartClawGameReq) (*pb.StartClawGameResp, error) {
	if req.PlayerID <= 0 || req.MachineID <= 0 {
		return nil, fmt.Errorf("invalid player ID or machine ID")
	}

	results, err := s.PreDetermineCatchResults(ctx, req.MachineID)
	if err != nil {
		return nil, fmt.Errorf("failed to pre-determine catch results: %w", err)
	}

	err = s.PlayMachine(ctx, req.PlayerID, req.MachineID)
	if err != nil {
		return nil, fmt.Errorf("failed to charge player: %w", err)
	}

	gameID, err := s.repo.AddGameHistory(ctx, req.PlayerID, &domain.ClawMachineGameRecord{
		PlayerID:      req.PlayerID,
		ClawMachineID: req.MachineID,
		CreatedBy:     fmt.Sprintf("%d", req.PlayerID),
		UpdatedBy:     fmt.Sprintf("%d", req.PlayerID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create game history: %w", err)
	}

	protoResults := make([]*pb.ClawResult, 0, len(results))
	for _, result := range results {
		clawResult := &pb.ClawResult{
			ItemID:  result.ItemID,
			Catched: &result.Success,
		}
		protoResults = append(protoResults, clawResult)
	}

	err = s.redis.StoreGameResults(ctx, gameID, results)
	if err != nil {
		// Log error but don't fail the request
		logger.GetSugar().Warnf("failed to store game results in Redis: %v", err)
	}

	return &pb.StartClawGameResp{
		GameID:  gameID,
		Results: protoResults,
	}, nil
}

func (s *ClawMachineGRPCServices) GetClawMachineInfo(
	ctx context.Context,
	req *pb.GetClawMachineInfoReq,
) (*pb.GetClawMachineInfoResp, error) {
	var machines []*pb.ClawMachine

	// get all machines
	if req.MachineID == 0 {
		machineDomainList, err := s.repo.GetAllClawMachines(ctx)
		if err != nil {
			return nil, err
		}

		for _, resp := range machineDomainList {
			items := make([]*pb.Item, 0, len(resp.Items))
			for _, it := range resp.Items {
				items = append(items, &pb.Item{
					ItemID:          it.Item.ID,
					Name:            it.Item.Name,
					Rarity:          it.Item.Rarity,
					SpawnPercentage: it.Item.SpawnPercentage,
					CatchPercentage: it.Item.CatchPercentage,
					MaxItemSpawned:  it.Item.MaxItemSpawned,
				})
			}

			machines = append(machines, &pb.ClawMachine{
				MachineID: resp.ID,
				Name:      resp.Name,
				Price:     resp.Price,
				MaxItem:   resp.MaxItem,
				Items:     items,
			})
		}
	} else {
		resp, err := s.repo.GetClawMachineInfo(ctx, req.MachineID)
		if err != nil {
			return nil, err
		}

		items := make([]*pb.Item, 0, len(resp.Items))
		for _, it := range resp.Items {
			items = append(items, &pb.Item{
				ItemID:          it.Item.ID,
				Name:            it.Item.Name,
				Rarity:          it.Item.Rarity,
				SpawnPercentage: it.Item.SpawnPercentage,
				CatchPercentage: it.Item.CatchPercentage,
				MaxItemSpawned:  it.Item.MaxItemSpawned,
			})
		}

		machines = append(machines, &pb.ClawMachine{
			MachineID: resp.ID,
			Name:      resp.Name,
			Price:     resp.Price,
			MaxItem:   resp.MaxItem,
			Items:     items,
		})
	}

	return &pb.GetClawMachineInfoResp{
		Machine: machines,
	}, nil
}

func (s *ClawMachineGRPCServices) CreateClawMachine(ctx context.Context, req *pb.CreateClawMachineReq) (*pb.CreateClawMachineResp, error) {
	c := &domain.ClawMachine{
		Name:    req.Name,
		Price:   req.Price,
		MaxItem: req.MaxItem,
		Items:   make([]domain.ClawMachineItem, 0, len(req.Items)),
	}

	for _, item := range req.Items {
		c.Items = append(c.Items, domain.ClawMachineItem{
			ItemID: item.ItemID,
		})
	}

	created, err := s.repo.CreateClawMachine(ctx, c)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.Item, 0, len(created.Items))
	for _, it := range created.Items {
		items = append(items, &pb.Item{
			ItemID:          it.Item.ID,
			Name:            it.Item.Name,
			Rarity:          it.Item.Rarity,
			SpawnPercentage: it.Item.SpawnPercentage,
			CatchPercentage: it.Item.CatchPercentage,
			MaxItemSpawned:  it.Item.MaxItemSpawned,
		})
	}

	return &pb.CreateClawMachineResp{
		Machine: &pb.ClawMachine{
			MachineID: created.ID,
			Name:      created.Name,
			Price:     created.Price,
			MaxItem:   created.MaxItem,
			Items:     items,
		},
	}, nil
}

func (s *ClawMachineGRPCServices) CreateClawItems(ctx context.Context, req *pb.CreateClawItemsReq) (*pb.CreateClawItemsResp, error) {
	items := make([]domain.ClawItem, 0, len(req.ClawItems))
	for _, item := range req.ClawItems {
		items = append(items, domain.ClawItem{
			Name:            item.Name,
			Rarity:          item.Rarity,
			SpawnPercentage: item.SpawnPercentage,
			CatchPercentage: item.CatchPercentage,
			MaxItemSpawned:  item.MaxItemSpawned,
		})
	}

	resp, err := s.repo.CreateClawItems(ctx, &items)
	if err != nil {
		return nil, err
	}
	createdItems := make([]*pb.Item, 0, len(*resp))
	for _, it := range *resp {
		createdItems = append(createdItems, &pb.Item{
			ItemID:          it.ID,
			Name:            it.Name,
			Rarity:          it.Rarity,
			SpawnPercentage: it.SpawnPercentage,
			CatchPercentage: it.CatchPercentage,
			MaxItemSpawned:  it.MaxItemSpawned,
		})
	}

	return &pb.CreateClawItemsResp{
		ClawItems: createdItems,
	}, nil
}

func (s *ClawMachineGRPCServices) CreateClawPlayer(ctx context.Context, req *pb.CreateClawPlayerReq) (*pb.CreateClawPlayerResp, error) {
	cp := &domain.ClawPlayer{
		Player: domain.Player{
			ID:       req.Player.BasePlayer.PlayerID,
			UserName: req.Player.BasePlayer.UserName,
		},
		Coin:    req.Player.Coin,
		Diamond: req.Player.Diamond,
	}

	created, err := s.repo.CreateClawPlayer(ctx, cp)
	if err != nil {
		return nil, err
	}

	clawPlayer := &pb.ClawPlayer{
		BasePlayer: &player.Player{
			PlayerID: created.Player.ID,
			UserName: created.Player.UserName,
		},
		Coin:    created.Coin,
		Diamond: created.Diamond,
	}

	return &pb.CreateClawPlayerResp{
		Player: clawPlayer,
	}, nil
}

func (s *ClawMachineGRPCServices) AdjustPlayerCoin(ctx context.Context, req *pb.AdjustPlayerCoinReq) (*pb.AdjustPlayerCoinResp, error) {
	updated, err := s.repo.AdjustPlayerCoin(ctx, req.PlayerID, req.Amount, req.Type)
	if err != nil {
		return nil, err
	}

	return &pb.AdjustPlayerCoinResp{
		PlayerID:       updated.Player.ID,
		AdjustedAmount: updated.Coin,
	}, nil
}

func (s *ClawMachineGRPCServices) AdjustPlayerDiamond(ctx context.Context, req *pb.AdjustPlayerDiamondReq) (*pb.AdjustPlayerDiamondResp, error) {
	updated, err := s.repo.AdjustPlayerDiamond(ctx, req.PlayerID, req.Amount, req.Type)
	if err != nil {
		return nil, err
	}

	return &pb.AdjustPlayerDiamondResp{
		PlayerID:       updated.Player.ID,
		AdjustedAmount: updated.Diamond,
	}, nil
}

func (s *ClawMachineGRPCServices) AddTouchedItemRecord(
	ctx context.Context,
	req *pb.AddTouchedItemRecordReq,
) (*pb.AddTouchedItemRecordResp, error) {
	var storedResults []CatchResult
	if err := s.redis.GetGameResults(ctx, req.GameID, &storedResults); err != nil {
		return nil, fmt.Errorf("failed to load game results: %w", err)
	}

	var catched bool
	if req.ItemID != 0 {
		var serverResult *CatchResult
		for i := range storedResults {
			if storedResults[i].ItemID == req.ItemID {
				serverResult = &storedResults[i]
				break
			}
		}

		if serverResult == nil {
			return nil, fmt.Errorf("item %d not found in game %d", req.ItemID, req.GameID)
		}

		catched = serverResult.Success

		if req.Catched != nil && *req.Catched != catched {
			logger.GetSugar().Warnf(
				"client desync - game=%d item=%d client=%t server=%t",
				req.GameID, req.ItemID, *req.Catched, catched,
			)
		}
	}
	var itemID *int64
	itemID = &req.ItemID
	if req.ItemID == 0 {
		itemID = nil
		catched = false
	}

	game, err := s.repo.AddTouchedItemRecord(ctx, req.GameID, itemID, catched)
	if err != nil {
		return nil, fmt.Errorf("failed to persist record: %w", err)
	}

	payout := int64(0)
	if catched {
		payout = GetRarityValue(game.TouchedItem.Rarity, game.Machine.Price)
	}

	// Update RTP state for payout only (revenue already tracked in PlayMachine)
	err = s.repo.UpdateClawMachineRTP(ctx, game.Machine.ID, 0, payout)
	if err != nil {
		return nil, fmt.Errorf("failed to update claw machine RTP: %w", err)
	}

	_ = s.redis.DeleteGameResults(ctx, req.GameID)

	return &pb.AddTouchedItemRecordResp{
		GameID:  req.GameID,
		ItemID:  req.ItemID,
		Catched: &catched,
	}, nil
}

func (s *ClawMachineGRPCServices) DeleteClawPlayer(ctx context.Context, req *pb.DeleteClawPlayerReq) (*pb.DeleteClawPlayerResp, error) {
	err := s.repo.DeleteClawPlayer(ctx, req.PlayerID)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteClawPlayerResp{
		PlayerID: req.PlayerID,
		Success:  true,
	}, nil
}

func (s *ClawMachineGRPCServices) GetGameHistory(ctx context.Context, req *pb.GetGameHistoryReq) (*pb.GetGameHistoryResp, error) {
	gameRecords, err := s.repo.GetGameHistory(ctx, req.PlayerID)
	if err != nil {
		return nil, err
	}

	protoRecords := make([]*pb.ClawMachineGameRecord, 0, len(gameRecords))
	for _, record := range gameRecords {
		protoRecords = append(protoRecords, &pb.ClawMachineGameRecord{
			GameID:        record.ID,
			ClawMachineID: record.ClawMachineID,
			PlayerID:      record.PlayerID,
			TouchedItemID: *record.TouchedItemID,
			Catched:       *record.Catched,
			CreatedAt:     record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.GetGameHistoryResp{
		GameRecords: protoRecords,
	}, nil
}

func (s *ClawMachineGRPCServices) UpdateClawMachineItems(ctx context.Context, req *pb.UpdateClawMachineItemsReq) (*pb.UpdateClawMachineItemsResp, error) {
	items := make([]domain.ClawMachineItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.ClawMachineItem{
			ItemID: item.ItemID,
		})
	}

	updatedMachine, err := s.repo.UpdateClawMachineItems(ctx, req.MachineID, items)
	if err != nil {
		return nil, err
	}

	// Validate that the machine was updated successfully
	if updatedMachine == nil {
		return nil, fmt.Errorf("failed to update machine items: machine not found")
	}

	return &pb.UpdateClawMachineItemsResp{
		MachineID: req.MachineID,
		Success:   true,
	}, nil
}

func (s *ClawMachineGRPCServices) DeleteClawMachine(ctx context.Context, req *pb.DeleteClawMachineReq) (*pb.DeleteClawMachineResp, error) {
	err := s.repo.DeleteClawMachine(ctx, req.MachineID)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteClawMachineResp{
		MachineID: req.MachineID,
		Success:   true,
	}, nil
}

func (s *ClawMachineGRPCServices) DeleteClawItems(ctx context.Context, req *pb.DeleteClawItemsReq) (*pb.DeleteClawItemsResp, error) {
	err := s.repo.DeleteClawItems(ctx, req.ItemIDs)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteClawItemsResp{
		ItemIDs: req.ItemIDs,
		Success: true,
	}, nil
}

func (s *ClawMachineGRPCServices) UpdateClawMachineTargetRTP(ctx context.Context, req *pb.UpdateClawMachineTargetRTPReq) (*pb.UpdateClawMachineTargetRTPResp, error) {
	err := s.repo.UpdateClawMachineTargetRTP(ctx, req.MachineID, req.TargetRTP)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateClawMachineTargetRTPResp{
		MachineID: req.MachineID,
		Success:   true,
	}, nil
}
