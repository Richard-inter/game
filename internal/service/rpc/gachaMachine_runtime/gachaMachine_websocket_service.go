package gachaMachine_runtime

import (
	"context"
	"errors"
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/cache"
	"github.com/Richard-inter/game/internal/domain"
	"github.com/Richard-inter/game/internal/repository"
	"github.com/Richard-inter/game/pkg/logger"
	pb "github.com/Richard-inter/game/pkg/protocol/gachaMachine_Websocket"
	fbs "github.com/Richard-inter/game/pkg/protocol/gachaMachine_Websocket/gachaMachine"
)

type GachaMachineWebsocketService struct {
	pb.UnimplementedGachaMachineRuntimeServiceServer
	repo      repository.GachaMachineRepository
	redis     *cache.RedisClient
	streamKey string
	log       *zap.SugaredLogger
}

func NewGachaMachineWebsocketService(repo repository.GachaMachineRepository, redis *cache.RedisClient, streamKey string) *GachaMachineWebsocketService {
	return &GachaMachineWebsocketService{
		repo:      repo,
		redis:     redis,
		streamKey: streamKey,
		log:       logger.GetSugar(),
	}
}

func (_ *GachaMachineWebsocketService) buildEnvelopeResponse(messageType fbs.MessageType, payloadBytes []byte) *pb.RuntimeResponse {
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

func (s *GachaMachineWebsocketService) GetPullResultWs(ctx context.Context, req *pb.RuntimeRequest) (*pb.RuntimeResponse, error) {
	getPullResultReq := fbs.GetRootAsGetPullResultWsReq(req.Payload, 0)
	playerID := getPullResultReq.PlayerId()
	machineID := getPullResultReq.MachineId()
	pullCount := getPullResultReq.PullCount()

	if playerID <= 0 || machineID <= 0 || pullCount <= 0 {
		return nil, errors.New("invalid player ID, machine ID, or pull count")
	}

	err := s.PlayMachine(ctx, playerID, machineID, pullCount)
	if err != nil {
		s.log.Errorf("Failed to play machine: %v", err)
		return s.createErrorResponse(err), nil
	}

	session := &domain.GachaPullSession{
		GachaMachineID: machineID,
		PlayerID:       playerID,
		PullCount:      pullCount,
	}

	if pullCount == 1 {
		itemID, err := s.PullGachaSingle(ctx, machineID, playerID)
		if err != nil {
			return nil, err
		}

		if err := s.redis.PublishGachaEvent(ctx, s.streamKey, session, itemID); err != nil {
			fmt.Println("error publish")
			s.log.Errorf("Failed to publish gacha pull history to stream: %v", err)
			// Continue even if stream publish fails
		}

		// Create the FlatBuffers response
		builder := flatbuffers.NewBuilder(256)

		// ---- Build item_ids vector (length = 1) ----
		fbs.GetPullResultWsRespStartItemIdsVector(builder, 1)
		builder.PrependInt64(itemID)
		itemIdsVector := builder.EndVector(1)

		// ---- Build response ----
		fbs.GetPullResultWsRespStart(builder)
		fbs.GetPullResultWsRespAddItemIds(builder, itemIdsVector)
		respOffset := fbs.GetPullResultWsRespEnd(builder)

		builder.Finish(respOffset)
		respBytes := builder.FinishedBytes()

		return s.buildEnvelopeResponse(
			fbs.MessageTypeGetPullResultWsResp,
			respBytes,
		), nil
	}

	itemIDs, err := s.PullGachaByMachineIDMulti(ctx, machineID, playerID, int(pullCount))
	if err != nil {
		return nil, err
	}

	// Publish all pull histories to stream with session data in a single message
	if err := s.redis.PublishGachaEvent(ctx, s.streamKey, session, itemIDs); err != nil {
		s.log.Errorf("Failed to publish gacha pull history to stream: %v", err)
		// Continue even if stream publish fails
	}

	builder := flatbuffers.NewBuilder(1024)

	itemIDsOffsets := make([]flatbuffers.UOffsetT, len(itemIDs))
	for i, id := range itemIDs {
		itemIDsOffsets[i] = fbs.GetPullResultWsRespStartItemIdsVector(builder, int(id))
	}
	fbs.GetPullResultWsRespStartItemIdsVector(builder, len(itemIDs))
	for i := len(itemIDsOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(itemIDsOffsets[i])
	}
	itemIDsVector := builder.EndVector(len(itemIDsOffsets))

	fbs.GetPullResultWsRespStart(builder)
	fbs.GetPullResultWsRespAddItemIds(builder, itemIDsVector)
	respOffset := fbs.GetPullResultWsRespEnd(builder)

	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(fbs.MessageTypeGetPullResultWsResp, respBytes), nil
}

func (s *GachaMachineWebsocketService) GetPlayerInfoWs(
	ctx context.Context,
	req *pb.RuntimeRequest,
) (*pb.RuntimeResponse, error) {

	getPlayerInfoReq := fbs.GetRootAsGetPlayerInfoWsReq(req.Payload, 0)
	playerID := getPlayerInfoReq.PlayerId()

	if playerID <= 0 {
		return nil, errors.New("invalid player ID")
	}

	resp, err := s.repo.GetGachaPlayerInfo(ctx, playerID)
	if err != nil {
		s.log.Errorf("Failed to get player info: %v", err)
		return s.createErrorResponse(err), nil
	}

	builder := flatbuffers.NewBuilder(256)

	// ---- Strings ----
	usernameOffset := builder.CreateString(resp.Player.UserName)

	// ---- Build response ----
	fbs.GetPlayerInfoWsRespStart(builder)
	fbs.GetPlayerInfoWsRespAddPlayerId(builder, resp.Player.ID)
	fbs.GetPlayerInfoWsRespAddUsername(builder, usernameOffset)
	fbs.GetPlayerInfoWsRespAddCoin(builder, resp.Coin)
	fbs.GetPlayerInfoWsRespAddDiamond(builder, resp.Diamond)
	respOffset := fbs.GetPlayerInfoWsRespEnd(builder)

	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(
		fbs.MessageTypeGetPlayerInfoWsResp,
		respBytes,
	), nil
}

func (s *GachaMachineWebsocketService) GetMachineInfoWs(
	ctx context.Context,
	req *pb.RuntimeRequest,
) (*pb.RuntimeResponse, error) {
	getMachineInfoReq := fbs.GetRootAsGetMachineInfoWsReq(req.Payload, 0)
	machineID := getMachineInfoReq.MachineId()

	if machineID <= 0 {
		return nil, errors.New("invalid machine ID")
	}

	resp, err := s.repo.GetGachaMachineInfo(ctx, machineID)
	if err != nil {
		s.log.Errorf("Failed to get machine info: %v", err)
		return s.createErrorResponse(err), nil
	}

	builder := flatbuffers.NewBuilder(2048)

	// ---- Build Items vector ----
	itemOffsets := make([]flatbuffers.UOffsetT, len(resp.Items))
	for i, item := range resp.Items {
		nameOffset := builder.CreateString(item.Item.Name)
		rarityOffset := builder.CreateString(item.Item.Rarity)

		fbs.ItemsStart(builder)
		fbs.ItemsAddItemId(builder, item.ItemID)
		fbs.ItemsAddName(builder, nameOffset)
		fbs.ItemsAddRarity(builder, rarityOffset)
		fbs.ItemsAddPullWeight(builder, item.Item.PullWeight)
		itemOffsets[i] = fbs.ItemsEnd(builder)
	}

	fbs.GetMachineInfoWsRespStartItemsVector(builder, len(itemOffsets))
	for i := len(itemOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(itemOffsets[i])
	}
	itemsVector := builder.EndVector(len(itemOffsets))

	// ---- Strings ----
	nameOffset := builder.CreateString(resp.Name)

	// ---- Build response ----
	fbs.GetMachineInfoWsRespStart(builder)
	fbs.GetMachineInfoWsRespAddMachineId(builder, resp.ID)
	fbs.GetMachineInfoWsRespAddName(builder, nameOffset)
	fbs.GetMachineInfoWsRespAddPrice(builder, resp.Price)
	fbs.GetMachineInfoWsRespAddPriceTimesTen(builder, resp.PriceTimesTen)
	fbs.GetMachineInfoWsRespAddSuperRarePity(builder, resp.SuperRarePity)
	fbs.GetMachineInfoWsRespAddUltraRarePity(builder, resp.UltraRarePity)
	fbs.GetMachineInfoWsRespAddItems(builder, itemsVector)
	respOffset := fbs.GetMachineInfoWsRespEnd(builder)

	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(
		fbs.MessageTypeGetMachineInfoWsResp,
		respBytes,
	), nil
}

func (s *GachaMachineWebsocketService) createErrorResponse(err error) *pb.RuntimeResponse {
	builder := flatbuffers.NewBuilder(256)

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
