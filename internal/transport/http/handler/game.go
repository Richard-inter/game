package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Game struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func ListGames(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		games := []Game{
			{ID: "1", Name: "Game 1", Status: "active"},
			{ID: "2", Name: "Game 2", Status: "inactive"},
		}

		logger.Info("List games requested")
		c.JSON(http.StatusOK, gin.H{"games": games})
	}
}

func GetGame(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		game := Game{ID: id, Name: "Game " + id, Status: "active"}

		logger.WithField("game_id", id).Info("Get game requested")
		c.JSON(http.StatusOK, game)
	}
}

func CreateGame(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var game Game
		if err := c.ShouldBindJSON(&game); err != nil {
			logger.WithError(err).Error("Invalid game data")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		game.ID = "3"
		game.Status = "active"

		logger.WithField("game", game).Info("Create game requested")
		c.JSON(http.StatusCreated, game)
	}
}

func UpdateGame(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var game Game
		if err := c.ShouldBindJSON(&game); err != nil {
			logger.WithError(err).Error("Invalid game data")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		game.ID = id

		logger.WithField("game_id", id).Info("Update game requested")
		c.JSON(http.StatusOK, game)
	}
}

func DeleteGame(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		logger.WithField("game_id", id).Info("Delete game requested")
		c.JSON(http.StatusOK, gin.H{"message": "Game deleted successfully"})
	}
}
