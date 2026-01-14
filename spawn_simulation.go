package main

import (
	"fmt"
	"math/rand/v2"
)

// SpawnConfig represents spawn algorithm configuration
type SpawnConfig struct {
	MaxOutput int // max items in machine output
}

// SpawnItem represents a machine item for spawn algorithm
type SpawnItem struct {
	ID           int64
	Name         string
	SpawnPercent int // absolute probability (0-100)
	MaxPerRound  int // soft cap to prevent RNG spikes
}

// Roll returns true if a percentage-based roll succeeds
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

// SpawnWithControls generates exactly MaxOutput items
// using production-style control layers
func SpawnWithControls(items []SpawnItem, config SpawnConfig) []SpawnItem {
	result := make([]SpawnItem, 0, config.MaxOutput)
	counts := make(map[int64]int)

	// Keep spawning until we have exactly MaxOutput items
	for len(result) < config.MaxOutput {
		// Create a weighted list of available items (respecting max per round)
		availableItems := make([]SpawnItem, 0)
		weights := make([]int, 0)
		totalWeight := 0

		for _, item := range items {
			// Skip if this item has reached its max per round
			if counts[item.ID] >= item.MaxPerRound {
				continue
			}

			// Add to available items with its spawn percentage as weight
			availableItems = append(availableItems, item)
			weights = append(weights, item.SpawnPercent)
			totalWeight += item.SpawnPercent
		}

		// If no items available (all reached max), break to avoid infinite loop
		if len(availableItems) == 0 {
			break
		}

		// Weighted random selection from available items
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

func main() {
	// Your test data
	items := []SpawnItem{
		{ID: 1, Name: "Gold Coin", SpawnPercent: 22, MaxPerRound: 6},
		{ID: 2, Name: "Plastic Ring", SpawnPercent: 18, MaxPerRound: 5},
		{ID: 3, Name: "Keychain Toy", SpawnPercent: 15, MaxPerRound: 4},
		{ID: 4, Name: "Mini Plush", SpawnPercent: 12, MaxPerRound: 3},
		{ID: 5, Name: "Sticker Pack", SpawnPercent: 10, MaxPerRound: 3},
		{ID: 6, Name: "LED Toy", SpawnPercent: 8, MaxPerRound: 2},
		{ID: 7, Name: "Metal Figurine", SpawnPercent: 6, MaxPerRound: 2},
		{ID: 8, Name: "Gaming Mouse Toy", SpawnPercent: 5, MaxPerRound: 1},
		{ID: 9, Name: "Smartwatch Replica", SpawnPercent: 3, MaxPerRound: 1},
		{ID: 10, Name: "Diamond Toy", SpawnPercent: 1, MaxPerRound: 1},
	}

	config := SpawnConfig{
		MaxOutput: 10, // maxItem from your data
	}

	fmt.Println("=== Spawn Simulation Results (10 runs) ===")
	fmt.Println()

	for run := 1; run <= 10; run++ {
		fmt.Printf("Run %d:\n", run)

		spawnedItems := SpawnWithControls(items, config)

		fmt.Printf("  Spawned %d items:\n", len(spawnedItems))
		for _, item := range spawnedItems {
			fmt.Printf("    - ID:%d %s\n", item.ID, item.Name)
		}

		// Show counts by item type
		counts := make(map[int64]int)
		for _, item := range spawnedItems {
			counts[item.ID]++
		}

		fmt.Printf("  Summary:\n")
		for _, item := range items {
			if count, exists := counts[item.ID]; exists {
				fmt.Printf("    %s: %d\n", item.Name, count)
			}
		}
		fmt.Println()
	}
}
