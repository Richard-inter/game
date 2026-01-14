package clawmachine

import (
	"context"
	"fmt"

	"github.com/Richard-inter/game/internal/cache"
	"github.com/Richard-inter/game/internal/domain"
	"github.com/Richard-inter/game/internal/repository"
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

	// 1. Fetch domain model
	domainPlayer, err := s.repo.GetClawPlayerInfo(req.PlayerID)
	if err != nil {
		return nil, err
	}

	// 2. Map domain → protobuf
	clawPlayer := &pb.ClawPlayer{
		BasePlayer: &player.Player{
			PlayerID: domainPlayer.Player.ID,
			UserName: domainPlayer.Player.UserName,
		},
		Coin:    domainPlayer.Coin,
		Diamond: domainPlayer.Diamond,
	}

	// 3. Return response
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

	gameID, err := s.repo.AddGameHistory(req.PlayerID, &domain.ClawMachineGameRecord{
		PlayerID:      req.PlayerID,
		ClawMachineID: req.MachineID,
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
		fmt.Printf("Warning: failed to store game results in Redis: %v\n", err)
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
	resp, err := s.repo.GetClawMachineInfo(req.MachineID)
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

	return &pb.GetClawMachineInfoResp{
		Machine: &pb.ClawMachine{
			MachineID: resp.ID,
			Name:      resp.Name,
			Price:     resp.Price,
			MaxItem:   resp.MaxItem,
			Items:     items,
		},
	}, nil
}

func (s *ClawMachineGRPCServices) CreateClawMachine(ctx context.Context, req *pb.CreateClawMachineReq) (*pb.CreateClawMachineResp, error) {
	// Create domain claw machine with items
	c := &domain.ClawMachine{
		Name:    req.Name,
		Price:   req.Price,
		MaxItem: req.MaxItem,
		Items:   make([]domain.ClawMachineItem, 0, len(req.Items)),
	}

	// Convert protobuf items to domain items
	for _, item := range req.Items {
		c.Items = append(c.Items, domain.ClawMachineItem{
			ItemID: item.ItemID,
		})
	}

	// Create claw machine with items in a single transaction
	created, err := s.repo.CreateClawMachine(c)
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
	// Convert protobuf items to domain items
	items := make([]domain.Item, 0, len(req.ClawItems))
	for _, item := range req.ClawItems {
		items = append(items, domain.Item{
			Name:            item.Name,
			Rarity:          item.Rarity,
			SpawnPercentage: item.SpawnPercentage,
			CatchPercentage: item.CatchPercentage,
			MaxItemSpawned:  item.MaxItemSpawned,
		})
	}

	resp, err := s.repo.CreateClawItems(&items)
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
	// Create domain claw player
	cp := &domain.ClawPlayer{
		Player: domain.Player{
			ID:       req.Player.BasePlayer.PlayerID,
			UserName: req.Player.BasePlayer.UserName,
		},
		Coin:    req.Player.Coin,
		Diamond: req.Player.Diamond,
	}

	created, err := s.repo.CreateClawPlayer(cp)
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
	updated, err := s.repo.AdjustPlayerCoin(req.PlayerID, req.Amount, req.Type)
	if err != nil {
		return nil, err
	}

	return &pb.AdjustPlayerCoinResp{
		PlayerID:       updated.Player.ID,
		AdjustedAmount: updated.Coin,
	}, nil
}

func (s *ClawMachineGRPCServices) AdjustPlayerDiamond(ctx context.Context, req *pb.AdjustPlayerDiamondReq) (*pb.AdjustPlayerDiamondResp, error) {
	updated, err := s.repo.AdjustPlayerDiamond(req.PlayerID, req.Amount, req.Type)
	if err != nil {
		return nil, err
	}

	return &pb.AdjustPlayerDiamondResp{
		PlayerID:       updated.Player.ID,
		AdjustedAmount: updated.Diamond,
	}, nil
}

func (s *ClawMachineGRPCServices) AddTouchedItemRecord(ctx context.Context, req *pb.AddTouchedItemRecordReq) (*pb.AddTouchedItemRecordResp, error) {
	var storedResults []CatchResult
	err := s.redis.GetGameResults(ctx, req.GameID, &storedResults)
	if err != nil {
		return nil, fmt.Errorf("failed to load game results from Redis: %w", err)
	}

	var foundItem *CatchResult
	for _, result := range storedResults {
		if result.ItemID == req.ItemID {
			foundItem = &result
			break
		}
	}

	if foundItem.Success != *req.Catched {
		err := s.redis.DeleteGameResults(ctx, req.GameID)
		if err != nil {
			fmt.Printf("Warning: failed to delete game results from Redis: %v\n", err)
		}
		return nil, fmt.Errorf("catched value mismatch: expected %t, got %t", foundItem.Success, *req.Catched)
	}

	err = s.repo.AddTouchedItemRecord(req.GameID, req.ItemID, *req.Catched)
	if err != nil {
		return nil, fmt.Errorf("failed to update touched item record: %w", err)
	}

	err = s.redis.DeleteGameResults(ctx, req.GameID)
	if err != nil {
		// Log error but don't fail the request since validation passed
		fmt.Printf("Warning: failed to delete game results from Redis: %v\n", err)
	}

	return &pb.AddTouchedItemRecordResp{
		GameID:  req.GameID,
		ItemID:  req.ItemID,
		Catched: req.Catched,
	}, nil
}
