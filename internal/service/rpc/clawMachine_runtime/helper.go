package clawmachine_runtime

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/Richard-inter/game/internal/domain"
)

// Rarity value multipliers
var BaseRarityMultiplier = map[string]float64{
	"common":    0.05,
	"uncommon":  0.10,
	"rare":      0.25,
	"very rare": 0.50,
	"epic":      1.20,
	"legend":    2.50,
}

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
	SpawnPercent int   // absolute probability (0-100)
	MaxPerRound  int   // soft cap to prevent RNG spikes
	Value        int64 // RTP value based on rarity
}

type RTPState struct {
	TotalSpent   int64
	TotalPayout  int64
	TargetRTP    float64
	TotalWagered int64
	TotalWon     int64
}

// GetRarityValue returns the value for a given rarity
func GetRarityValue(rarity string, price int64) int64 {
	return int64(BaseRarityMultiplier[rarity] * float64(price))
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
func SpawnWithControls(items []SpawnItem, config SpawnConfig, rtpState *domain.ClawMachineRTPState) []SpawnItem {
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
			adjusted := AdjustSpawnWeight(
				item.SpawnPercent,
				item.Value,
				CalculateRTPDelta(rtpState),
			)
			weights = append(weights, adjusted)
			totalWeight += adjusted
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

func (s *ClawMachineWebsocketService) SpawnMachineItems(
	ctx context.Context,
	machineID int64,
) ([]int64, error) {
	clawMachine, err := s.repo.GetClawMachineInfo(ctx, machineID)
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
			Value:        GetRarityValue(item.Item.Rarity, clawMachine.Price),
		})
	}

	// Get RTP state for this machine
	rtpState, err := s.repo.GetClawMachineRTPState(ctx, machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine RTP state: %w", err)
	}

	spawnedItems := SpawnWithControls(spawnItems, config, rtpState)

	// Convert SpawnItem results to item IDs
	spawnedIDs := make([]int64, 0, len(spawnedItems))
	for _, spawnedItem := range spawnedItems {
		spawnedIDs = append(spawnedIDs, spawnedItem.ID)
	}

	return spawnedIDs, nil
}

// PreDetermineCatchResults generates a list of pre-determined catch results for all items
func (s *ClawMachineWebsocketService) PreDetermineCatchResults(
	ctx context.Context,
	machineID int64,
) ([]*CatchResult, error) {
	// Get machine info to access items and their catch percentages
	clawMachine, err := s.repo.GetClawMachineInfo(ctx, machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine info: %w", err)
	}

	if len(clawMachine.Items) == 0 {
		return nil, fmt.Errorf("no items in machine to catch from")
	}

	results := make([]*CatchResult, 0, len(clawMachine.Items))

	rtpState, err := s.repo.GetClawMachineRTPState(ctx, machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine RTP state: %w", err)
	}
	expectedRevenue := rtpState.TotalRevenue + clawMachine.Price
	maxPayout := int64(float64(expectedRevenue) * rtpState.TargetRTP)
	for _, item := range clawMachine.Items {
		catchWeight := item.Item.CatchPercentage
		if catchWeight == 0 {
			return nil, fmt.Errorf("database error: item %s (ID: %d) has zero catch percentage", item.Item.Name, item.Item.ID)
		}

		catchSuccess := Roll(int(catchWeight))

		if catchSuccess && rtpState.TotalPayout+GetRarityValue(item.Item.Rarity, clawMachine.Price) > maxPayout {
			catchSuccess = false
		}

		results = append(results, &CatchResult{
			ItemID:  item.Item.ID,
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
	clawMachine, err := s.repo.GetClawMachineInfo(ctx, machineID)
	if err != nil {
		return fmt.Errorf("failed to get machine info: %w", err)
	}

	_, err = s.repo.AdjustPlayerCoin(ctx, playerID, int64(clawMachine.Price), "minus")
	if err != nil {
		return fmt.Errorf("failed to adjust player coin: %w", err)
	}

	return nil
}

func AdjustSpawnWeight(
	base int,
	itemValue int64,
	rtpDelta float64,
) int {
	if rtpDelta > 0 && itemValue > 0 {
		penalty := int(float64(base) * rtpDelta)
		return max(1, base-penalty)
	}

	if rtpDelta < 0 {
		boost := int(float64(base) * -rtpDelta * 0.5)
		return base + boost
	}

	return base
}

func CalculateRTPDelta(state *domain.ClawMachineRTPState) float64 {
	if state == nil || state.TotalRevenue == 0 {
		return 0
	}
	currentRTP := (float64(state.TotalPayout) / float64(state.TotalRevenue)) * 100.0
	return (currentRTP - state.TargetRTP) / 100.0
}
