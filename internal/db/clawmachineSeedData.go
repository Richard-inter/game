package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Richard-inter/game/internal/domain"
)

// SeedClawMachineDataIfEmpty
// Idempotent: safe to run multiple times
func SeedClawMachineDataIfEmpty(db *gorm.DB) error {
	fmt.Println("Running claw machine seed (idempotent)")
	return SeedClawMachineData(db)
}

func SeedClawMachineData(db *gorm.DB) error {
	now := time.Now()
	createdBy := "seed"

	return db.Transaction(func(tx *gorm.DB) error {

		// 1️⃣ Player (NO manual ID)
		player := domain.ClawPlayer{
			Player: domain.Player{
				UserName: "clawPlayer",
			},
			Coin:      1000,
			Diamond:   50,
			IsActive:  true,
			CreatedAt: now,
			CreatedBy: createdBy,
			UpdatedAt: now,
			UpdatedBy: createdBy,
		}

		if err := tx.
			Where("player_user_name = ?", "clawPlayer").
			FirstOrCreate(&player).Error; err != nil {
			return err
		}

		// 2️⃣ Claw Machine
		machine := domain.ClawMachine{
			Name:      "Beginner Claw Machine",
			Price:     10,
			MaxItem:   10,
			IsActive:  true,
			CreatedAt: now,
			CreatedBy: createdBy,
			UpdatedAt: now,
			UpdatedBy: createdBy,
		}

		if err := tx.
			Where("name = ?", machine.Name).
			FirstOrCreate(&machine).Error; err != nil {
			return err
		}

		// 3️⃣ Claw Items (NO manual IDs)
		rarityDistribution := []struct {
			Rarity string
			Count  int
			Spawn  float64
			Catch  float64
		}{
			{"common", 4, 30, 60},
			{"uncommon", 2, 20, 45},
			{"rare", 2, 15, 30},
			{"very rare", 1, 8, 18},
			{"epic", 1, 5, 10},
		}

		var items []domain.ClawItem

		for _, r := range rarityDistribution {
			for i := 0; i < r.Count; i++ {
				items = append(items, domain.ClawItem{
					Name:            fmt.Sprintf("%s Item %d", r.Rarity, i+1),
					Rarity:          r.Rarity,
					SpawnPercentage: r.Spawn,
					CatchPercentage: r.Catch,
					MaxItemSpawned:  100,
					IsActive:        true,
					CreatedAt:       now,
					CreatedBy:       createdBy,
					UpdatedAt:       now,
					UpdatedBy:       createdBy,
				})
			}
		}

		// Insert items safely (no duplicates on restart)
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&items).Error; err != nil {
			return err
		}

		// Reload items to get their IDs
		var storedItems []domain.ClawItem
		if err := tx.Find(&storedItems).Error; err != nil {
			return err
		}

		// 4️⃣ Machine ↔ Items mapping
		for _, item := range storedItems {
			machineItem := domain.ClawMachineItem{
				ClawMachineID: machine.ID,
				ItemID:        item.ID,
				IsActive:      true,
				CreatedAt:     now,
				CreatedBy:     createdBy,
				UpdatedAt:     now,
				UpdatedBy:     createdBy,
			}

			if err := tx.
				Where("claw_machine_id = ? AND item_id = ?", machine.ID, item.ID).
				FirstOrCreate(&machineItem).Error; err != nil {
				return err
			}
		}

		// 5️⃣ RTP State
		rtp := domain.ClawMachineRTPState{
			ClawMachineID: machine.ID,
			TargetRTP:     60.0,
			IsActive:      true,
			CreatedAt:     now,
			CreatedBy:     createdBy,
			UpdatedAt:     now,
			UpdatedBy:     createdBy,
		}

		if err := tx.
			Where("claw_machine_id = ?", machine.ID).
			FirstOrCreate(&rtp).Error; err != nil {
			return err
		}

		return nil
	})
}
