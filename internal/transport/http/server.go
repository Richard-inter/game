package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/transport/grpc"
	"github.com/Richard-inter/game/internal/transport/http/handler"

	_ "github.com/Richard-inter/game/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Game Service API
// @version 1.0
// @description This is the API documentation for the Game Service including Player, ClawMachine, GachaMachine, and WhackAMole endpoints.
// @host localhost:8080
// @BasePath /api/v1

const (
	httpStatusNoContent = 204
)

type Server struct {
	config     *config.ServiceConfig
	logger     *zap.SugaredLogger
	server     *http.Server
	engine     *gin.Engine
	grpcClient *grpc.ClientManager
}

func NewServer(cfg *config.ServiceConfig, logger *zap.SugaredLogger, grpcClient *grpc.ClientManager) *Server {
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

	s.logger.Infow("Starting HTTP server", "address", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server failed to start: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	s.logger.Infow("Shutting down HTTP server")
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
		s.logger.Fatalw("Failed to create player handler", "error", err)
	}

	clawMachineHandler, err := handler.NewClawMachineHandler(s.logger, s.grpcClient)
	if err != nil {
		s.logger.Fatalw("Failed to create claw machine handler", "error", err)
	}

	gachaMachineHandler, err := handler.NewGachaMachineHandler(s.logger, s.grpcClient)
	if err != nil {
		s.logger.Fatalw("Failed to create gacha machine handler", "error", err)
	}

	whackAMoleHandler, err := handler.NewWhackAMoleHandler(s.logger, s.grpcClient)
	if err != nil {
		s.logger.Fatalw("Failed to create whack-a-mole handler", "error", err)
	}

	// Swagger endpoint
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	s.engine.GET("/health", handler.HealthCheck(s.logger.Desugar()))

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
			// items
			clawMachine.POST("/createClawItems", clawMachineHandler.HandleCreateClawItems)

			// machine
			clawMachine.POST("/createClawMachine", clawMachineHandler.HandleCreateClawMachine)
			clawMachine.GET("/getClawMachineInfo/:machineID", clawMachineHandler.HandleGetClawMachineInfo)

			// player
			clawMachine.GET("/getClawPlayerInfo/:playerID", clawMachineHandler.HandleGetClawPlayerInfo)
			clawMachine.POST("/createClawPlayer", clawMachineHandler.HandleCreateClawPlayer)
			clawMachine.POST("/adjustPlayerCoin", clawMachineHandler.HandleAdjustPlayerCoin)
			clawMachine.POST("/adjustPlayerDiamond", clawMachineHandler.HandleAdjustPlayerDiamond)

			// game
			clawMachine.POST("/startClawGame", clawMachineHandler.HandleStartClawGame)
			clawMachine.POST("/addTouchedItemRecord", clawMachineHandler.HandleAddTouchedItemRecord)
		}

		gachaMachine := v1.Group("/gachaMachine")
		{
			// items
			gachaMachine.POST("/createGachaItems", gachaMachineHandler.HandleCreateGachaItems)

			// machine
			gachaMachine.POST("/createGachaMachine", gachaMachineHandler.HandleCreateGachaMachine)
			gachaMachine.GET("/getGachaMachineInfo/:machineID", gachaMachineHandler.HandleGetGachaMachineInfo)

			// player
			gachaMachine.GET("/getGachaPlayerInfo/:playerID", gachaMachineHandler.HandleGetGachaPlayerInfo)
			gachaMachine.POST("/createGachaPlayer", gachaMachineHandler.HandleCreateGachaPlayer)
			gachaMachine.POST("/adjustPlayerCoin", gachaMachineHandler.HandleAdjustPlayerCoin)
			gachaMachine.POST("/adjustPlayerDiamond", gachaMachineHandler.HandleAdjustPlayerDiamond)

			// game
			gachaMachine.POST("/getPullResult", gachaMachineHandler.HandleGetPullResult)
		}

		whackAMole := v1.Group("/whackAMole")
		{
			// player
			whackAMole.POST("/createWhackAMolePlayer", whackAMoleHandler.HandleCreateWhackAMolePlayer)
			whackAMole.GET("/getWhackAMolePlayer/:id", whackAMoleHandler.HandleGetPlayerInfo)

			// leaderboard
			whackAMole.GET("/leaderboard/:limit", whackAMoleHandler.HandleGetLeaderboard)
			whackAMole.POST("/updateScore", whackAMoleHandler.HandleUpdateScore)

			// mole weight configs
			whackAMole.POST("/createMoleWeightConfig", whackAMoleHandler.HandleCreateMoleWeightConfig)
			whackAMole.GET("/getMoleWeightConfig/:id", whackAMoleHandler.HandleGetMoleWeightConfig)
			whackAMole.POST("/updateMoleWeightConfig", whackAMoleHandler.HandleUpdateMoleWeightConfig)
		}
	}
}
