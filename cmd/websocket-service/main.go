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
	"github.com/sirupsen/logrus"
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
	log := logger.GetLogger()

	log.WithFields(logrus.Fields{
		"version":    Version,
		"build_time": BuildTime,
		"go_version": GoVersion,
		"service":    "websocket-service",
	}).Info("Starting WebSocket Service")

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/websocket-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
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
		log.WithFields(logrus.Fields{
			"address": server.Addr,
			"path":    cfg.WebSocket.Path,
		}).Info("WebSocket Service starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start WebSocket service")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down WebSocket Service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("WebSocket service shutdown error")
	}

	log.Info("WebSocket Service stopped")
}

func handleWebSocket(upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("Failed to upgrade WebSocket connection")
		return
	}
	defer conn.Close()

	log.WithFields(logrus.Fields{
		"client_ip":  r.RemoteAddr,
		"user_agent": r.UserAgent(),
	}).Info("WebSocket client connected")

	// Handle WebSocket messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithError(err).Error("WebSocket error")
			}
			break
		}

		log.WithFields(logrus.Fields{
			"message_type": messageType,
			"message":      string(message),
		}).Debug("Received WebSocket message")

		// Echo message back (placeholder)
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.WithError(err).Error("Failed to send WebSocket message")
			break
		}
	}

	log.WithField("client_ip", r.RemoteAddr).Info("WebSocket client disconnected")
}
