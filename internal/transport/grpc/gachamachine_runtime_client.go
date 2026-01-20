package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	runtimepb "github.com/Richard-inter/game/pkg/protocol/gachaMachine_Websocket"
)

type GachaMachineRuntimeClient struct {
	client runtimepb.GachaMachineRuntimeServiceClient
	conn   *grpc.ClientConn
}

func NewGachaMachineRuntimeClient(address string) (*GachaMachineRuntimeClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gachamachine runtime service: %w", err)
	}

	return &GachaMachineRuntimeClient{
		client: runtimepb.NewGachaMachineRuntimeServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GachaMachineRuntimeClient) GetPullResultWs(ctx context.Context, req *runtimepb.RuntimeRequest) (*runtimepb.RuntimeResponse, error) {
	return c.client.GetPullResultWs(ctx, req)
}

func (c *GachaMachineRuntimeClient) GetPlayerInfoWs(ctx context.Context, req *runtimepb.RuntimeRequest) (*runtimepb.RuntimeResponse, error) {
	return c.client.GetPlayerInfoWs(ctx, req)
}

func (c *GachaMachineRuntimeClient) GetMachineInfoWs(ctx context.Context, req *runtimepb.RuntimeRequest) (*runtimepb.RuntimeResponse, error) {
	return c.client.GetMachineInfoWs(ctx, req)
}
