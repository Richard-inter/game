package clawmachine

import (
	"context"
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/Richard-inter/game/internal/cache"
	"github.com/Richard-inter/game/internal/domain"
	"github.com/Richard-inter/game/internal/repository"
	pb "github.com/Richard-inter/game/pkg/protocol/clawMachine_Websocket"
	fbs "github.com/Richard-inter/game/pkg/protocol/clawMachine_Websocket/clawMachine"
)

type ClawMachineWebsocketService struct {
	pb.UnimplementedClawMachineRuntimeServiceServer
	repo  repository.ClawMachineRepository
	redis *cache.RedisClient
}

func NewClawMachineWebsocketService(repo repository.ClawMachineRepository, redis *cache.RedisClient) *ClawMachineWebsocketService {
	return &ClawMachineWebsocketService{
		repo:  repo,
		redis: redis,
	}
}

func (s *ClawMachineWebsocketService) StartClawGameWs(
	ctx context.Context,
	req *pb.RuntimeRequest,
) (*pb.RuntimeResponse, error) {

	// ---- 1. Decode StartClawGameReq (payload only, NOT Envelope) ----
	startReq := fbs.GetRootAsStartClawGameReq(req.Payload, 0)
	playerID := startReq.PlayerId()
	machineID := startReq.MachineId()

	fmt.Println("StartClawGameReq received")
	fmt.Println("  PlayerID :", playerID)
	fmt.Println("  MachineID:", machineID)

	// Validate input
	if playerID <= 0 || machineID <= 0 {
		return nil, fmt.Errorf("invalid player ID or machine ID")
	}

	// Pre-determine catch results
	results, err := s.PreDetermineCatchResults(ctx, int64(machineID))
	if err != nil {
		return nil, fmt.Errorf("failed to pre-determine catch results: %w", err)
	}

	// Charge player for playing
	err = s.PlayMachine(ctx, int64(playerID), int64(machineID))
	if err != nil {
		return nil, fmt.Errorf("failed to charge player: %w", err)
	}

	// Create game history record
	gameID, err := s.repo.AddGameHistory(int64(playerID), &domain.ClawMachineGameRecord{
		PlayerID:      int64(playerID),
		ClawMachineID: int64(machineID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create game history: %w", err)
	}

	// Store results in Redis for later validation
	err = s.redis.StoreGameResults(ctx, gameID, results)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to store game results in Redis: %v\n", err)
	}

	// ---- 3. Build StartClawGameResp ----
	builder := flatbuffers.NewBuilder(1024)

	// Build ClawResult objects from real results
	resultOffsets := make([]flatbuffers.UOffsetT, len(results))
	for i := len(results) - 1; i >= 0; i-- {
		fbs.ClawResultStart(builder)
		fbs.ClawResultAddItemId(builder, uint64(results[i].ItemID))
		fbs.ClawResultAddCatched(builder, results[i].Success)
		resultOffsets[i] = fbs.ClawResultEnd(builder)
	}

	// Build results vector
	fbs.StartClawGameRespStartResultsVector(builder, len(resultOffsets))
	for i := len(resultOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(resultOffsets[i])
	}
	resultsVector := builder.EndVector(len(resultOffsets))

	// Build StartClawGameResp table
	fbs.StartClawGameRespStart(builder)
	fbs.StartClawGameRespAddGameId(builder, uint64(gameID))
	fbs.StartClawGameRespAddResults(builder, resultsVector)
	respOffset := fbs.StartClawGameRespEnd(builder) // ✅ Must finish table

	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes() // safe, complete table

	// ---- 4. Wrap response in Envelope ----
	envBuilder := flatbuffers.NewBuilder(1024)
	payloadOffset := envBuilder.CreateByteVector(respBytes) // Wrap table as byte vector

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeStartClawGameResp)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	envOffset := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(envOffset)

	// ---- 5. Return Envelope bytes ----
	return &pb.RuntimeResponse{
		Payload: envBuilder.FinishedBytes(),
	}, nil
}

// PreDetermineCatchResults generates a list of pre-determined catch results for all items
func (s *ClawMachineWebsocketService) PreDetermineCatchResults(
	ctx context.Context,
	machineID int64,
) ([]*CatchResult, error) {
	// Get machine info to access items and their catch percentages
	clawMachine, err := s.repo.GetClawMachineInfo(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine info: %w", err)
	}

	if len(clawMachine.Items) == 0 {
		return nil, fmt.Errorf("no items in machine to catch from")
	}

	results := make([]*CatchResult, 0, len(clawMachine.Items))

	// Generate pre-determined result for each item in the machine
	for _, item := range clawMachine.Items {
		catchWeight := item.Item.CatchPercentage
		if catchWeight == 0 {
			return nil, fmt.Errorf("database error: item %s (ID: %d) has zero catch percentage", item.Item.Name, item.Item.ID)
		}

		// Determine if catch is successful based on the item's catch percentage
		catchSuccess := Roll(int(catchWeight))

		results = append(results, &CatchResult{
			ItemID:  item.ID,
			Name:    item.Item.Name,
			Success: catchSuccess,
		})
	}

	return results, nil
}

func (s *ClawMachineWebsocketService) PlayMachine(
	ctx context.Context,
	playerID int64,
	machineID int64,
) error {
	clawMachine, err := s.repo.GetClawMachineInfo(machineID)
	if err != nil {
		return fmt.Errorf("failed to get machine info: %w", err)
	}

	resp, err := s.repo.AdjustPlayerCoin(playerID, int64(clawMachine.Price), "minus")
	if err != nil {
		return fmt.Errorf("failed to adjust player coin: %w", err)
	}

	if resp.Coin < 0 {
		// Revert adjustment
		_, _ = s.repo.AdjustPlayerCoin(playerID, int64(clawMachine.Price), "plus")
		return fmt.Errorf("insufficient coins to play the claw machine")
	}

	return nil
}
