package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	runtimepb "github.com/Richard-inter/game/pkg/protocol/clawMachine_Websocket"
)

type ClawMachineRuntimeClient struct {
	client runtimepb.ClawMachineRuntimeServiceClient
	conn   *grpc.ClientConn
}

func NewClawMachineRuntimeClient(address string) (*ClawMachineRuntimeClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to clawmachine runtime service: %w", err)
	}

	return &ClawMachineRuntimeClient{
		client: runtimepb.NewClawMachineRuntimeServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *ClawMachineRuntimeClient) StartClawGameWs(ctx context.Context, req *runtimepb.RuntimeRequest) (*runtimepb.RuntimeResponse, error) {
	return c.client.StartClawGameWs(ctx, req)
}

func (c *ClawMachineRuntimeClient) AddTouchedItemRecordWs(ctx context.Context, req *runtimepb.RuntimeRequest) (*runtimepb.RuntimeResponse, error) {
	return c.client.AddTouchedItemRecordWs(ctx, req)
}

func (c *ClawMachineRuntimeClient) GetPlayerSnapshotWs(ctx context.Context, req *runtimepb.RuntimeRequest) (*runtimepb.RuntimeResponse, error) {
	return c.client.GetPlayerSnapshotWs(ctx, req)
}

func (c *ClawMachineRuntimeClient) GetMachineSnapshotWs(ctx context.Context, req *runtimepb.RuntimeRequest) (*runtimepb.RuntimeResponse, error) {
	return c.client.GetMachineSnapshotWs(ctx, req)
}
