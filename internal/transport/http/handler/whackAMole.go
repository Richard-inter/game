package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/transport/grpc"
	dto "github.com/Richard-inter/game/internal/transport/http/DTO"
	"github.com/Richard-inter/game/pkg/common"
	whackAMolepb "github.com/Richard-inter/game/pkg/protocol/whackAMole"
)

type WhackAMoleHandler struct {
	logger           *zap.SugaredLogger
	whackAMoleClient *grpc.WhackAMoleClient
}

func NewWhackAMoleHandler(
	logger *zap.SugaredLogger,
	grpcManager *grpc.ClientManager,
) (*WhackAMoleHandler, error) {
	whackAMoleClient, err := grpcManager.GetWhackAMoleClient()
	if err != nil {
		return nil, err
	}

	return &WhackAMoleHandler{
		logger:           logger,
		whackAMoleClient: whackAMoleClient,
	}, nil
}

// HandleCreateWhackAMolePlayer godoc
// @Summary Create a new Whack-A-Mole player
// @Description Create a new Whack-A-Mole player with username
// @Tags WhackAMole
// @Accept json
// @Produce json
// @Param request body dto.CreateWhackAMolePlayerRequest true "Whack-A-Mole player creation request"
// @Success 201 {object} map[string]interface{} "Whack-A-Mole player created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /whackAMole/createWhackAMolePlayer [post]
func (h *WhackAMoleHandler) HandleCreateWhackAMolePlayer(c *gin.Context) {
	var req dto.CreateWhackAMolePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	grpcReq := &whackAMolepb.CreateWhackAMolePlayerReq{
		PlayerId: req.PlayerID,
		Username: req.Username,
	}

	resp, err := h.whackAMoleClient.CreateWhackAMolePlayer(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to create Whack-A-Mole player", "error", err)
		common.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Infow("Successfully created Whack-A-Mole player", "playerID", req.PlayerID)
	common.SendCreated(c, resp)
}

// HandleGetPlayerInfo godoc
// @Summary Get Whack-A-Mole player information
// @Description Get Whack-A-Mole player information by player ID
// @Tags WhackAMole
// @Accept json
// @Produce json
// @Param id path int64 true "Player ID"
// @Success 200 {object} map[string]interface{} "Player info retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid player ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /whackAMole/getWhackAMolePlayer/{id} [get]
func (h *WhackAMoleHandler) HandleGetPlayerInfo(c *gin.Context) {
	playerIDStr := c.Param("id")
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		h.logger.Errorw("Invalid player ID", "error", err)
		common.SendError(c, http.StatusBadRequest, "Invalid player ID")
		return
	}

	resp, err := h.whackAMoleClient.GetPlayerInfo(c, &whackAMolepb.GetPlayerInfoReq{
		PlayerId: playerID,
	})
	if err != nil {
		h.logger.Errorw("Failed to get player info", "error", err)
		common.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved player info", "player_id", playerID)
	common.SendSuccessWithMessage(c, resp, "Player info retrieved successfully")
}

// HandleGetLeaderboard godoc
// @Summary Get Whack-A-Mole leaderboard
// @Description Get the top players leaderboard
// @Tags WhackAMole
// @Accept json
// @Produce json
// @Param limit path int true "Limit of entries to return"
// @Success 200 {object} map[string]interface{} "Leaderboard retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid limit parameter"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /whackAMole/leaderboard/{limit} [get]
func (h *WhackAMoleHandler) HandleGetLeaderboard(c *gin.Context) {
	limitStr := c.Param("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit <= 0 {
		h.logger.Errorw("Invalid limit parameter", "error", err)
		common.SendError(c, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	grpcReq := &whackAMolepb.GetLeaderboardReq{
		Limit: int32(limit),
	}

	resp, err := h.whackAMoleClient.GetLeaderboard(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to get leaderboard", "error", err)
		common.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved leaderboard", "entries", resp)
	common.SendSuccessWithMessage(c, resp, "Leaderboard retrieved successfully")
}

// HandleGetMoleWeightConfig godoc
// @Summary Get mole weight configuration
// @Description Get mole weight configuration by config ID
// @Tags WhackAMole
// @Accept json
// @Produce json
// @Param id path int64 true "Config ID"
// @Success 200 {object} map[string]interface{} "Mole weight config retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid config ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /whackAMole/getMoleWeightConfig/{id} [get]
func (h *WhackAMoleHandler) HandleGetMoleWeightConfig(c *gin.Context) {
	configIDStr := c.Param("id")
	configID, err := strconv.ParseInt(configIDStr, 10, 64)
	if err != nil {
		h.logger.Errorw("Invalid config ID", "error", err)
		common.SendError(c, http.StatusBadRequest, "Invalid config ID")
		return
	}

	resp, err := h.whackAMoleClient.GetMoleWeightConfig(c, &whackAMolepb.GetMoleWeightConfigReq{
		Id: configID,
	})
	if err != nil {
		h.logger.Errorw("Failed to get mole weight config", "error", err)
		common.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved mole weight config", "configID", configID)
	common.SendSuccessWithMessage(c, resp, "Mole weight config retrieved successfully")
}

// HandleCreateMoleWeightConfig godoc
// @Summary Create mole weight configuration
// @Description Create a new mole weight configuration
// @Tags WhackAMole
// @Accept json
// @Produce json
// @Param request body dto.CreateMoleWeightConfigRequest true "Mole weight config creation request"
// @Success 201 {object} map[string]interface{} "Mole weight config created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /whackAMole/createMoleWeightConfig [post]
func (h *WhackAMoleHandler) HandleCreateMoleWeightConfig(c *gin.Context) {
	var req dto.CreateMoleWeightConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	grpcReq := &whackAMolepb.CreateMoleWeightConfigReq{
		MoleType: req.MoleType,
		Weight:   req.Weight,
	}

	resp, err := h.whackAMoleClient.CreateMoleWeightConfig(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to create mole weight config", "error", err)
		common.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Infow("Successfully created mole weight config", "moleType", req.MoleType, "weight", req.Weight)
	common.SendCreated(c, resp)
}

// HandleUpdateMoleWeightConfig godoc
// @Summary Update mole weight configuration
// @Description Update an existing mole weight configuration
// @Tags WhackAMole
// @Accept json
// @Produce json
// @Param request body dto.UpdateMoleWeightConfigRequest true "Mole weight config update request"
// @Success 200 {object} map[string]interface{} "Mole weight config updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /whackAMole/updateMoleWeightConfig [post]
func (h *WhackAMoleHandler) HandleUpdateMoleWeightConfig(c *gin.Context) {
	var req dto.UpdateMoleWeightConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	grpcReq := &whackAMolepb.UpdateMoleWeightConfigReq{
		Id:       req.ID,
		MoleType: req.MoleType,
		Weight:   req.Weight,
	}

	resp, err := h.whackAMoleClient.UpdateMoleWeightConfig(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to update mole weight config", "error", err)
		common.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Infow("Successfully updated mole weight config", "configID", req.ID)
	common.SendSuccessWithMessage(c, resp, "Mole weight config updated successfully")
}

// HandleUpdateScore godoc
// @Summary Update player score
// @Description Update a player's score in Whack-A-Mole
// @Tags WhackAMole
// @Accept json
// @Produce json
// @Param request body dto.UpdateScoreRequest true "Update score request"
// @Success 200 {object} map[string]interface{} "Player score updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /whackAMole/updateScore [post]
func (h *WhackAMoleHandler) HandleUpdateScore(c *gin.Context) {
	var req dto.UpdateScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	grpcReq := &whackAMolepb.UpdateScoreReq{
		PlayerId: req.PlayerID,
		Score:    req.Score,
	}

	resp, err := h.whackAMoleClient.UpdateScore(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to update player score", "error", err)
		common.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Infow("Successfully updated player score", "playerId", req.PlayerID, "score", req.Score)
	common.SendSuccessWithMessage(c, resp, "Player score updated successfully")
}
