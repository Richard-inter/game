package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/1nterdigital/game/internal/domain"
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

		// 3️⃣ Claw Items (NO manual IDs) - create each individually with FirstOrCreate
		// 3️⃣ Claw Items (fixed dataset, idempotent)
		seedItems := []domain.ClawItem{
			{
				Name:            "Gecko",
				Rarity:          "Common",
				SpawnPercentage: 80,
				CatchPercentage: 90,
				MaxItemSpawned:  3,
			},
			{
				Name:            "Herring",
				Rarity:          "Common",
				SpawnPercentage: 80,
				CatchPercentage: 90,
				MaxItemSpawned:  3,
			},
			{
				Name:            "Monkey",
				Rarity:          "Uncommon",
				SpawnPercentage: 60,
				CatchPercentage: 70,
				MaxItemSpawned:  3,
			},
			{
				Name:            "Muskrat",
				Rarity:          "Rare",
				SpawnPercentage: 40,
				CatchPercentage: 50,
				MaxItemSpawned:  3,
			},
			{
				Name:            "Pudu",
				Rarity:          "VeryRare",
				SpawnPercentage: 30,
				CatchPercentage: 30,
				MaxItemSpawned:  2,
			},
			{
				Name:            "Sparrow",
				Rarity:          "Epic",
				SpawnPercentage: 20,
				CatchPercentage: 10,
				MaxItemSpawned:  2,
			},
			{
				Name:            "Squid",
				Rarity:          "Legend",
				SpawnPercentage: 10,
				CatchPercentage: 5,
				MaxItemSpawned:  1,
			},
		}

		var items []domain.ClawItem

		for _, item := range seedItems {
			item.CreatedAt = now
			item.CreatedBy = createdBy
			item.UpdatedAt = now
			item.UpdatedBy = createdBy
			item.IsActive = true

			// Idempotent insert
			if err := tx.
				Where("name = ?", item.Name).
				FirstOrCreate(&item).Error; err != nil {
				return err
			}

			items = append(items, item)
		}

		// 4️⃣ Machine ↔ Items mapping (use FirstOrCreate to prevent duplicates)
		for i := range items {
			item := &items[i]
		for i := range items {
			item := &items[i]
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

		// 5️⃣ RTP State - use configurable default RTP
		rtp := domain.ClawMachineRTPState{
			ClawMachineID: machine.ID,
			TargetRTP:     85.0, // Default target RTP - will be configurable via service config
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
