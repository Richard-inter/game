package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Richard-inter/game/internal/transport/grpc"
	"github.com/Richard-inter/game/pkg/common"
	player "github.com/Richard-inter/game/pkg/protocol/player"
)

type PlayerHandler struct {
	logger       *logrus.Logger
	playerClient *grpc.PlayerClient
}

func NewPlayerHandler(
	logger *logrus.Logger,
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

func (h *PlayerHandler) HandleCreatePlayer(c *gin.Context) {
	var grpcReq player.CreatePlayerReq
	if err := c.ShouldBindJSON(&grpcReq); err != nil {
		h.logger.WithError(err).Error("Invalid request body")
		common.SendError(c, http.StatusBadRequest, "Username is required")
		return
	}

	resp, err := h.playerClient.CreatePlayer(c, &grpcReq)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create player")
		common.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.WithField("username", grpcReq.UserName).Info("Successfully created player")
	common.SendCreated(c, resp.Player)
}

func (h *PlayerHandler) HandleGetPlayerInfo(c *gin.Context) {
	playerIDStr := c.Param("id")
	var playerID int64
	if _, err := fmt.Sscanf(playerIDStr, "%d", &playerID); err != nil {
		h.logger.WithError(err).Error("Invalid player ID")
		common.SendError(c, http.StatusBadRequest, "Invalid player ID")
		return
	}

	resp, err := h.playerClient.GetPlayerInfo(c, &player.GetPlayerInfoReq{
		PlayerID: playerID,
	})
	if err != nil {
		h.logger.WithError(err).Error("Failed to get player info")
		common.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.WithField("player_id", playerID).Info("Successfully retrieved player info")
	common.SendSuccessWithMessage(c, resp.Player, "Player info retrieved successfully")
}
