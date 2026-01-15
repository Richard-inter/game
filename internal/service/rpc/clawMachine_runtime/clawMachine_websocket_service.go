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

func (s *ClawMachineWebsocketService) GetPlayerInfoWs(
	ctx context.Context,
	req *pb.RuntimeRequest,
) (*pb.RuntimeResponse, error) {
	startReq := fbs.GetRootAsGetPlayerInfoWsReq(req.Payload, 0)
	playerID := startReq.PlayerId()

	// 1. Fetch domain model
	domainPlayer, err := s.repo.GetClawPlayerInfo(int64(playerID))
	if err != nil {
		return nil, err
	}

	builder := flatbuffers.NewBuilder(1024)
	usernameOffset := builder.CreateString(domainPlayer.Player.UserName)

	// Start table
	fbs.GetPlayerInfoWsRespStart(builder)

	// Add fields
	fbs.GetPlayerInfoWsRespAddPlayerId(builder, uint64(domainPlayer.Player.ID))
	fbs.GetPlayerInfoWsRespAddUsername(builder, usernameOffset)
	fbs.GetPlayerInfoWsRespAddCoin(builder, domainPlayer.Coin)
	fbs.GetPlayerInfoWsRespAddDiamond(builder, domainPlayer.Diamond)

	// End table
	resp := fbs.GetPlayerInfoWsRespEnd(builder)
	builder.Finish(resp)
	respBytes := builder.FinishedBytes()

	// Wrap response in Envelope
	envBuilder := flatbuffers.NewBuilder(1024)
	payloadOffset := envBuilder.CreateByteVector(respBytes)

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeGetPlayerInfoWsResp)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	envOffset := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(envOffset)

	// Return RuntimeResponse with envelope
	return &pb.RuntimeResponse{
		Payload: envBuilder.FinishedBytes(),
	}, nil
}

func (s *ClawMachineWebsocketService) AddTouchedItemRecordWs(
	ctx context.Context,
	req *pb.RuntimeRequest,
) (*pb.RuntimeResponse, error) {
	startReq := fbs.GetRootAsAddTouchedItemRecordReq(req.Payload, 0)
	gameID := startReq.GameId()
	itemID := startReq.ItemId()
	catched := startReq.Catched()

	// Load stored results from Redis
	var storedResults []CatchResult
	err := s.redis.GetGameResults(ctx, int64(gameID), &storedResults)
	if err != nil {
		return nil, fmt.Errorf("failed to load game results from Redis: %w", err)
	}

	var foundItem *CatchResult
	for _, result := range storedResults {
		if result.ItemID == int64(itemID) {
			foundItem = &result
			break
		}
	}

	if foundItem.Success != catched {
		err := s.redis.DeleteGameResults(ctx, int64(gameID))
		if err != nil {
			fmt.Printf("Warning: failed to delete game results from Redis: %v\n", err)
		}
		return nil, fmt.Errorf("catched value mismatch: expected %t, got %t", foundItem.Success, catched)
	}

	err = s.repo.AddTouchedItemRecord(int64(gameID), int64(itemID), catched)
	if err != nil {
		return nil, fmt.Errorf("failed to update touched item record: %w", err)
	}

	err = s.redis.DeleteGameResults(ctx, int64(gameID))
	if err != nil {
		// Log error but don't fail the request since validation passed
		fmt.Printf("Warning: failed to delete game results from Redis: %v\n", err)
	}

	// Build AddTouchedItemRecordResp
	builder := flatbuffers.NewBuilder(256)

	fbs.AddTouchedItemRecordRespStart(builder)
	fbs.AddTouchedItemRecordRespAddGameId(builder, uint64(gameID))
	fbs.AddTouchedItemRecordRespAddItemId(builder, uint64(itemID))
	fbs.AddTouchedItemRecordRespAddCatched(builder, catched)
	respOffset := fbs.AddTouchedItemRecordRespEnd(builder)
	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes()

	// Wrap response in Envelope
	envBuilder := flatbuffers.NewBuilder(256)
	payloadOffset := envBuilder.CreateByteVector(respBytes)

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeAddTouchedItemRecordResp)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	envOffset := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(envOffset)

	return &pb.RuntimeResponse{
		Payload: envBuilder.FinishedBytes(),
	}, nil
}
