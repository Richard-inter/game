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

func (_ *ClawMachineWebsocketService) buildEnvelopeResponse(messageType fbs.MessageType, payloadBytes []byte) *pb.RuntimeResponse {
	builder := flatbuffers.NewBuilder(len(payloadBytes) + 256)
	payloadOffset := builder.CreateByteVector(payloadBytes)

	fbs.EnvelopeStart(builder)
	fbs.EnvelopeAddType(builder, messageType)
	fbs.EnvelopeAddPayload(builder, payloadOffset)
	envOffset := fbs.EnvelopeEnd(builder)
	builder.Finish(envOffset)

	return &pb.RuntimeResponse{
		Payload: builder.FinishedBytes(),
	}
}

func (s *ClawMachineWebsocketService) StartClawGameWs(
	ctx context.Context,
	req *pb.RuntimeRequest,
) (*pb.RuntimeResponse, error) {
	startReq := fbs.GetRootAsStartClawGameReq(req.Payload, 0)
	playerID := startReq.PlayerId()
	machineID := startReq.MachineId()

	fmt.Println("StartClawGameReq received")
	fmt.Println("  PlayerID :", playerID)
	fmt.Println("  MachineID:", machineID)

	if playerID <= 0 || machineID <= 0 {
		return nil, fmt.Errorf("invalid player ID or machine ID")
	}

	results, err := s.PreDetermineCatchResults(ctx, int64(machineID))
	if err != nil {
		return nil, fmt.Errorf("failed to pre-determine catch results: %w", err)
	}

	err = s.PlayMachine(ctx, int64(playerID), int64(machineID))
	if err != nil {
		return nil, fmt.Errorf("failed to charge player: %w", err)
	}

	gameID, err := s.repo.AddGameHistory(ctx, int64(playerID), &domain.ClawMachineGameRecord{
		PlayerID:      int64(playerID),
		ClawMachineID: int64(machineID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create game history: %w", err)
	}

	err = s.redis.StoreGameResults(ctx, gameID, results)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to store game results in Redis: %v\n", err)
	}

	builder := flatbuffers.NewBuilder(1024)
	resultOffsets := make([]flatbuffers.UOffsetT, len(results))
	for i := len(results) - 1; i >= 0; i-- {
		fbs.ClawResultStart(builder)
		fbs.ClawResultAddItemId(builder, uint64(results[i].ItemID))
		fbs.ClawResultAddCatched(builder, results[i].Success)
		resultOffsets[i] = fbs.ClawResultEnd(builder)
	}

	fbs.StartClawGameRespStartResultsVector(builder, len(resultOffsets))
	for i := len(resultOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(resultOffsets[i])
	}
	resultsVector := builder.EndVector(len(resultOffsets))

	fbs.StartClawGameRespStart(builder)
	fbs.StartClawGameRespAddGameId(builder, uint64(gameID))
	fbs.StartClawGameRespAddResults(builder, resultsVector)
	respOffset := fbs.StartClawGameRespEnd(builder)

	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(fbs.MessageTypeStartClawGameResp, respBytes), nil
}

func (s *ClawMachineWebsocketService) GetPlayerInfoWs(
	ctx context.Context,
	req *pb.RuntimeRequest,
) (*pb.RuntimeResponse, error) {
	startReq := fbs.GetRootAsGetPlayerInfoWsReq(req.Payload, 0)
	playerID := startReq.PlayerId()

	domainPlayer, err := s.repo.GetClawPlayerInfo(ctx, int64(playerID))
	if err != nil {
		return nil, err
	}

	builder := flatbuffers.NewBuilder(1024)
	usernameOffset := builder.CreateString(domainPlayer.Player.UserName)

	fbs.GetPlayerInfoWsRespStart(builder)

	fbs.GetPlayerInfoWsRespAddPlayerId(builder, uint64(domainPlayer.Player.ID))
	fbs.GetPlayerInfoWsRespAddUsername(builder, usernameOffset)
	fbs.GetPlayerInfoWsRespAddCoin(builder, domainPlayer.Coin)
	fbs.GetPlayerInfoWsRespAddDiamond(builder, domainPlayer.Diamond)

	resp := fbs.GetPlayerInfoWsRespEnd(builder)
	builder.Finish(resp)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(fbs.MessageTypeGetPlayerInfoWsResp, respBytes), nil
}

func (s *ClawMachineWebsocketService) AddTouchedItemRecordWs(
	ctx context.Context,
	req *pb.RuntimeRequest,
) (*pb.RuntimeResponse, error) {
	startReq := fbs.GetRootAsAddTouchedItemRecordReq(req.Payload, 0)
	gameID := startReq.GameId()
	itemID := startReq.ItemId()
	catched := startReq.Catched()

	var storedResults []CatchResult
	err := s.redis.GetGameResults(ctx, int64(gameID), &storedResults)
	if err != nil {
		return nil, fmt.Errorf("failed to load game results from Redis: %w", err)
	}

	var foundItem *CatchResult
	for i := range storedResults {
		if storedResults[i].ItemID == int64(itemID) {
			foundItem = &storedResults[i]
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

	err = s.repo.AddTouchedItemRecord(ctx, int64(gameID), int64(itemID), catched)
	if err != nil {
		return nil, fmt.Errorf("failed to update touched item record: %w", err)
	}

	err = s.redis.DeleteGameResults(ctx, int64(gameID))
	if err != nil {
		// Log error but don't fail the request since validation passed
		fmt.Printf("Warning: failed to delete game results from Redis: %v\n", err)
	}

	builder := flatbuffers.NewBuilder(256)

	fbs.AddTouchedItemRecordRespStart(builder)
	fbs.AddTouchedItemRecordRespAddGameId(builder, gameID)
	fbs.AddTouchedItemRecordRespAddItemId(builder, itemID)
	fbs.AddTouchedItemRecordRespAddCatched(builder, catched)
	respOffset := fbs.AddTouchedItemRecordRespEnd(builder)
	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(fbs.MessageTypeAddTouchedItemRecordResp, respBytes), nil
}

func (s *ClawMachineWebsocketService) SpawnItemWs(
	ctx context.Context,
	req *pb.RuntimeRequest,
) (*pb.RuntimeResponse, error) {
	startReq := fbs.GetRootAsSpawnItemReq(req.Payload, 0)
	machineID := startReq.MachineId()

	result, err := s.SpawnMachineItems(ctx, int64(machineID))
	if err != nil {
		return nil, fmt.Errorf("failed to spawn item: %w", err)
	}

	items := make([]uint64, len(result))
	for i, v := range result {
		items[i] = uint64(v)
	}

	builder := flatbuffers.NewBuilder(256)

	fbs.SpawnItemRespStartItemsVector(builder, len(result))
	for i := len(result) - 1; i >= 0; i-- {
		builder.PrependUint64(uint64(result[i]))
	}
	itemsVector := builder.EndVector(len(result))

	fbs.SpawnItemRespStart(builder)
	fbs.SpawnItemRespAddItems(builder, itemsVector)
	respOffset := fbs.SpawnItemRespEnd(builder)
	builder.Finish(respOffset)

	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(fbs.MessageTypeSpawnItemResp, respBytes), nil
}
