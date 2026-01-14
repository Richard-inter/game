package clawmachine

import (
	"context"
	"fmt"

	"github.com/Richard-inter/game/internal/cache"
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

func (s *ClawMachineWebsocketService) StartClawGameWs(ctx context.Context, req *pb.RuntimeRequest) (*pb.RuntimeResponse, error) {
	// Get the payload bytes from the request
	payload := req.GetPayload()

	// Parse the StartClawGame payload from FlatBuffer
	envelope := fbs.GetRootAsEnvelope(payload, 0)

	// Get the payload bytes from envelope
	payloadBytes := envelope.PayloadBytes()
	if len(payloadBytes) > 0 {
		// Parse as StartClawGameReq
		startReq := fbs.GetRootAsStartClawGameReq(payloadBytes, 0)
		playerID := startReq.PlayerId()
		machineID := startReq.MachineId()

		fmt.Printf("Parsed StartClawGame payload:\n")
		fmt.Printf("  PlayerID: %d\n", playerID)
		fmt.Printf("  MachineID: %d\n", machineID)
	} else {
		fmt.Printf("Empty payload in envelope\n")
	}

	return &pb.RuntimeResponse{
		Payload: []byte("StartClawGame request received"),
	}, nil
}
