package clawmachine

import (
	"context"
	"fmt"
	"math/rand/v2"
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

func (s *ClawMachineWebsocketService) PreDetermineCatchResults(
	ctx context.Context,
	machineID int64,
) ([]*CatchResult, error) {
	// Get machine info to access items and their catch percentages
	clawMachine, err := s.repo.GetClawMachineInfo(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine info: %w", err)
	}

	if len(clawMachine.Items) == 0 {
		return nil, fmt.Errorf("no items in machine to catch from")
	}

	results := make([]*CatchResult, 0, len(clawMachine.Items))

	// Generate pre-determined result for each item in the machine
	for _, item := range clawMachine.Items {
		catchWeight := item.Item.CatchPercentage
		if catchWeight == 0 {
			return nil, fmt.Errorf("database error: item %s (ID: %d) has zero catch percentage", item.Item.Name, item.Item.ID)
		}

		// Determine if catch is successful based on the item's catch percentage
		catchSuccess := Roll(int(catchWeight))

		results = append(results, &CatchResult{
			ItemID:  item.ID,
			Name:    item.Item.Name,
			Success: catchSuccess,
		})
	}

	return results, nil
}

func (s *ClawMachineWebsocketService) PlayMachine(
	ctx context.Context,
	playerID int64,
	machineID int64,
) error {
	clawMachine, err := s.repo.GetClawMachineInfo(machineID)
	if err != nil {
		return fmt.Errorf("failed to get machine info: %w", err)
	}

	_, err = s.repo.AdjustPlayerCoin(playerID, int64(clawMachine.Price), "minus")
	if err != nil {
		return fmt.Errorf("failed to adjust player coin: %w", err)
	}

	return nil
}
