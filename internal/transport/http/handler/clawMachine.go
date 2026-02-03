package handler

import (
	"fmt"
	"strconv"

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

// HandleCreateClawMachine godoc
// @Summary Create a new claw machine
// @Description Create a new claw machine with the provided details
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.CreateClawMachineRequest true "Claw machine creation request"
// @Success 201 {object} map[string]interface{} "Claw machine created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/createClawMachine [post]
func (h *ClawMachineHandler) HandleCreateClawMachine(c *gin.Context) {
	var req dto.CreateClawMachineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	// Convert DTO to gRPC request
	grpcReq := &clawMachine.CreateClawMachineReq{
		Name:    req.Name,
		Price:   req.Price,
		MaxItem: req.MaxItem,
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
	common.SendCreated(c, resp)
}

// HandleGetClawMachineInfo godoc
// @Summary Get claw machine information
// @Description Get claw machine information by machine ID
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.GetClawMachineInfoRequest true "Get machine info request"
// @Success 200 {object} map[string]interface{} "Claw machine info retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/getClawMachineInfo [post]
func (h *ClawMachineHandler) HandleGetClawMachineInfo(c *gin.Context) {
	var req dto.GetClawMachineInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &clawMachine.GetClawMachineInfoReq{
		MachineID: req.MachineID,
	}

	resp, err := h.clawMachineClient.GetClawMachineInfo(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to get claw machine info", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved claw machine info", "machine_id", req.MachineID)
	common.SendSuccess(c, resp)
}

// HandleCreateClawItems godoc
// @Summary Create claw items
// @Description Create multiple claw items with their properties
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.CreateClawItemsRequest true "Claw items creation request"
// @Success 201 {object} map[string]interface{} "Claw items created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/createClawItems [post]
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
			Name:            item.Name,
			Rarity:          item.Rarity,
			SpawnPercentage: item.SpawnPercentage,
			CatchPercentage: item.CatchPercentage,
			MaxItemSpawned:  item.MaxItemSpawned,
		})
	}

	resp, err := h.clawMachineClient.CreateClawItems(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to create claw items", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully created claw items", "item_count", len(req.ClawItems))
	common.SendCreated(c, resp)
}

// HandleGetClawPlayerInfo godoc
// @Summary Get claw player information
// @Description Get claw player information by player ID
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.GetClawPlayerInfoRequest true "Get player info request"
// @Success 200 {object} map[string]interface{} "Claw player info retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/getClawPlayerInfo [post]
func (h *ClawMachineHandler) HandleGetClawPlayerInfo(c *gin.Context) {
	var req dto.GetClawPlayerInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &clawMachine.GetClawPlayerInfoReq{
		PlayerID: req.PlayerID,
	}

	resp, err := h.clawMachineClient.GetClawPlayerInfo(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to get claw player info", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved claw player info", "player_id", req.PlayerID)
	common.SendSuccess(c, resp)
}

// HandleCreateClawPlayer godoc
// @Summary Create a new claw player
// @Description Create a new claw player with initial coin and diamond balance
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.CreateClawPlayerRequest true "Claw player creation request"
// @Success 201 {object} map[string]interface{} "Claw player created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/createClawPlayer [post]
func (h *ClawMachineHandler) HandleCreateClawPlayer(c *gin.Context) {
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
	common.SendCreated(c, resp)
}

// HandleAdjustPlayerCoin godoc
// @Summary Adjust player coin balance
// @Description Add or subtract coins from a player's balance
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.AdjustPlayerCoinRequest true "Adjust coin request"
// @Success 200 {object} map[string]interface{} "Player coin adjusted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/adjustPlayerCoin [post]
func (h *ClawMachineHandler) HandleAdjustPlayerCoin(c *gin.Context) {
	var req dto.AdjustPlayerCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &clawMachine.AdjustPlayerCoinReq{
		PlayerID: req.PlayerID,
		Amount:   req.Amount,
		Type:     req.Type,
	}

	resp, err := h.clawMachineClient.AdjustPlayerCoin(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to adjust player coin", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully adjusted player coin", "player_id", req.PlayerID, "amount", req.Amount, "type", req.Type)
	common.SendSuccess(c, resp)
}

// HandleAdjustPlayerDiamond godoc
// @Summary Adjust player diamond balance
// @Description Add or subtract diamonds from a player's balance
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.AdjustPlayerDiamondRequest true "Adjust diamond request"
// @Success 200 {object} map[string]interface{} "Player diamond adjusted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/adjustPlayerDiamond [post]
func (h *ClawMachineHandler) HandleAdjustPlayerDiamond(c *gin.Context) {
	var req dto.AdjustPlayerDiamondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &clawMachine.AdjustPlayerDiamondReq{
		PlayerID: req.PlayerID,
		Amount:   req.Amount,
		Type:     req.Type,
	}

	resp, err := h.clawMachineClient.AdjustPlayerDiamond(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to adjust player diamond", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully adjusted player diamond", "player_id", req.PlayerID, "amount", req.Amount, "type", req.Type)
	common.SendSuccess(c, resp)
}

// HandleStartClawGame godoc
// @Summary Start a claw game
// @Description Start a new claw game session for a player
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.StartClawGameRequest true "Start game request"
// @Success 200 {object} map[string]interface{} "Claw game started successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/startClawGame [post]
func (h *ClawMachineHandler) HandleStartClawGame(c *gin.Context) {
	var req dto.StartClawGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &clawMachine.StartClawGameReq{
		PlayerID:  req.PlayerID,
		MachineID: req.MachineID,
	}
	resp, err := h.clawMachineClient.StartClawGame(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to start claw game", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully started claw game", "player_id", req.PlayerID, "machine_id", req.MachineID)
	common.SendSuccess(c, resp)
}

// HandleAddTouchedItemRecord godoc
// @Summary Add touched item record
// @Description Record an item that was touched during a claw game
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.AddTouchedItemRecordRequest true "Add touched item record request"
// @Success 200 {object} map[string]interface{} "Touched item record added successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/addTouchedItemRecord [post]
func (h *ClawMachineHandler) HandleAddTouchedItemRecord(c *gin.Context) {
	var req dto.AddTouchedItemRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &clawMachine.AddTouchedItemRecordReq{
		GameID:  req.GameID,
		ItemID:  req.ItemID,
		Catched: req.Catched,
	}

	resp, err := h.clawMachineClient.AddTouchedItemRecord(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to add touched item record", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully added touched item record", "game_id", req.GameID, "item_id", req.ItemID, "catched", req.Catched)
	common.SendSuccess(c, resp)
}

// HandleDeleteClawPlayer godoc
// @Summary Delete a claw player
// @Description Delete a claw player by player ID
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.DeleteClawPlayerRequest true "Delete player request"
// @Success 200 {object} map[string]interface{} "Claw player deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/deleteClawPlayer [post]
func (h *ClawMachineHandler) HandleDeleteClawPlayer(c *gin.Context) {
	var req dto.DeleteClawPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &clawMachine.DeleteClawPlayerReq{
		PlayerID: req.PlayerID,
	}

	resp, err := h.clawMachineClient.DeleteClawPlayer(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to delete claw player", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully deleted claw player", "player_id", req.PlayerID)
	common.SendSuccess(c, resp)
}

// HandleGetGameHistory godoc
// @Summary Get player game history
// @Description Get game history for a claw player
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.GetGameHistoryRequest true "Get game history request"
// @Success 200 {object} map[string]interface{} "Game history retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/getGameHistory [get]
func (h *ClawMachineHandler) HandleGetGameHistory(c *gin.Context) {
	playerIDStr := c.Query("playerID")
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		h.logger.Errorw("Invalid player ID", "error", err)
		common.SendError(c, 400, "Invalid player ID")
		return
	}

	grpcReq := &clawMachine.GetGameHistoryReq{
		PlayerID: playerID,
	}

	resp, err := h.clawMachineClient.GetGameHistory(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to get game history", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved game history", "player_id", playerID)
	common.SendSuccess(c, resp)
}

// HandleUpdateClawMachineItems godoc
// @Summary Update claw machine items
// @Description Update items in a claw machine
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.UpdateClawMachineItemsRequest true "Update machine items request"
// @Success 200 {object} map[string]interface{} "Claw machine items updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/updateClawMachineItems [post]
func (h *ClawMachineHandler) HandleUpdateClawMachineItems(c *gin.Context) {
	var req dto.UpdateClawMachineItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &clawMachine.UpdateClawMachineItemsReq{
		MachineID: req.MachineID,
	}

	for _, item := range req.Items {
		grpcReq.Items = append(grpcReq.Items, &clawMachine.Items{
			ItemID: item.ItemID,
		})
	}

	resp, err := h.clawMachineClient.UpdateClawMachineItems(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to update claw machine items", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully updated claw machine items", "machine_id", req.MachineID)
	common.SendSuccess(c, resp)
}

// HandleDeleteClawMachine godoc
// @Summary Delete a claw machine
// @Description Delete a claw machine by machine ID
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.DeleteClawMachineRequest true "Delete machine request"
// @Success 200 {object} map[string]interface{} "Claw machine deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/deleteClawMachine [post]
func (h *ClawMachineHandler) HandleDeleteClawMachine(c *gin.Context) {
	var req dto.DeleteClawMachineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &clawMachine.DeleteClawMachineReq{
		MachineID: req.MachineID,
	}

	resp, err := h.clawMachineClient.DeleteClawMachine(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to delete claw machine", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully deleted claw machine", "machine_id", req.MachineID)
	common.SendSuccess(c, resp)
}

// HandleDeleteClawItems godoc
// @Summary Delete claw items
// @Description Delete multiple claw items by their IDs
// @Tags ClawMachine
// @Accept json
// @Produce json
// @Param request body dto.DeleteClawItemsRequest true "Delete claw items request"
// @Success 200 {object} map[string]interface{} "Claw items deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /clawMachine/deleteClawItems [delete]
func (h *ClawMachineHandler) HandleDeleteClawItems(c *gin.Context) {
	var req dto.DeleteClawItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, 400, "Invalid request body")
		return
	}

	grpcReq := &clawMachine.DeleteClawItemsReq{
		ItemIDs: req.ItemIDs,
	}

	resp, err := h.clawMachineClient.DeleteClawItems(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to delete claw items", "error", err)
		common.SendError(c, 500, err.Error())
		return
	}

	h.logger.Infow("Successfully deleted claw items", "item_count", len(req.ItemIDs))
	common.SendSuccess(c, resp)
}
