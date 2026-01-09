package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/transport/grpc"
	dto "github.com/Richard-inter/game/internal/transport/http/DTO"
	"github.com/Richard-inter/game/pkg/common"
	"github.com/Richard-inter/game/pkg/protocol/clawMachine"
	"github.com/Richard-inter/game/pkg/protocol/player"
)

type ClawMachineHandler struct {
	logger            *zap.SugaredLogger
	clawMachineClient *grpc.ClawMachineClient
}

func NewClawMachineHandler(
	logger *zap.SugaredLogger,
	grpcManager *grpc.ClientManager,
) (*ClawMachineHandler, error) {
	clawMachineClient, err := grpcManager.GetClawMachineClient()
	if err != nil {
		return nil, err
	}

	return &ClawMachineHandler{
		logger:            logger,
		clawMachineClient: clawMachineClient,
	}, nil
}

func (h *ClawMachineHandler) HandleCreateClawMachine(c *gin.Context) {
	var req dto.CreateClawMachineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	// Convert DTO to gRPC request
	grpcReq := &clawMachine.CreateClawMachineReq{
		Name:  req.Name,
		Price: req.Price,
	}

	for _, item := range req.Items {
		grpcReq.Items = append(grpcReq.Items, &clawMachine.Items{
			ItemID: item.ItemID,
		})
	}

	resp, err := h.clawMachineClient.CreateClawMachine(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to create claw machine", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully created claw machine", "machine_name", req.Name)
	common.SendCreated(c, resp.Machine)
}

func (h *ClawMachineHandler) HandleGetClawMachineInfo(c *gin.Context) {
	machineIDParam := c.Param("machineID")
	var machineID int64
	_, err := fmt.Sscan(machineIDParam, &machineID)
	if err != nil {
		h.logger.Errorw("Invalid machine ID", "error", err)
		common.SendError(c, 400, "Invalid machine ID")
		return
	}

	grpcReq := &clawMachine.GetClawMachineInfoReq{
		MachineID: machineID,
	}

	resp, err := h.clawMachineClient.GetClawMachineInfo(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to get claw machine info", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved claw machine info", "machine_id", machineID)
	common.SendSuccess(c, resp.Machine)
}

func (h *ClawMachineHandler) HandleCreateClawItems(c *gin.Context) {
	var req dto.CreateClawItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	fmt.Printf("ClawItems: %+v\n", req.ClawItems)

	// Convert DTO to gRPC request
	grpcReq := &clawMachine.CreateClawItemsReq{}
	for _, item := range req.ClawItems {
		grpcReq.ClawItems = append(grpcReq.ClawItems, &clawMachine.CreateItemReq{
			Name:   item.Name,
			Rarity: item.Rarity,
			Weight: item.Weight,
		})
	}

	resp, err := h.clawMachineClient.CreateClawItems(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to create claw items", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully created claw items", "item_count", len(req.ClawItems))
	common.SendCreated(c, resp.ClawItems)
}

func (h *ClawMachineHandler) HandleGetClawPlayerInfo(c *gin.Context) {
	playerIDParam := c.Param("playerID")
	var playerID int64
	_, err := fmt.Sscan(playerIDParam, &playerID)
	if err != nil {
		h.logger.Errorw("Invalid player ID", "error", err)
		common.SendError(c, 400, "Invalid player ID")
		return
	}

	grpcReq := &clawMachine.GetClawPlayerInfoReq{
		PlayerID: playerID,
	}

	resp, err := h.clawMachineClient.GetClawPlayerInfo(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to get claw player info", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved claw player info", "player_id", playerID)
	common.SendSuccess(c, resp.Player)
}

func (h *ClawMachineHandler) HandleCreatePlayer(c *gin.Context) {
	var req dto.CreateClawPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	// Convert DTO to gRPC request
	grpcReq := &clawMachine.CreateClawPlayerReq{
		Player: &clawMachine.ClawPlayer{
			BasePlayer: &player.Player{
				PlayerID: req.PlayerID,
				UserName: req.UserName,
			},
			Coin:    req.Coin,
			Diamond: req.Diamond,
		},
	}

	resp, err := h.clawMachineClient.CreateClawPlayer(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to create claw player", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully created claw player", "player_id", req.PlayerID)
	common.SendCreated(c, resp.Player)
}
