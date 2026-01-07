package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

func HealthCheck(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Service:   "game-server",
			Version:   "1.0.0",
		}

		logger.WithFields(logrus.Fields{
			"status":    response.Status,
			"service":   response.Service,
			"client_ip": c.ClientIP(),
		}).Info("Health check requested")

		c.JSON(http.StatusOK, response)
	}
}
