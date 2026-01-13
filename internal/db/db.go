package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/domain"
)

// InitDB initializes a database connection
func InitDB(cfg *config.ServiceConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&domain.Player{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// InitPlayerDB initializes player database connection
func InitPlayerDB(cfg *config.ServiceConfig) (*gorm.DB, error) {
	dsn := cfg.GetPlayerDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to player database: %w", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&domain.Player{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate player database: %w", err)
	}

	return db, nil
}

// InitClawmachineDB initializes clawmachine database connection
func InitClawmachineDB(cfg *config.ServiceConfig) (*gorm.DB, error) {
	dsn := cfg.GetClawmachineDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to clawmachine database: %w", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&domain.ClawMachine{}, &domain.ClawMachineItem{}, &domain.Item{}, &domain.ClawPlayer{}, &domain.ClawMachineGameRecord{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate clawmachine database: %w", err)
	}

	return db, nil
}
