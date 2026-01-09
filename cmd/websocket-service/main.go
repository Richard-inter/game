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
	"github.com/Richard-inter/game/pkg/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 5 * time.Second
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	log := logger.GetSugar()

	log.Infow("Starting WebSocket Service",
		"version", Version,
		"build_time", BuildTime,
		"go_version", GoVersion,
		"service", "websocket-service",
	)

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/websocket-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.Fatalw("Failed to load configuration", "error", err)
	}

	// WebSocket upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  cfg.WebSocket.ReadBufferSize,
		WriteBufferSize: cfg.WebSocket.WriteBufferSize,
		CheckOrigin: func(_ *http.Request) bool {
			return true // Allow all origins for now
		},
	}

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc(cfg.Service.Path, func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(upgrader, w, r, log)
	})

	// Add health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"websocket-service"}`)
	})

	server := &http.Server{
		Addr:         cfg.GetServiceAddr(),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Infow("WebSocket Service starting", "address", server.Addr, "path", cfg.WebSocket.Path)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("Failed to start WebSocket service", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infow("Shutting down WebSocket Service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorw("WebSocket service shutdown error", "error", err)
	}

	log.Infow("WebSocket Service stopped")
}

func handleWebSocket(upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request, log *zap.SugaredLogger) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorw("Failed to upgrade WebSocket connection", "error", err)
		return
	}
	defer conn.Close()

	log.Infow("WebSocket client connected", "client_ip", r.RemoteAddr, "user_agent", r.UserAgent())

	// Handle WebSocket messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorw("WebSocket error", "error", err)
			}
			break
		}

		log.Debugw("Received WebSocket message", "message_type", messageType, "message", string(message))

		// Echo message back (placeholder)
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Errorw("Failed to send WebSocket message", "error", err)
			break
		}
	}

	log.Infow("WebSocket client disconnected", "client_ip", r.RemoteAddr)
}
