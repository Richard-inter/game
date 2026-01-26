package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/db"
	"github.com/Richard-inter/game/internal/repository"
	"github.com/Richard-inter/game/pkg/logger"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	log := logger.GetSugar()

	// Load configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/rpc-whackamole-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.Fatalw("Failed to load configuration", "error", err)
	}

	// Initialize database
	database, err := db.InitWhackAMoleDB(cfg)
	if err != nil {
		log.Fatalw("Failed to initialize database", "error", err)
	}

	// Create WhackAMole repository
	whackAMoleRepo := repository.NewWhackAMoleRepository(database)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create ticker for scheduling (1 second interval)
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	log.Info("Starting leaderboard recalculation scheduler (every 1 second)...")

	// Main scheduling loop
	for {
		select {
		case <-ctx.Done():
			log.Info("Scheduler shutting down...")
			return

		case <-sigChan:
			log.Info("Received shutdown signal, stopping scheduler...")
			cancel()
			return

		case <-ticker.C:
			log.Info("Running scheduled leaderboard recalculation...")

			err := whackAMoleRepo.RecalculateLeaderboard(context.Background())
			if err != nil {
				log.Errorw("Failed to recalculate leaderboard", "error", err)
			} else {
				log.Info("Leaderboard recalculation completed successfully!")
			}
		}
	}
}
