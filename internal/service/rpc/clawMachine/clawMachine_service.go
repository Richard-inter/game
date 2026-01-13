package clawmachine

import (
	"context"
	"fmt"
	"time"

	"github.com/Richard-inter/game/internal/domain"
	"github.com/Richard-inter/game/internal/repository"
	pb "github.com/Richard-inter/game/pkg/protocol/clawMachine"
	"github.com/Richard-inter/game/pkg/protocol/player"
)

// ClawMachineGRPCService implements the ClawMachineService gRPC service
type ClawMachineGRPCServices struct {
	pb.UnimplementedClawMachineServiceServer
	repo repository.ClawMachineRepository
}

// NewClawMachineGRPCService creates a new ClawMachineGRPCService
func NewClawMachineGRPCService(repo repository.ClawMachineRepository) *ClawMachineGRPCServices {
	return &ClawMachineGRPCServices{
		repo: repo,
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
	// 1. Validate request
	if req.PlayerID <= 0 || req.MachineID <= 0 {
		return nil, fmt.Errorf("invalid player ID or machine ID")
	}

	// 2. Get pre-determined catch results
	results, err := s.PreDetermineCatchResults(ctx, req.MachineID)
	fmt.Println(results)
	if err != nil {
		return nil, fmt.Errorf("failed to pre-determine catch results: %w", err)
	}

	// 3. Charge player for playing the machine
	err = s.PlayMachine(ctx, req.PlayerID, req.MachineID)
	if err != nil {
		return nil, fmt.Errorf("failed to charge player: %w", err)
	}

	// 4. Convert domain results to protobuf
	protoResults := make([]*pb.ClawResult, 0, len(results))
	for _, result := range results {
		clawResult := &pb.ClawResult{
			ItemID:  result.ItemID,
			Catched: &result.Success,
		}
		protoResults = append(protoResults, clawResult)
	}

	// 5. Generate unique game ID
	gameID := time.Now().UnixNano()

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
