package handler

import (
	"net/http"

	"github.com/Richard-inter/game/internal/transport/grpc"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Player struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Score    int    `json:"score"`
}

func ListPlayers(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		players := []Player{
			{ID: "1", Username: "player1", Email: "player1@example.com", Score: 100},
			{ID: "2", Username: "player2", Email: "player2@example.com", Score: 200},
		}

		logger.Info("List players requested")
		c.JSON(http.StatusOK, gin.H{"players": players})
	}
}

func GetPlayer(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		player := Player{ID: id, Username: "player" + id, Email: "player" + id + "@example.com", Score: 150}

		logger.WithField("player_id", id).Info("Get player requested")
		c.JSON(http.StatusOK, player)
	}
}

func CreatePlayer(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var player Player
		if err := c.ShouldBindJSON(&player); err != nil {
			logger.WithError(err).Error("Invalid player data")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		player.ID = "3"
		player.Score = 0

		logger.WithField("player", player).Info("Create player requested")
		c.JSON(http.StatusCreated, player)
	}
}

func UpdatePlayer(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var player Player
		if err := c.ShouldBindJSON(&player); err != nil {
			logger.WithError(err).Error("Invalid player data")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		player.ID = id

		logger.WithField("player_id", id).Info("Update player requested")
		c.JSON(http.StatusOK, player)
	}
}

func DeletePlayer(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		logger.WithField("player_id", id).Info("Delete player requested")
		c.JSON(http.StatusNoContent, gin.H{"message": "Player deleted successfully"})
	}
}

func HandleGetPlayerInfo(logger *logrus.Logger, grpcClient *grpc.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call the gRPC service using the client
		resp, err := grpcClient.GetPlayerInfo(c, 123)
		if err != nil {
			logger.WithError(err).Error("Failed to call gRPC service")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call gRPC service"})
			return
		}

		logger.Info("Successfully called gRPC player service")
		c.JSON(http.StatusOK, gin.H{
			"message":  "Successfully connected to gRPC service",
			"response": resp,
		})
	}
}
