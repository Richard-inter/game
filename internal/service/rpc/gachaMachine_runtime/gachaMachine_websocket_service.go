package gachaMachine_runtime

import (
	"context"
	"errors"

	flatbuffers "github.com/google/flatbuffers/go"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/pkg/logger"
	"github.com/Richard-inter/game/pkg/protocol/gachaMachine"
	pb "github.com/Richard-inter/game/pkg/protocol/gachaMachine_Websocket"
	fbs "github.com/Richard-inter/game/pkg/protocol/gachaMachine_Websocket/gachaMachine"
)

type GachaMachineWebsocketService struct {
	pb.UnimplementedGachaMachineRuntimeServiceServer
	gachaService gachaMachine.GachaMachineServiceServer
	log          *zap.SugaredLogger
}

func NewGachaMachineWebsocketService(gachaService gachaMachine.GachaMachineServiceServer) *GachaMachineWebsocketService {
	return &GachaMachineWebsocketService{
		gachaService: gachaService,
		log:          logger.GetSugar(),
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

	s.log.Infof("GetPullResultWs received")
	s.log.Infof("  PlayerID : %d", playerID)
	s.log.Infof("  MachineID: %d", machineID)
	s.log.Infof("  PullCount: %d", pullCount)

	if playerID <= 0 || machineID <= 0 || pullCount <= 0 {
		return nil, errors.New("invalid player ID, machine ID, or pull count")
	}

	// Create the gRPC request
	grpcReq := &gachaMachine.GetPullResultReq{
		PlayerID:  playerID,
		MachineID: machineID,
		PullCount: int32(pullCount),
	}

	// Call the gRPC service
	grpcResp, err := s.gachaService.GetPullResult(ctx, grpcReq)
	if err != nil {
		s.log.Errorf("Failed to get pull result: %v", err)
		return s.createErrorResponse(err), nil
	}

	// Create the FlatBuffers response
	builder := flatbuffers.NewBuilder(1024)

	// Create item IDs vector
	fbs.GetPullResultWsRespStartItemIdsVector(builder, len(grpcResp.ItemIDs))
	for i := len(grpcResp.ItemIDs) - 1; i >= 0; i-- {
		builder.PrependInt64(grpcResp.ItemIDs[i])
	}
	itemIdsVector := builder.EndVector(len(grpcResp.ItemIDs))

	fbs.GetPullResultWsRespStart(builder)
	fbs.GetPullResultWsRespAddItemIds(builder, itemIdsVector)
	respOffset := fbs.GetPullResultWsRespEnd(builder)

	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(fbs.MessageTypeGetPullResultWs, respBytes), nil
}

func (s *GachaMachineWebsocketService) GetPlayerInfoWs(ctx context.Context, req *pb.RuntimeRequest) (*pb.RuntimeResponse, error) {
	getPlayerInfoReq := fbs.GetRootAsGetPlayerInfoWsReq(req.Payload, 0)
	playerID := getPlayerInfoReq.PlayerId()

	s.log.Infof("GetPlayerInfoWs received")
	s.log.Infof("  PlayerID: %d", playerID)

	if playerID <= 0 {
		return nil, errors.New("invalid player ID")
	}

	// Create the gRPC request
	grpcReq := &gachaMachine.GetGachaPlayerInfoReq{
		PlayerID: playerID,
	}

	// Call the gRPC service
	grpcResp, err := s.gachaService.GetGachaPlayerInfo(ctx, grpcReq)
	if err != nil {
		s.log.Errorf("Failed to get player info: %v", err)
		return s.createErrorResponse(err), nil
	}

	// Create the FlatBuffers response
	builder := flatbuffers.NewBuilder(1024)

	// Create strings
	usernameOffset := builder.CreateString(grpcResp.Player.BasePlayer.UserName)

	fbs.GetPlayerInfoWsRespStart(builder)
	fbs.GetPlayerInfoWsRespAddPlayerId(builder, grpcResp.Player.BasePlayer.PlayerID)
	fbs.GetPlayerInfoWsRespAddUsername(builder, usernameOffset)
	fbs.GetPlayerInfoWsRespAddCoin(builder, grpcResp.Player.Coin)
	fbs.GetPlayerInfoWsRespAddDiamond(builder, grpcResp.Player.Diamond)
	respOffset := fbs.GetPlayerInfoWsRespEnd(builder)

	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(fbs.MessageTypeGetPlayerInfoWs, respBytes), nil
}

func (s *GachaMachineWebsocketService) GetMachineInfoWs(ctx context.Context, req *pb.RuntimeRequest) (*pb.RuntimeResponse, error) {
	getMachineInfoReq := fbs.GetRootAsGetMachineInfoWsReq(req.Payload, 0)
	machineID := getMachineInfoReq.MachineId()

	s.log.Infof("GetMachineInfoWs received")
	s.log.Infof("  MachineID: %d", machineID)

	if machineID <= 0 {
		return nil, errors.New("invalid machine ID")
	}

	// Create the gRPC request
	grpcReq := &gachaMachine.GetGachaMachineInfoReq{
		MachineID: machineID,
	}

	// Call the gRPC service
	grpcResp, err := s.gachaService.GetGachaMachineInfo(ctx, grpcReq)
	if err != nil {
		s.log.Errorf("Failed to get machine info: %v", err)
		return s.createErrorResponse(err), nil
	}

	// Create the FlatBuffers response
	builder := flatbuffers.NewBuilder(1024)

	if len(grpcResp.Machine) == 0 {
		return nil, errors.New("machine not found")
	}

	machine := grpcResp.Machine[0] // Take the first machine

	// Create strings
	nameOffset := builder.CreateString(machine.Name)

	fbs.GetMachineInfoWsRespStart(builder)
	fbs.GetMachineInfoWsRespAddMachineId(builder, machine.MachineID)
	fbs.GetMachineInfoWsRespAddName(builder, nameOffset)
	fbs.GetMachineInfoWsRespAddPrice(builder, machine.Price)
	fbs.GetMachineInfoWsRespAddPriceTimesTen(builder, machine.PriceTimesTen)
	fbs.GetMachineInfoWsRespAddSuperRarePity(builder, machine.SuperRarePity)
	fbs.GetMachineInfoWsRespAddUltraRarePity(builder, machine.UltraRarePity)
	respOffset := fbs.GetMachineInfoWsRespEnd(builder)

	builder.Finish(respOffset)
	respBytes := builder.FinishedBytes()

	return s.buildEnvelopeResponse(fbs.MessageTypeGetMachineInfoWs, respBytes), nil
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
