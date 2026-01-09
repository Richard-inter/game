package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

func HealthCheck(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Service:   "game-server",
			Version:   "1.0.0",
		}

		logger.Info("Health check requested",
			zap.String("status", response.Status),
			zap.String("service", response.Service),
			zap.String("client_ip", c.ClientIP()),
		)

		c.JSON(http.StatusOK, response)
	}
}
