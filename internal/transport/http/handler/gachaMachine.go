package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/transport/grpc"
	dto "github.com/Richard-inter/game/internal/transport/http/DTO"

	// dto "github.com/Richard-inter/game/internal/transport/http/DTO"
	"github.com/Richard-inter/game/pkg/common"
	"github.com/Richard-inter/game/pkg/protocol/gachaMachine"
	"github.com/Richard-inter/game/pkg/protocol/player"
)

type GachaMachineHandler struct {
	logger             *zap.SugaredLogger
	gachaMachineClient *grpc.GachaMachineClient
}

func NewGachaMachineHandler(
	logger *zap.SugaredLogger,
	grpcManager *grpc.ClientManager,
) (*GachaMachineHandler, error) {
	gachaMachineClient, err := grpcManager.GetGachaMachineClient()
	if err != nil {
		return nil, err
	}

	return &GachaMachineHandler{
		logger:             logger,
		gachaMachineClient: gachaMachineClient,
	}, nil
}

// HandleCreateGachaPlayer godoc
// @Summary Create a new gacha player
// @Description Create a new gacha player with initial coin and diamond balance
// @Tags GachaMachine
// @Accept json
// @Produce json
// @Param request body dto.CreateGachaPlayerReq true "Gacha player creation request"
// @Success 200 {object} map[string]interface{} "Gacha player created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /gachaMachine/createGachaPlayer [post]
func (h *GachaMachineHandler) HandleCreateGachaPlayer(c *gin.Context) {
	var req dto.CreateGachaPlayerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind JSON", "error", err)
		common.SendError(c, 400, "Invalid request payload")
		return
	}

	grpcReq := &gachaMachine.CreateGachaPlayerReq{
		Player: &gachaMachine.GachaPlayer{
			BasePlayer: &player.Player{
				PlayerID: req.PlayerID,
				UserName: req.UserName,
			},
			Coin:    req.Coin,
			Diamond: req.Diamond,
		},
	}

	resp, err := h.gachaMachineClient.CreateGachaPlayer(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to create gacha player", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully created gacha player", "player_id", req.PlayerID)
	common.SendSuccess(c, resp)
}

// HandleGetGachaPlayerInfo godoc
// @Summary Get gacha player information
// @Description Get gacha player information by player ID
// @Tags GachaMachine
// @Accept json
// @Produce json
// @Param playerID path int64 true "Player ID"
// @Success 200 {object} map[string]interface{} "Gacha player info retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid player ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /gachaMachine/getGachaPlayerInfo/{playerID} [get]
func (h *GachaMachineHandler) HandleGetGachaPlayerInfo(c *gin.Context) {
	playerIDParam := c.Param("playerID")
	var playerID int64
	_, err := fmt.Sscan(playerIDParam, &playerID)
	if err != nil {
		h.logger.Errorw("Invalid player ID", "error", err)
		common.SendError(c, 400, "Invalid player ID")
		return
	}

	grpcReq := &gachaMachine.GetGachaPlayerInfoReq{
		PlayerID: playerID,
	}

	resp, err := h.gachaMachineClient.GetGachaPlayerInfo(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to get gacha player info", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved gacha player info", "player_id", playerID)
	common.SendSuccess(c, resp)
}

// HandleAdjustPlayerCoin godoc
// @Summary Adjust gacha player coin balance
// @Description Add or subtract coins from a gacha player's balance
// @Tags GachaMachine
// @Accept json
// @Produce json
// @Param request body dto.AdjustGachaPlayerCoinRequest true "Adjust coin request"
// @Success 200 {object} map[string]interface{} "Player coin adjusted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /gachaMachine/adjustPlayerCoin [post]
func (h *GachaMachineHandler) HandleAdjustPlayerCoin(c *gin.Context) {
	var req dto.AdjustGachaPlayerCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &gachaMachine.AdjustPlayerCoinReq{
		PlayerID: req.PlayerID,
		Amount:   req.Amount,
		Type:     req.Type,
	}

	resp, err := h.gachaMachineClient.AdjustPlayerCoin(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to adjust player coin", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully adjusted player coin", "player_id", req.PlayerID, "amount", req.Amount, "type", req.Type)
	common.SendSuccess(c, resp)
}

// HandleAdjustPlayerDiamond godoc
// @Summary Adjust gacha player diamond balance
// @Description Add or subtract diamonds from a gacha player's balance
// @Tags GachaMachine
// @Accept json
// @Produce json
// @Param request body dto.AdjustGachaPlayerDiamondRequest true "Adjust diamond request"
// @Success 200 {object} map[string]interface{} "Player diamond adjusted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /gachaMachine/adjustPlayerDiamond [post]
func (h *GachaMachineHandler) HandleAdjustPlayerDiamond(c *gin.Context) {
	var req dto.AdjustGachaPlayerDiamondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &gachaMachine.AdjustPlayerDiamondReq{
		PlayerID: req.PlayerID,
		Amount:   req.Amount,
		Type:     req.Type,
	}

	resp, err := h.gachaMachineClient.AdjustPlayerDiamond(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to adjust player diamond", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully adjusted player diamond", "player_id", req.PlayerID, "amount", req.Amount, "type", req.Type)
	common.SendSuccess(c, resp)
}

// HandleCreateGachaItems godoc
// @Summary Create gacha items
// @Description Create multiple gacha items with their properties
// @Tags GachaMachine
// @Accept json
// @Produce json
// @Param request body dto.CreateGachaItemsRequest true "Gacha items creation request"
// @Success 201 {object} map[string]interface{} "Gacha items created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /gachaMachine/createGachaItems [post]
func (h *GachaMachineHandler) HandleCreateGachaItems(c *gin.Context) {
	var req dto.CreateGachaItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	fmt.Printf("GachaItems: %+v\n", req.GachaItems)

	// Convert DTO to gRPC request
	grpcReq := &gachaMachine.CreateGachaItemsReq{}
	for _, item := range req.GachaItems {
		grpcReq.GachaItems = append(grpcReq.GachaItems, &gachaMachine.CreateGachaItemReq{
			Name:       item.Name,
			Rarity:     item.Rarity,
			PullWeight: item.PullWeight,
		})
	}

	resp, err := h.gachaMachineClient.CreateGachaItems(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to create gacha items", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully created gacha items", "item_count", len(req.GachaItems))
	common.SendCreated(c, resp)
}

// HandleCreateGachaMachine godoc
// @Summary Create a new gacha machine
// @Description Create a new gacha machine with the provided details
// @Tags GachaMachine
// @Accept json
// @Produce json
// @Param request body dto.CreateGachaMachineRequest true "Gacha machine creation request"
// @Success 201 {object} map[string]interface{} "Gacha machine created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /gachaMachine/createGachaMachine [post]
func (h *GachaMachineHandler) HandleCreateGachaMachine(c *gin.Context) {
	var req dto.CreateGachaMachineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &gachaMachine.CreateGachaMachineReq{
		Name:          req.Name,
		Price:         req.Price,
		PriceTimesTen: req.PriceTimesTen,
		SuperRarePity: req.SuperRarePity,
		UltraRarePity: req.UltraRarePity,
	}

	for _, item := range req.GachaItems {
		grpcReq.Items = append(grpcReq.Items, &gachaMachine.Items{
			ItemID: item.ItemID,
		})
	}

	resp, err := h.gachaMachineClient.CreateGachaMachine(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to create gacha machine", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully created gacha machine", "machine_name", req.Name)
	common.SendCreated(c, resp)
}

// HandleGetGachaMachineInfo godoc
// @Summary Get gacha machine information
// @Description Get gacha machine information by machine ID
// @Tags GachaMachine
// @Accept json
// @Produce json
// @Param machineID path int64 true "Machine ID"
// @Success 200 {object} map[string]interface{} "Gacha machine info retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid machine ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /gachaMachine/getGachaMachineInfo/{machineID} [get]
func (h *GachaMachineHandler) HandleGetGachaMachineInfo(c *gin.Context) {
	machineIDParam := c.Param("machineID")
	var machineID int64
	_, err := fmt.Sscan(machineIDParam, &machineID)
	if err != nil {
		h.logger.Errorw("Invalid machine ID", "error", err)
		common.SendError(c, 400, "Invalid machine ID")
		return
	}

	grpcReq := &gachaMachine.GetGachaMachineInfoReq{
		MachineID: machineID,
	}

	resp, err := h.gachaMachineClient.GetGachaMachineInfo(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to get gacha machine info", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved gacha machine info", "machine_id", machineID)
	common.SendSuccess(c, resp)
}

// HandleGetPullResult godoc
// @Summary Get gacha pull result
// @Description Execute a gacha pull and get the result
// @Tags GachaMachine
// @Accept json
// @Produce json
// @Param request body dto.GetPullResultRequest true "Pull result request"
// @Success 200 {object} map[string]interface{} "Pull result retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /gachaMachine/getPullResult [post]
func (h *GachaMachineHandler) HandleGetPullResult(c *gin.Context) {
	var req dto.GetPullResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &gachaMachine.GetPullResultReq{
		MachineID: req.MachineID,
		PlayerID:  req.PlayerID,
		PullCount: req.PullCount,
	}

	resp, err := h.gachaMachineClient.GetPullResult(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to get pull result", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved pull result", "machine_id", req.MachineID, "player_id", req.PlayerID)
	common.SendSuccess(c, resp)
}
