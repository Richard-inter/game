package main

import (
	"fmt"
	"log"

	"github.com/Richard-inter/game/internal/domain"
	clawmachine "github.com/Richard-inter/game/internal/service/rpc/clawMachine"
)

// Mock repository for testing
type mockRepo struct{}

func (m *mockRepo) GetClawMachineInfo(machineID int64) (*domain.ClawMachine, error) {
	return &domain.ClawMachine{
		ID:      machineID,
		Name:    "Test Machine",
		Price:   100,
		MaxItem: 10,
		Items: []domain.ClawMachineItem{
			{
				Item: domain.Item{
					ID:              1,
					Name:            "Common Item",
					Rarity:          "Common",
					SpawnPercentage: 70,
					CatchPercentage: 80,
					MaxItemSpawned:  5,
				},
			},
			{
				Item: domain.Item{
					ID:              2,
					Name:            "Rare Item",
					Rarity:          "Rare",
					SpawnPercentage: 20,
					CatchPercentage: 30,
					MaxItemSpawned:  2,
				},
			},
		},
	}, nil
}

func (m *mockRepo) CreateClawMachine(clawMachine *domain.ClawMachine) (*domain.ClawMachine, error) {
	return clawMachine, nil
}

func (m *mockRepo) CreateClawItems(items *[]domain.Item) (*[]domain.Item, error) {
	return items, nil
}

func (m *mockRepo) CreateClawPlayer(clawPlayer *domain.ClawPlayer) (*domain.ClawPlayer, error) {
	return clawPlayer, nil
}

func (m *mockRepo) GetClawPlayerInfo(playerID int64) (*domain.ClawPlayer, error) {
	return &domain.ClawPlayer{
		Player: domain.Player{
			ID:       playerID,
			UserName: "TestPlayer",
		},
		Coin:    100,
		Diamond: 50,
	}, nil
}

func (m *mockRepo) UpdateClawMachineItems(clawMachineID int64, items []domain.ClawMachineItem) error {
	return nil
}

func main() {
	// Create service with mock repository
	service := clawmachine.NewClawMachineGRPCService(&mockRepo{})

	// Test 1: Spawn items (returns IDs)
	fmt.Println("=== Testing Spawn Function (Returns IDs) ===")
	spawnedIDs, err := service.SpawnMachineItems(nil, 123)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Spawned item IDs: %v\n", spawnedIDs)
	fmt.Printf("Count: %d\n", len(spawnedIDs))

	// Test 2: Determine catch item
	fmt.Println("\n=== Testing Catch Function ===")
	catchResult, err := service.DetermineCatchItem(nil, 123)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Selected catch item: %+v\n", catchResult)
	fmt.Printf("Item ID: %d, Name: %s, Success: %t\n",
		catchResult.ItemID, catchResult.Name, catchResult.Success)
}
