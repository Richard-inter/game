package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/transport/grpc"
	"github.com/Richard-inter/game/internal/transport/http/handler"
)

const (
	httpStatusNoContent = 204
)

type Server struct {
	config     *config.ServiceConfig
	logger     *logrus.Logger
	server     *http.Server
	engine     *gin.Engine
	grpcClient *grpc.ClientManager
}

func NewServer(cfg *config.ServiceConfig, logger *logrus.Logger, grpcClient *grpc.ClientManager) *Server {
	return &Server{
		config:     cfg,
		logger:     logger,
		grpcClient: grpcClient,
	}
}

func (s *Server) Start() error {
	// Initialize Gin engine
	s.engine = gin.New()

	// Add middleware
	s.setupMiddleware()

	// Setup routes
	s.setupRoutes()

	// Create HTTP server
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Service.Host, s.config.Service.Port),
		Handler:      s.engine,
		ReadTimeout:  time.Duration(s.config.Service.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Service.WriteTimeout) * time.Second,
	}

	s.logger.WithFields(logrus.Fields{
		"address": s.server.Addr,
	}).Info("Starting HTTP server")

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server failed to start: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}

func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.engine.Use(gin.Recovery())

	// Logger middleware
	s.engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
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
	s.engine.Use(func(c *gin.Context) {
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

func (s *Server) setupRoutes() {
	// Create player handler
	playerHandler, err := handler.NewPlayerHandler(s.logger, s.grpcClient)
	if err != nil {
		s.logger.WithError(err).Fatal("Failed to create player handler")
	}

	clawMachineHandler, err := handler.NewClawMachineHandler(s.logger, s.grpcClient)
	if err != nil {
		s.logger.WithError(err).Fatal("Failed to create claw machine handler")
	}

	// Health check
	s.engine.GET("/health", handler.HealthCheck(s.logger))

	// API version 1
	v1 := s.engine.Group("/api/v1")
	{
		player := v1.Group("/player")
		{
			player.POST("/create", playerHandler.HandleCreatePlayer)
			player.GET("/info/:id", playerHandler.HandleGetPlayerInfo)
		}

		clawMachine := v1.Group("/clawMachine")
		{
			clawMachine.POST("/createClawMachine", clawMachineHandler.HandleCreateClawMachine)
			clawMachine.POST("/createClawItems", clawMachineHandler.HandleCreateClawItems)
		}
	}
}
