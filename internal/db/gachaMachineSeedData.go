package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/Richard-inter/game/internal/domain"
)

// SeedGachaMachineDataIfEmpty
// Idempotent: safe to run multiple times
func SeedGachaMachineDataIfEmpty(db *gorm.DB) error {
	fmt.Println("Running gacha machine seed (idempotent)")
	return SeedGachaMachineData(db)
}

func SeedGachaMachineData(db *gorm.DB) error {
	now := time.Now()
	createdBy := "seed"

	return db.Transaction(func(tx *gorm.DB) error {

		// 1️⃣ Player (NO manual ID)
		player := domain.GachaPlayer{
			Player: domain.Player{
				UserName: "gachaPlayer",
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
			Where("player_user_name = ?", "gachaPlayer").
			FirstOrCreate(&player).Error; err != nil {
			return err
		}

		// 2️⃣ Gacha Machine
		machine := domain.GachaMachine{
			Name:          "Beginner Gacha Machine",
			Price:         10,
			PriceTimesTen: 90,
			SuperRarePity: 50,
			UltraRarePity: 100,
			IsActive:      true,
			CreatedAt:     now,
			CreatedBy:     createdBy,
			UpdatedAt:     now,
			UpdatedBy:     createdBy,
		}

		if err := tx.
			Where("name = ?", machine.Name).
			FirstOrCreate(&machine).Error; err != nil {
			return err
		}

		// 3️⃣ Gacha Items (NO manual IDs) - create each individually with FirstOrCreate
		rarityDistribution := []struct {
			Rarity string
			Count  int
			Weight int32
		}{
			{"common", 4, 30},
			{"uncommon", 2, 20},
			{"rare", 2, 15},
			{"very_rare", 1, 8},
			{"super_rare", 1, 5},
			{"ultra_rare", 1, 2},
		}

		var items []domain.GachaItem

		for _, r := range rarityDistribution {
			for i := 0; i < r.Count; i++ {
				item := domain.GachaItem{
					Name:       fmt.Sprintf("%s Item %d", r.Rarity, i+1),
					Rarity:     r.Rarity,
					PullWeight: r.Weight,
					IsActive:   true,
					CreatedAt:  now,
					CreatedBy:  createdBy,
					UpdatedAt:  now,
					UpdatedBy:  createdBy,
				}

				// Create item with FirstOrCreate to prevent duplicates
				if err := tx.Where("name = ?", item.Name).FirstOrCreate(&item).Error; err != nil {
					return err
				}

				items = append(items, item)
			}
		}

		// 4️⃣ Machine ↔ Items mapping (use FirstOrCreate to prevent duplicates)
		for _, item := range items {
			machineItem := domain.GachaMachineItem{
				GachaMachineID: machine.ID,
				ItemID:         item.ID,
				IsActive:       true,
				CreatedAt:      now,
				CreatedBy:      createdBy,
				UpdatedAt:      now,
				UpdatedBy:      createdBy,
			}

			if err := tx.
				Where("gacha_machine_id = ? AND item_id = ?", machine.ID, item.ID).
				FirstOrCreate(&machineItem).Error; err != nil {
				return err
			}
		}

		// 5️⃣ RTP State - use configurable default RTP
		rtp := domain.GachaMachineRTPState{
			GachaMachineID: machine.ID,
			TargetRTP:      85.0, // Default target RTP - will be configurable via service config
			IsActive:       true,
			CreatedAt:      now,
			CreatedBy:      createdBy,
			UpdatedAt:      now,
			UpdatedBy:      createdBy,
		}

		if err := tx.
			Where("gacha_machine_id = ?", machine.ID).
			FirstOrCreate(&rtp).Error; err != nil {
			return err
		}

		return nil
	})
}
