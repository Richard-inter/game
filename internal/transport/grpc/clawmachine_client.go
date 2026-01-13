package grpc

import (
	"context"
	"fmt"

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

func (c *ClawMachineClient) GetClawPlayerInfo(ctx context.Context, req *clawmachinepb.GetClawPlayerInfoReq) (*clawmachinepb.GetClawPlayerInfoResp, error) {
	return c.client.GetClawPlayerInfo(ctx, req)
}

func (c *ClawMachineClient) StartClawGame(ctx context.Context, req *clawmachinepb.StartClawGameReq) (*clawmachinepb.StartClawGameResp, error) {
	return c.client.StartClawGame(ctx, req)
}

func (c *ClawMachineClient) GetClawMachineInfo(ctx context.Context, req *clawmachinepb.GetClawMachineInfoReq) (*clawmachinepb.GetClawMachineInfoResp, error) {
	return c.client.GetClawMachineInfo(ctx, req)
}

func (c *ClawMachineClient) CreateClawMachine(ctx context.Context, req *clawmachinepb.CreateClawMachineReq) (*clawmachinepb.CreateClawMachineResp, error) {
	return c.client.CreateClawMachine(ctx, req)
}

func (c *ClawMachineClient) CreateClawItems(ctx context.Context, req *clawmachinepb.CreateClawItemsReq) (*clawmachinepb.CreateClawItemsResp, error) {
	return c.client.CreateClawItems(ctx, req)
}

func (c *ClawMachineClient) CreateClawPlayer(ctx context.Context, req *clawmachinepb.CreateClawPlayerReq) (*clawmachinepb.CreateClawPlayerResp, error) {
	return c.client.CreateClawPlayer(ctx, req)
}

func (c *ClawMachineClient) AdjustPlayerCoin(ctx context.Context, req *clawmachinepb.AdjustPlayerCoinReq) (*clawmachinepb.AdjustPlayerCoinResp, error) {
	return c.client.AdjustPlayerCoin(ctx, req)
}

func (c *ClawMachineClient) AdjustPlayerDiamond(ctx context.Context, req *clawmachinepb.AdjustPlayerDiamondReq) (*clawmachinepb.AdjustPlayerDiamondResp, error) {
	return c.client.AdjustPlayerDiamond(ctx, req)
}

func (c *ClawMachineClient) Close() error {
	return c.conn.Close()
}
