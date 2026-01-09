package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/pkg/logger"
	"go.uber.org/zap"
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

	log.Infow("Starting TCP Service",
		"version", Version,
		"build_time", BuildTime,
		"go_version", GoVersion,
		"service", "tcp-service",
	)

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/tcp-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.Fatalw("Failed to load configuration", "error", err)
	}

	// Create listener
	lc := net.ListenConfig{}
	listener, err := lc.Listen(context.Background(), "tcp", cfg.GetServiceAddr())
	if err != nil {
		log.Fatalw("Failed to listen on TCP port", "error", err)
	}
	defer listener.Close()

	log.Infow("TCP Service starting", "address", listener.Addr().String())

	// Channel to signal shutdown
	quit := make(chan struct{})

	// Start connection handler
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-quit:
					return // Shutdown in progress
				default:
					log.Errorw("TCP accept error", "error", err)
					continue
				}
			}

			go handleConnection(conn, cfg, log)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Infow("Shutting down TCP Service...")

	// Signal shutdown
	close(quit)
	listener.Close()

	// Give connections time to close gracefully
	time.Sleep(2 * time.Second)

	log.Infow("TCP Service stopped")
}

func handleConnection(conn net.Conn, cfg *config.ServiceConfig, log *zap.SugaredLogger) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Infow("TCP client connected", "client", clientAddr)

	// Set read/write timeouts
	if cfg.TCP.ReadTimeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(time.Duration(cfg.TCP.ReadTimeout) * time.Second)); err != nil {
			log.Errorw("Failed to set read deadline", "error", err)
		}
	}
	if cfg.TCP.WriteTimeout > 0 {
		if err := conn.SetWriteDeadline(time.Now().Add(time.Duration(cfg.TCP.WriteTimeout) * time.Second)); err != nil {
			log.Errorw("Failed to set write deadline", "error", err)
		}
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Bytes()

		log.Debugw("Received TCP message", "client", clientAddr, "message", string(message))

		// Process message (placeholder)
		response := processMessage(message, clientAddr, log)

		// Send response
		if _, err := conn.Write(response); err != nil {
			log.Errorw("Failed to send TCP response", "error", err)
			break
		}

		// Reset read deadline
		if cfg.TCP.ReadTimeout > 0 {
			if err := conn.SetReadDeadline(time.Now().Add(time.Duration(cfg.TCP.ReadTimeout) * time.Second)); err != nil {
				log.Errorw("Failed to reset read deadline", "error", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Errorw("TCP client error", "error", err, "client", clientAddr)
	}

	log.Infow("TCP client disconnected", "client", clientAddr)
}

func processMessage(message []byte, clientAddr string, _ *zap.SugaredLogger) []byte {
	// Simple echo server with prefix
	response := fmt.Sprintf("TCP Service Echo [%s]: %s", clientAddr, string(message))
	return []byte(response)
}
