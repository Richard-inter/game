package clawmachine

import (
	"context"
	"fmt"

	"github.com/Richard-inter/game/internal/cache"
	"github.com/Richard-inter/game/internal/repository"
	pb "github.com/Richard-inter/game/pkg/protocol/clawMachine_Websocket"
	fbs "github.com/Richard-inter/game/pkg/protocol/clawMachine_Websocket/clawMachine"
	flatbuffers "github.com/google/flatbuffers/go"
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

	// ---- 2. Prepare response data ----
	gameID := uint64(8)
	resultsData := []struct {
		ItemID  uint64
		Catched bool
	}{
		{1, true}, {2, false}, {3, true}, {4, false}, {5, false},
		{6, false}, {7, false}, {8, false}, {9, false}, {10, false},
	}

	// ---- 3. Build StartClawGameResp ----
	builder := flatbuffers.NewBuilder(1024)

	// Build ClawResult objects
	resultOffsets := make([]flatbuffers.UOffsetT, len(resultsData))
	for i := len(resultsData) - 1; i >= 0; i-- {
		fbs.ClawResultStart(builder)
		fbs.ClawResultAddItemId(builder, resultsData[i].ItemID)
		fbs.ClawResultAddCatched(builder, resultsData[i].Catched)
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
	fbs.StartClawGameRespAddGameId(builder, gameID)
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
