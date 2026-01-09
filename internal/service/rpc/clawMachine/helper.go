package clawmachine

import (
	"context"
	"fmt"
	"math/rand/v2"

	pb "github.com/Richard-inter/game/pkg/protocol/clawMachine"
)

type SpawnConfig struct {
	MaxOutput int // max items in machine output
}

type CatchResult struct {
	ItemID  int64  `json:"itemID"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
}

type SpawnItem struct {
	ID           int64
	SpawnPercent int // absolute probability (0-100)
	MaxPerRound  int // soft cap to prevent RNG spikes
}

func Roll(percent int) bool {
	if percent <= 0 {
		return false
	}
	if percent >= 100 {
		return true
	}
	return rand.IntN(100) < percent
}

// AdjustForPity slightly increases chances for very rare items
// after multiple failed passes
// NOTE: Commented out for future use - currently not using pity system
/*
func AdjustForPity(item SpawnItem, pass int) int {
	// Example rule:
	// Items <= 5% get a boost after pass 1
	if item.SpawnPercent <= 5 && pass >= 1 {
		boost := pass * 5 // +5% per extra pass
		adjusted := item.SpawnPercent + boost
		if adjusted > 100 {
			return 100
		}
		return adjusted
	}
	return item.SpawnPercent
}
*/

func SpawnWithControls(items []SpawnItem, config SpawnConfig) []SpawnItem {
	result := make([]SpawnItem, 0, config.MaxOutput)
	counts := make(map[int64]int)

	for len(result) < config.MaxOutput {
		availableItems := make([]SpawnItem, 0)
		weights := make([]int, 0)
		totalWeight := 0

		for _, item := range items {
			if counts[item.ID] >= item.MaxPerRound {
				continue
			}

			availableItems = append(availableItems, item)
			weights = append(weights, item.SpawnPercent)
			totalWeight += item.SpawnPercent
		}

		if len(availableItems) == 0 {
			break
		}

		selection := rand.IntN(totalWeight)
		currentWeight := 0

		for idx, weight := range weights {
			currentWeight += weight
			if selection < currentWeight {
				selectedItem := availableItems[idx]
				result = append(result, selectedItem)
				counts[selectedItem.ID]++
				break
			}
		}
	}

	return result
}

func (s *ClawMachineGRPCServices) GetMachineItems(
	ctx context.Context,
	machineID int64,
) ([]*pb.Item, error) {
	clawMachine, err := s.repo.GetClawMachineInfo(machineID)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.Item, 0, len(clawMachine.Items))
	for _, item := range clawMachine.Items {
		items = append(items, &pb.Item{
			ItemID:          item.Item.ID,
			Name:            item.Item.Name,
			Rarity:          item.Item.Rarity,
			SpawnPercentage: item.Item.SpawnPercentage,
			CatchPercentage: item.Item.CatchPercentage,
			MaxItemSpawned:  item.Item.MaxItemSpawned,
		})
	}

	return items, nil
}

func (s *ClawMachineGRPCServices) SpawnMachineItems(
	ctx context.Context,
	machineID int64,
) ([]int64, error) {
	clawMachine, err := s.repo.GetClawMachineInfo(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine info: %w", err)
	}
	config := SpawnConfig{
		MaxOutput: int(clawMachine.MaxItem), // Use machine's MaxItem as output cap
	}

	spawnItems := make([]SpawnItem, 0, len(clawMachine.Items))
	for _, item := range clawMachine.Items {
		spawnItems = append(spawnItems, SpawnItem{
			ID:           item.Item.ID,
			SpawnPercent: int(item.Item.SpawnPercentage),
			MaxPerRound:  int(item.Item.MaxItemSpawned),
		})
	}

	spawnedItems := SpawnWithControls(spawnItems, config)

	// Convert SpawnItem results to item IDs
	spawnedIDs := make([]int64, 0, len(spawnedItems))
	for _, spawnedItem := range spawnedItems {
		spawnedIDs = append(spawnedIDs, spawnedItem.ID)
	}

	return spawnedIDs, nil
}

// TODO
func (s *ClawMachineGRPCServices) DetermineCatchItem(
	ctx context.Context,
	machineID int64,
) (*CatchResult, error) {
	// Get machine info to access items and their catch percentages
	clawMachine, err := s.repo.GetClawMachineInfo(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine info: %w", err)
	}

	if len(clawMachine.Items) == 0 {
		return nil, fmt.Errorf("no items in machine to catch from")
	}

	// Create weighted selection based on actual catch percentages
	weights := make([]int64, 0, len(clawMachine.Items))
	totalWeight := int64(0)

	for _, item := range clawMachine.Items {
		// Get catch percentage from machine items
		catchWeight := item.Item.CatchPercentage
		if catchWeight == 0 {
			catchWeight = 50 // Default 50% if not set
		}

		weights = append(weights, catchWeight)
		totalWeight += catchWeight
	}

	// Weighted random selection
	selection := rand.IntN(int(totalWeight))
	currentWeight := int64(0)

	for idx, weight := range weights {
		currentWeight += weight
		if int64(selection) < currentWeight {
			return &CatchResult{
				ItemID:  clawMachine.Items[idx].ID,
				Name:    clawMachine.Items[idx].Item.Name,
				Success: true, // If selected for catch, it's successful
			}, nil
		}
	}

	// Fallback to last item
	lastIndex := len(clawMachine.Items) - 1
	return &CatchResult{
		ItemID:  clawMachine.Items[lastIndex].ID,
		Name:    clawMachine.Items[lastIndex].Item.Name,
		Success: true,
	}, nil
}
