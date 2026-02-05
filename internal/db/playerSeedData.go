package db

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Richard-inter/game/internal/domain"
)

// SeedPlayerData seeds initial players with specified usernames
// Idempotent: safe to run multiple times
func SeedPlayerData(db *gorm.DB) error {
	// Check if any players already exist
	var count int64
	if err := db.Model(&domain.Player{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check if players exist: %w", err)
	}

	// Only seed if no players exist
	if count > 0 {
		fmt.Println("Players already exist, skipping player seeding")
		return nil
	}

	fmt.Println("Seeding initial players...")

	// Players to create
	players := []domain.Player{
		{UserName: "clawPlayer"},
		{UserName: "gachaPlayer"},
		{UserName: "molePlayer"},
	}

	// Create players
	for _, player := range players {
		if err := db.Create(&player).Error; err != nil {
			return fmt.Errorf("failed to create player %s: %w", player.UserName, err)
		}

		fmt.Printf("Created player: %s (ID: %d)\n", player.UserName, player.ID)
	}

	fmt.Println("Player seeding completed successfully")
	return nil
}
