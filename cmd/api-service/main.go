package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/transport/http/handler"
	"github.com/Richard-inter/game/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	shutdownTimeout     = 5 * time.Second
	httpStatusNoContent = 204
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	log := logger.GetLogger()

	log.WithFields(logrus.Fields{
		"version":    Version,
		"build_time": BuildTime,
		"go_version": GoVersion,
		"service":    "api-service",
	}).Info("Starting API Service")

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/api-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	// Set Gin mode
	if cfg.Service.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin engine
	engine := gin.New()

	// Add middleware
	setupMiddleware(engine, cfg, log)

	// Setup routes
	setupRoutes(engine, log)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.GetServiceAddr(),
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.Service.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Service.WriteTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithFields(logrus.Fields{
			"address": server.Addr,
		}).Info("API Service starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start API service")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down API Service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("API service shutdown error")
	}

	log.Info("API Service stopped")
}

func setupMiddleware(engine *gin.Engine, _ *config.ServiceConfig, _ *logrus.Logger) {
	// Recovery middleware
	engine.Use(gin.Recovery())

	// Logger middleware
	engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// CORS middleware
	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(httpStatusNoContent)
			return
		}

		c.Next()
	})
}

func setupRoutes(engine *gin.Engine, log *logrus.Logger) {
	// Health check
	engine.GET("/health", handler.HealthCheck(log))

	// API version 1
	v1 := engine.Group("/api/v1")
	{
		// Game routes
		games := v1.Group("/games")
		{
			games.GET("", handler.ListGames(log))
			games.GET("/:id", handler.GetGame(log))
			games.POST("", handler.CreateGame(log))
			games.PUT("/:id", handler.UpdateGame(log))
			games.DELETE("/:id", handler.DeleteGame(log))
		}

		// Player routes
		players := v1.Group("/players")
		{
			players.GET("", handler.ListPlayers(log))
			players.GET("/:id", handler.GetPlayer(log))
			players.POST("", handler.CreatePlayer(log))
			players.PUT("/:id", handler.UpdatePlayer(log))
			players.DELETE("/:id", handler.DeletePlayer(log))
		}
	}
}
