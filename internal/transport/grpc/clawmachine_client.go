package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	clawmachinepb "github.com/Richard-inter/game/pkg/protocol/clawMachine"
)

type ClawMachineClient struct {
	client clawmachinepb.ClawMachineServiceClient
	conn   *grpc.ClientConn
}

func NewClawMachineClient(address string) (*ClawMachineClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to clawmachine service: %w", err)
	}

	return &ClawMachineClient{
		client: clawmachinepb.NewClawMachineServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *ClawMachineClient) GetClawPlayerInfo(ctx context.Context, playerID int64) (*clawmachinepb.GetClawPlayerInfoResp, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &clawmachinepb.GetClawPlayerInfoReq{
		PlayerID: playerID,
	}

	return c.client.GetClawPlayerInfo(ctx, req)
}

func (c *ClawMachineClient) StartClawGame(ctx context.Context, machineID int64, playerID int64) (*clawmachinepb.StartClawGameResp, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &clawmachinepb.StartClawGameReq{
		MachineID: machineID,
		PlayerID:  playerID,
	}

	return c.client.StartClawGame(ctx, req)
}

func (c *ClawMachineClient) Close() error {
	return c.conn.Close()
}
