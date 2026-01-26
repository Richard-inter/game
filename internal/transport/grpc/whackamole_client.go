package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	whackAMolepb "github.com/Richard-inter/game/pkg/protocol/whackAMole"
)

type WhackAMoleClient struct {
	client whackAMolepb.WhackAMoleServiceClient
	conn   *grpc.ClientConn
}

func NewWhackAMoleClient(address string) (*WhackAMoleClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to whack a mole service: %w", err)
	}

	return &WhackAMoleClient{
		client: whackAMolepb.NewWhackAMoleServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *WhackAMoleClient) CreateWhackAMolePlayer(ctx context.Context, req *whackAMolepb.CreateWhackAMolePlayerReq) (*whackAMolepb.CreateWhackAMolePlayerResp, error) {
	return c.client.CreateWhackAMolePlayer(ctx, req)
}

func (c *WhackAMoleClient) GetPlayerInfo(ctx context.Context, req *whackAMolepb.GetPlayerInfoReq) (*whackAMolepb.GetPlayerInfoResp, error) {
	return c.client.GetPlayerInfo(ctx, req)
}

func (c *WhackAMoleClient) GetLeaderboard(ctx context.Context, req *whackAMolepb.GetLeaderboardReq) (*whackAMolepb.GetLeaderboardResp, error) {
	return c.client.GetLeaderboard(ctx, req)
}

func (c *WhackAMoleClient) GetMoleWeightConfig(ctx context.Context, req *whackAMolepb.GetMoleWeightConfigReq) (*whackAMolepb.GetMoleWeightConfigResp, error) {
	return c.client.GetMoleWeightConfig(ctx, req)
}

func (c *WhackAMoleClient) UpdateScore(ctx context.Context, req *whackAMolepb.UpdateScoreReq) (*whackAMolepb.UpdateScoreResp, error) {
	return c.client.UpdateScore(ctx, req)
}

func (c *WhackAMoleClient) CreateMoleWeightConfig(ctx context.Context, req *whackAMolepb.CreateMoleWeightConfigReq) (*whackAMolepb.CreateMoleWeightConfigResp, error) {
	return c.client.CreateMoleWeightConfig(ctx, req)
}

func (c *WhackAMoleClient) UpdateMoleWeightConfig(ctx context.Context, req *whackAMolepb.UpdateMoleWeightConfigReq) (*whackAMolepb.UpdateMoleWeightConfigResp, error) {
	return c.client.UpdateMoleWeightConfig(ctx, req)
}

func (c *WhackAMoleClient) Close() error {
	return c.conn.Close()
}
