package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/transport/grpc"
	dto "github.com/Richard-inter/game/internal/transport/http/DTO"
	"github.com/Richard-inter/game/pkg/common"
	player "github.com/Richard-inter/game/pkg/protocol/player"
)

type PlayerHandler struct {
	logger       *zap.SugaredLogger
	playerClient *grpc.PlayerClient
}

func NewPlayerHandler(
	logger *zap.SugaredLogger,
	grpcManager *grpc.ClientManager,
) (*PlayerHandler, error) {
	playerClient, err := grpcManager.GetPlayerClient()
	if err != nil {
		return nil, err
	}

	return &PlayerHandler{
		logger:       logger,
		playerClient: playerClient,
	}, nil
}

// HandleCreatePlayer godoc
// @Summary Create a new player
// @Description Create a new player with the provided username
// @Tags Player
// @Accept json
// @Produce json
// @Param request body dto.CreatePlayerRequest true "Player creation request"
// @Success 201 {object} map[string]interface{} "Player created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /player/create [post]
func (h *PlayerHandler) HandleCreatePlayer(c *gin.Context) {
	var req dto.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request body", "error", err)
		common.SendError(c, http.StatusBadRequest, "Username is required")
		return
	}

	// Convert DTO to gRPC request
	grpcReq := &player.CreatePlayerReq{
		UserName: req.UserName,
	}

	resp, err := h.playerClient.CreatePlayer(c, grpcReq)
	if err != nil {
		h.logger.Errorw("Failed to create player", "error", err)
		common.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Infow("Successfully created player", "username", req.UserName)
	common.SendCreated(c, resp.Player)
}

// HandleGetPlayerInfo godoc
// @Summary Get player information
// @Description Get player information by player ID
// @Tags Player
// @Accept json
// @Produce json
// @Param id path int64 true "Player ID"
// @Success 200 {object} map[string]interface{} "Player info retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid player ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /player/info/{id} [get]
func (h *PlayerHandler) HandleGetPlayerInfo(c *gin.Context) {
	playerIDStr := c.Param("id")
	var playerID int64
	if _, err := fmt.Sscanf(playerIDStr, "%d", &playerID); err != nil {
		h.logger.Errorw("Invalid player ID", "error", err)
		common.SendError(c, http.StatusBadRequest, "Invalid player ID")
		return
	}

	resp, err := h.playerClient.GetPlayerInfo(c, &player.GetPlayerInfoReq{
		PlayerID: playerID,
	})
	if err != nil {
		h.logger.Errorw("Failed to get player info", "error", err)
		common.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Infow("Successfully retrieved player info", "player_id", playerID)
	common.SendSuccessWithMessage(c, resp.Player, "Player info retrieved successfully")
}
