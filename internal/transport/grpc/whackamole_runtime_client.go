package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	runtimepb "github.com/Richard-inter/game/pkg/protocol/whackAMole_Websocket"
)

type WhackAMoleRuntimeClient struct {
	client runtimepb.WhackAMoleRuntimeServiceClient
	conn   *grpc.ClientConn
}

func NewWhackAMoleRuntimeClient(address string) (*WhackAMoleRuntimeClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to whackamole runtime service: %w", err)
	}

	return &WhackAMoleRuntimeClient{
		client: runtimepb.NewWhackAMoleRuntimeServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *WhackAMoleRuntimeClient) GetMoleWeight(ctx context.Context, req *runtimepb.RuntimeRequest) (*runtimepb.RuntimeResponse, error) {
	return c.client.GetMoleWeight(ctx, req)
}

func (c *WhackAMoleRuntimeClient) GetLeaderboard(ctx context.Context, req *runtimepb.RuntimeRequest) (*runtimepb.RuntimeResponse, error) {
	return c.client.GetLeaderboard(ctx, req)
}

func (c *WhackAMoleRuntimeClient) Close() error {
	return c.conn.Close()
}
