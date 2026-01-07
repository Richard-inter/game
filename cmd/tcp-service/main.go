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
	"github.com/sirupsen/logrus"
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
		"service":    "tcp-service",
	}).Info("Starting TCP Service")

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/tcp-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	// Create listener
	lc := net.ListenConfig{}
	listener, err := lc.Listen(context.Background(), "tcp", cfg.GetServiceAddr())
	if err != nil {
		log.WithError(err).Fatal("Failed to listen on TCP port")
	}
	defer listener.Close()

	log.WithFields(logrus.Fields{
		"address": listener.Addr().String(),
	}).Info("TCP Service starting")

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
					log.WithError(err).Error("TCP accept error")
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

	log.Info("Shutting down TCP Service...")

	// Signal shutdown
	close(quit)
	listener.Close()

	// Give connections time to close gracefully
	time.Sleep(2 * time.Second)

	log.Info("TCP Service stopped")
}

func handleConnection(conn net.Conn, cfg *config.ServiceConfig, log *logrus.Logger) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.WithField("client", clientAddr).Info("TCP client connected")

	// Set read/write timeouts
	if cfg.TCP.ReadTimeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(time.Duration(cfg.TCP.ReadTimeout) * time.Second)); err != nil {
			log.WithError(err).Error("Failed to set read deadline")
		}
	}
	if cfg.TCP.WriteTimeout > 0 {
		if err := conn.SetWriteDeadline(time.Now().Add(time.Duration(cfg.TCP.WriteTimeout) * time.Second)); err != nil {
			log.WithError(err).Error("Failed to set write deadline")
		}
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Bytes()

		log.WithFields(logrus.Fields{
			"client":  clientAddr,
			"message": string(message),
		}).Debug("Received TCP message")

		// Process message (placeholder)
		response := processMessage(message, clientAddr, log)

		// Send response
		if _, err := conn.Write(response); err != nil {
			log.WithError(err).Error("Failed to send TCP response")
			break
		}

		// Reset read deadline
		if cfg.TCP.ReadTimeout > 0 {
			if err := conn.SetReadDeadline(time.Now().Add(time.Duration(cfg.TCP.ReadTimeout) * time.Second)); err != nil {
				log.WithError(err).Error("Failed to reset read deadline")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.WithError(err).WithField("client", clientAddr).Error("TCP client error")
	}

	log.WithField("client", clientAddr).Info("TCP client disconnected")
}

func processMessage(message []byte, clientAddr string, _ *logrus.Logger) []byte {
	// Simple echo server with prefix
	response := fmt.Sprintf("TCP Service Echo [%s]: %s", clientAddr, string(message))
	return []byte(response)
}
