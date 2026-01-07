package clawmachine

import (
	"github.com/Richard-inter/game/pkg/protocol/clawMachine"
)

// PlayerGRPCService implements the PlayerService gRPC service
type ClawMachineGRPCServices struct {
	clawMachine.UnimplementedClawMachineServiceServer
}

// NewClawMachineGRPCService creates a new ClawMachineGRPCService
func NewClawMachineGRPCService() *ClawMachineGRPCServices {
	return &ClawMachineGRPCServices{}
}
