package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gachamachinepb "github.com/Richard-inter/game/pkg/protocol/gachaMachine"
)

type GachaMachineClient struct {
	client gachamachinepb.GachaMachineServiceClient
	conn   *grpc.ClientConn
}

func NewGachaMachineClient(address string) (*GachaMachineClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gachamachine service: %w", err)
	}

	return &GachaMachineClient{
		client: gachamachinepb.NewGachaMachineServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GachaMachineClient) CreateGachaMachine(ctx context.Context, req *gachamachinepb.CreateGachaMachineReq) (*gachamachinepb.CreateGachaMachineResp, error) {
	return c.client.CreateGachaMachine(ctx, req)
}

func (c *GachaMachineClient) GetGachaMachineInfo(ctx context.Context, req *gachamachinepb.GetGachaMachineInfoReq) (*gachamachinepb.GetGachaMachineInfoResp, error) {
	return c.client.GetGachaMachineInfo(ctx, req)
}

func (c *GachaMachineClient) CreateGachaItems(ctx context.Context, req *gachamachinepb.CreateGachaItemsReq) (*gachamachinepb.CreateGachaItemsResp, error) {
	return c.client.CreateGachaItems(ctx, req)
}

func (c *GachaMachineClient) CreateGachaPlayer(ctx context.Context, req *gachamachinepb.CreateGachaPlayerReq) (*gachamachinepb.CreateGachaPlayerResp, error) {
	return c.client.CreateGachaPlayer(ctx, req)
}

func (c *GachaMachineClient) GetGachaPlayerInfo(ctx context.Context, req *gachamachinepb.GetGachaPlayerInfoReq) (*gachamachinepb.GetGachaPlayerInfoResp, error) {
	return c.client.GetGachaPlayerInfo(ctx, req)
}

func (c *GachaMachineClient) AdjustPlayerCoin(ctx context.Context, req *gachamachinepb.AdjustPlayerCoinReq) (*gachamachinepb.AdjustPlayerCoinResp, error) {
	return c.client.AdjustPlayerCoin(ctx, req)
}

func (c *GachaMachineClient) AdjustPlayerDiamond(ctx context.Context, req *gachamachinepb.AdjustPlayerDiamondReq) (*gachamachinepb.AdjustPlayerDiamondResp, error) {
	return c.client.AdjustPlayerDiamond(ctx, req)
}

func (c *GachaMachineClient) GetPullResult(ctx context.Context, req *gachamachinepb.GetPullResultReq) (*gachamachinepb.GetPullResultResp, error) {
	return c.client.GetPullResult(ctx, req)
}

func (c *GachaMachineClient) GetPullTimesTenResult(ctx context.Context, req *gachamachinepb.GetPullTimesTenResultReq) (*gachamachinepb.GetPullTimesTenResultResp, error) {
	return c.client.GetPullTimesTenResult(ctx, req)
}

func (c *GachaMachineClient) Close() error {
	return c.conn.Close()
}
