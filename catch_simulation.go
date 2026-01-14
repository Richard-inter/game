package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Richard-inter/game/internal/domain"
	clawmachine "github.com/Richard-inter/game/internal/service/rpc/clawMachine"
)

// Mock repository for testing
type MockRepo struct{}

func (m *MockRepo) GetClawMachineInfo(machineID int64) (*domain.ClawMachine, error) {
	// Mock data for testing
	items := []domain.ClawMachineItem{
		{
			ID:            1,
			ClawMachineID: machineID,
			ItemID:        101,
			Item: domain.Item{
				ID:              101,
				Name:            "Common Toy",
				Rarity:          "common",
				SpawnPercentage: 30,
				CatchPercentage: 70,
				MaxItemSpawned:  5,
			},
		},
		{
			ID:            2,
			ClawMachineID: machineID,
			ItemID:        102,
			Item: domain.Item{
				ID:              102,
				Name:            "Rare Toy",
				Rarity:          "rare",
				SpawnPercentage: 15,
				CatchPercentage: 40,
				MaxItemSpawned:  2,
			},
		},
		{
			ID:            3,
			ClawMachineID: machineID,
			ItemID:        103,
			Item: domain.Item{
				ID:              103,
				Name:            "Epic Toy",
				Rarity:          "epic",
				SpawnPercentage: 5,
				CatchPercentage: 20,
				MaxItemSpawned:  1,
			},
		},
	}

	return &domain.ClawMachine{
		ID:      machineID,
		Name:    "Test Machine",
		Price:   100,
		MaxItem: 10,
		Items:   items,
	}, nil
}

// Implement other required methods for the interface
func (m *MockRepo) CreateClawPlayer(clawPlayer *domain.ClawPlayer) (*domain.ClawPlayer, error) {
	return clawPlayer, nil
}

func (m *MockRepo) GetClawPlayerInfo(playerID int64) (*domain.ClawPlayer, error) {
	return &domain.ClawPlayer{}, nil
}

func (m *MockRepo) CreateClawMachine(clawMachine *domain.ClawMachine) (*domain.ClawMachine, error) {
	return clawMachine, nil
}

func (m *MockRepo) UpdateClawMachineItems(clawMachineID int64, items []domain.ClawMachineItem) error {
	return nil
}

func (m *MockRepo) CreateClawItems(items *[]domain.Item) (*[]domain.Item, error) {
	return items, nil
}

func runCatchSimulation() {
	// Create service instance with mock repository
	service := clawmachine.NewClawMachineGRPCService(&MockRepo{})

	ctx := context.Background()
	machineID := int64(1)

	fmt.Println("=== Claw Machine Catch Simulation ===")
	fmt.Printf("Machine ID: %d\n\n", machineID)

	// Test generating pre-determined results for all items
	fmt.Println("--- Pre-determined results for all items in machine ---")
	results, err := service.PreDetermineCatchResults(ctx, machineID)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Generated %d pre-determined results:\n", len(results))
	for i, catch := range results {
		status := "❌ Missed"
		if catch.Success {
			status = "✅ Caught"
		}
		fmt.Printf("Result %d: Item: %s (ID: %d) | Result: %s\n", i+1, catch.Name, catch.ItemID, status)
	}
	fmt.Println()

	// Run multiple simulations to show different random results
	for i := 0; i < 5; i++ {
		fmt.Printf("--- Simulation %d ---\n", i+1)

		result, err := service.PreDetermineCatchResults(ctx, machineID)
		if err != nil {
			log.Printf("Error in simulation %d: %v\n", i+1, err)
			continue
		}

		for _, catch := range result {
			status := "❌ Missed"
			if catch.Success {
				status = "✅ Caught"
			}
			fmt.Printf("Item: %s (ID: %d) | Result: %s\n", catch.Name, catch.ItemID, status)
		}
		fmt.Println()
	}
}

func main() {
	runCatchSimulation()
}
