package clawmachine_runtime

import (
	"context"
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/1nterdigital/game/internal/cache"
	"github.com/1nterdigital/game/internal/domain"
	"github.com/1nterdigital/game/internal/repository"
	"github.com/1nterdigital/game/pkg/constant"
	"github.com/1nterdigital/game/pkg/logger"
	pb "github.com/1nterdigital/game/pkg/protocol/clawMachine_Websocket"
	fbs "github.com/1nterdigital/game/pkg/protocol/clawMachine_Websocket/clawMachine"
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
	builder := flatbuffers.NewBuilder(len(payloadBytes) + constant.Byte256)
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

	results, err := s.PreDetermineCatchResults(ctx, machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to pre-determine catch results: %w", err)
	}

	err = s.PlayMachine(ctx, playerID, machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to charge player: %w", err)
	}

	gameID, err := s.repo.AddGameHistory(ctx, &domain.ClawMachineGameRecord{
		PlayerID:      playerID,
		ClawMachineID: machineID,
		TouchedItemID: nil,
		CreatedBy:     fmt.Sprintf("%d", playerID),
		UpdatedBy:     fmt.Sprintf("%d", playerID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create game history: %w", err)
	}

	err = s.redis.StoreGameResults(ctx, gameID, results)
	if err != nil {
		// Log error but don't fail the request
		logger.GetSugar().Warnf("failed to store game results in Redis: %v", err)
	}

	builder := flatbuffers.NewBuilder(constant.Byte1024)
	resultOffsets := make([]flatbuffers.UOffsetT, len(results))
	for i := len(results) - 1; i >= 0; i-- {
		fbs.ClawResultStart(builder)
		fbs.ClawResultAddItemId(builder, results[i].ItemID)
		fbs.ClawResultAddCatched(builder, results[i].Success)
		resultOffsets[i] = fbs.ClawResultEnd(builder)
	}

	fbs.StartClawGameRespStartResultsVector(builder, len(resultOffsets))
	for i := len(resultOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(resultOffsets[i])
	}
	resultsVector := builder.EndVector(len(resultOffsets))

	fbs.StartClawGameRespStart(builder)
	fbs.StartClawGameRespAddGameId(builder, gameID)
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

	domainPlayer, err := s.repo.GetClawPlayerInfo(ctx, playerID)
	if err != nil {
		return nil, err
	}

	builder := flatbuffers.NewBuilder(constant.Byte1024)
	usernameOffset := builder.CreateString(domainPlayer.Player.UserName)

	fbs.GetPlayerInfoWsRespStart(builder)

	fbs.GetPlayerInfoWsRespAddPlayerId(builder, domainPlayer.Player.ID)
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
	if err := s.redis.GetGameResults(ctx, gameID, &storedResults); err != nil {
		return nil, fmt.Errorf("failed to load game results: %w", err)
	}

	var serverCatched bool
	if itemID != 0 {
		var serverResult *CatchResult
		for i := range storedResults {
			if storedResults[i].ItemID == itemID {
				serverResult = &storedResults[i]
				break
			}
		}

		if serverResult == nil {
			return nil, fmt.Errorf("item %d not found in game %d", itemID, gameID)
		}

		serverCatched = serverResult.Success

		if serverCatched != catched {
			logger.GetSugar().Warnf(
				"client desync - game=%d item=%d client=%t server=%t",
				gameID, itemID, catched, serverCatched,
			)
		}
	}

	var itemIDPtr *int64
	itemIDPtr = &itemID
	if itemID == 0 {
		itemIDPtr = nil
		serverCatched = false
	}

	game, err := s.repo.AddTouchedItemRecord(ctx, gameID, itemIDPtr, serverCatched)
	if err != nil {
		return nil, fmt.Errorf("failed to persist record: %w", err)
	}

	err = s.repo.AddItemInventory(ctx, game.PlayerID, itemID, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to player inventory: %w", err)
	}

	payout := int64(0)
	if serverCatched {
		payout = GetRarityValue(game.TouchedItem.Rarity, game.Machine.Price)
	}

	err = s.repo.UpdateClawMachineRTP(ctx, game.ClawMachineID, game.TouchedItem.MaxItemSpawned, payout)
	if err != nil {
		return nil, fmt.Errorf("failed to update claw machine RTP: %w", err)
	}

	_ = s.redis.DeleteGameResults(ctx, gameID)

	builder := flatbuffers.NewBuilder(constant.Byte256)

	fbs.AddTouchedItemRecordRespStart(builder)
	fbs.AddTouchedItemRecordRespAddGameId(builder, gameID)
	fbs.AddTouchedItemRecordRespAddItemId(builder, itemID)
	fbs.AddTouchedItemRecordRespAddCatched(builder, serverCatched)
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

	result, err := s.SpawnMachineItems(ctx, machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to spawn item: %w", err)
	}

	items := make([]uint64, len(result))
	for i, v := range result {
		items[i] = uint64(v)
	}

	builder := flatbuffers.NewBuilder(constant.Byte256)

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

func (s *ClawMachineWebsocketService) GetPlayerInventoryWs(ctx context.Context, req *pb.RuntimeRequest) (*pb.RuntimeResponse, error) {
	getPlayerInventoryReq := fbs.GetRootAsGetPlayerInventoryWsReq(req.Payload, 0)
	playerID := getPlayerInventoryReq.PlayerId()

	if playerID <= 0 {
		return nil, fmt.Errorf("invalid player ID")
	}

	inventory, err := s.repo.GetPlayerInventory(ctx, playerID)
	if err != nil {
		logger.GetSugar().Errorf("Failed to get player inventory: %v", err)
		return s.createErrorResponse(err), nil
	}

	builder := flatbuffers.NewBuilder(constant.Byte2048)

	// Build PlayerInventoryItem vector
	itemOffsets := make([]flatbuffers.UOffsetT, len(inventory))
	for i := range inventory {
		item := &inventory[i]
		fbs.PlayerInventoryItemStart(builder)
		fbs.PlayerInventoryItemAddItemId(builder, item.ItemID)
		fbs.PlayerInventoryItemAddQuantity(builder, item.Quantity)
		itemOffsets[i] = fbs.PlayerInventoryItemEnd(builder)
	}

	fbs.GetPlayerInventoryWsRespStartItemsVector(builder, len(itemOffsets))
	for i := len(itemOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(itemOffsets[i])
	}
	itemsVectorOffset := builder.EndVector(len(itemOffsets))

	fbs.GetPlayerInventoryWsRespStart(builder)
	fbs.GetPlayerInventoryWsRespAddPlayerId(builder, playerID)
	fbs.GetPlayerInventoryWsRespAddItems(builder, itemsVectorOffset)
	respOffset := fbs.GetPlayerInventoryWsRespEnd(builder)

	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(fbs.MessageTypeGetPlayerInventoryWsResp, respBytes), nil
}

func (s *ClawMachineWebsocketService) GetGameHistoryWs(ctx context.Context, req *pb.RuntimeRequest) (*pb.RuntimeResponse, error) {
	getGameHistoryReq := fbs.GetRootAsGetGameHistoryWsReq(req.Payload, 0)
	playerID := getGameHistoryReq.PlayerId()

	if playerID <= 0 {
		return nil, fmt.Errorf("invalid player ID")
	}

	gameRecords, err := s.repo.GetGameHistory(ctx, playerID)
	if err != nil {
		logger.GetSugar().Errorf("Failed to get game history: %v", err)
		return s.createErrorResponse(err), nil
	}

	builder := flatbuffers.NewBuilder(constant.Byte2048)

	// Build GameRecord vector
	recordOffsets := make([]flatbuffers.UOffsetT, len(gameRecords))
	for i, record := range gameRecords {
		createdAtOffset := builder.CreateString(record.CreatedAt.Format("2006-01-02 15:04:05"))

		fbs.GameRecordStart(builder)
		fbs.GameRecordAddGameId(builder, record.ID)
		fbs.GameRecordAddPlayerId(builder, record.PlayerID)
		fbs.GameRecordAddMachineId(builder, record.ClawMachineID)

		var touchedItemID int64 = 0
		if record.TouchedItemID != nil {
			touchedItemID = *record.TouchedItemID
		}
		fbs.GameRecordAddTouchedItemId(builder, touchedItemID)

		var catched bool = false
		if record.Catched != nil {
			catched = *record.Catched
		}
		fbs.GameRecordAddCatched(builder, catched)
		fbs.GameRecordAddCreatedAt(builder, createdAtOffset)
		recordOffsets[i] = fbs.GameRecordEnd(builder)
	}

	fbs.GetGameHistoryWsRespStartRecordsVector(builder, len(recordOffsets))
	for i := len(recordOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(recordOffsets[i])
	}
	recordsVectorOffset := builder.EndVector(len(recordOffsets))

	fbs.GetGameHistoryWsRespStart(builder)
	fbs.GetGameHistoryWsRespAddPlayerId(builder, playerID)
	fbs.GetGameHistoryWsRespAddRecords(builder, recordsVectorOffset)
	respOffset := fbs.GetGameHistoryWsRespEnd(builder)

	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(fbs.MessageTypeGetGameHistoryWsResp, respBytes), nil
}

func (s *ClawMachineWebsocketService) createErrorResponse(err error) *pb.RuntimeResponse {
	builder := flatbuffers.NewBuilder(constant.Byte256)

	// Create error message string
	errorMsgOffset := builder.CreateString(err.Error())

	// Create ErrorResp
	fbs.ErrorRespStart(builder)
	fbs.ErrorRespAddCode(builder, 1)
	fbs.ErrorRespAddMessage(builder, errorMsgOffset)
	errorRespOffset := fbs.ErrorRespEnd(builder)

	builder.Finish(errorRespOffset)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(fbs.MessageTypeErrorResp, respBytes)
}
