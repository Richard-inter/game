package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	playerpb "github.com/Richard-inter/game/pkg/protocol/player"
)

type PlayerClient struct {
	client playerpb.PlayerServiceClient
	conn   *grpc.ClientConn
}

func NewPlayerClient(address string) (*PlayerClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to player service: %w", err)
	}

	return &PlayerClient{
		client: playerpb.NewPlayerServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *PlayerClient) GetPlayerInfo(ctx context.Context, req *playerpb.GetPlayerInfoReq) (*playerpb.GetPlayerInfoResp, error) {
	return c.client.GetPlayerInfo(ctx, req)
}

func (c *PlayerClient) CreatePlayer(ctx context.Context, req *playerpb.CreatePlayerReq) (*playerpb.CreatePlayerResp, error) {
	return c.client.CreatePlayer(ctx, req)
}

func (c *PlayerClient) Close() error {
	return c.conn.Close()
}
