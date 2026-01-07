package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	clawmachinepb "github.com/Richard-inter/game/pkg/protocol/clawMachine"
	playerpb "github.com/Richard-inter/game/pkg/protocol/player"
)

type Client struct {
	playerClient      playerpb.PlayerServiceClient
	clawmachineClient clawmachinepb.ClawMachineServiceClient
	conn              *grpc.ClientConn
}

type Config struct {
	PlayerServiceAddr      string
	ClawMachineServiceAddr string
}

func NewClient(cfg *Config) (*Client, error) {
	// Connect to player service
	playerConn, err := grpc.Dial(cfg.PlayerServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Connect to clawmachine service
	clawmachineConn, err := grpc.Dial(cfg.ClawMachineServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		playerConn.Close()
		return nil, err
	}

	return &Client{
		playerClient:      playerpb.NewPlayerServiceClient(playerConn),
		clawmachineClient: clawmachinepb.NewClawMachineServiceClient(clawmachineConn),
		conn:              playerConn, // Store one connection for closing
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetPlayerInfo(ctx context.Context, playerID int64) (*playerpb.GetPlayerInfoResp, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &playerpb.GetPlayerInfoReq{
		PlayerID: playerID,
	}

	return c.playerClient.GetPlayerInfo(ctx, req)
}

func (c *Client) GetClawPlayerInfo(ctx context.Context, playerID int64) (*clawmachinepb.GetClawPlayerInfoResp, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &clawmachinepb.GetClawPlayerInfoReq{
		PlayerID: playerID,
	}

	return c.clawmachineClient.GetClawPlayerInfo(ctx, req)
}

func (c *Client) StartClawGame(ctx context.Context, machineID int64, playerID int64) (*clawmachinepb.StartClawGameResp, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &clawmachinepb.StartClawGameReq{
		MachineID: machineID,
		PlayerID:  playerID,
	}

	return c.clawmachineClient.StartClawGame(ctx, req)
}
