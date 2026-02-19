package clawmachine

import (
	"context"
	crypto_rand "crypto/rand"
	"fmt"
	"math/big"

	"github.com/1nterdigital/game/internal/domain"
	"github.com/1nterdigital/game/pkg/logger"
	pb "github.com/1nterdigital/game/pkg/protocol/clawMachine"
)

const (
	// RTP adjustment factors
	RTPPenaltyFactor = 0.5 // Reduce spawn weight by up to 50% when RTP is too high
	RTPBoostFactor   = 0.3 // Increase spawn weight by up to 30% when RTP is too low

	// Probability constants
	MaxProbability = 100

	// RTP adjustment constants
	CatchPenaltyFactor = 0.02 // 2% adjustment per percentage point
	CatchBoostFactor   = 0.01 // 1% boost per percentage point
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
	if percent >= MaxProbability {
		return true
	}
	n, _ := crypto_rand.Int(crypto_rand.Reader, big.NewInt(int64(MaxProbability)))
	return int(n.Int64()) < percent
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

		randVal, _ := crypto_rand.Int(crypto_rand.Reader, big.NewInt(int64(totalWeight)))
		selection := int(randVal.Int64())
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
	clawMachine, err := s.repo.GetClawMachineInfo(ctx, machineID)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.Item, 0, len(clawMachine.Items))
	for i := range clawMachine.Items {
		item := &clawMachine.Items[i]
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
	clawMachine, err := s.repo.GetClawMachineInfo(ctx, machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine info: %w", err)
	}
	config := SpawnConfig{
		MaxOutput: int(clawMachine.MaxItem), // Use machine's MaxItem as output cap
	}

	spawnItems := make([]SpawnItem, 0, len(clawMachine.Items))
	for i := range clawMachine.Items {
		item := &clawMachine.Items[i]
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
func (s *ClawMachineGRPCServices) PreDetermineCatchResults(
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

	rtpDelta := CalculateRTPDelta(rtpState)

	for i := range clawMachine.Items {
		item := &clawMachine.Items[i]
		catchWeight := item.Item.CatchPercentage
		if catchWeight == 0 {
			return nil, fmt.Errorf("database error: item %s (ID: %d) has zero catch percentage", item.Item.Name, item.Item.ID)
		}

		// Apply RTP-based adjustment to catch probability
		adjustedCatchWeight := adjustCatchProbability(int(catchWeight), rtpDelta)
		catchSuccess := Roll(adjustedCatchWeight)

		// Additional safety check to prevent extreme overpayout
		if catchSuccess {
			expectedRevenue := rtpState.TotalRevenue + clawMachine.Price
			maxPayout := int64(float64(expectedRevenue) * (rtpState.TargetRTP + 10.0)) // Allow 10% buffer
			itemValue := GetRarityValue(item.Item.Rarity, clawMachine.Price)
			if rtpState.TotalPayout+itemValue > maxPayout {
				catchSuccess = false
			}
		}

		results = append(results, &CatchResult{
			ItemID:  item.Item.ID,
			Name:    item.Item.Name,
			Success: catchSuccess,
		})
	}

	return results, nil
}

func (s *ClawMachineGRPCServices) PlayMachine(
	ctx context.Context,
	playerID int64,
	machineID int64,
) error {
	clawMachine, err := s.repo.GetClawMachineInfo(ctx, machineID)
	if err != nil {
		return fmt.Errorf("failed to get machine info: %w", err)
	}

	_, err = s.repo.AdjustPlayerCoin(ctx, playerID, clawMachine.Price, "minus")
	if err != nil {
		return fmt.Errorf("failed to adjust player coin: %w", err)
	}

	// Update RTP state to track revenue for this play (regardless of outcome)
	err = s.repo.UpdateClawMachineRTP(ctx, machineID, clawMachine.Price, 0)
	if err != nil {
		// Log error but don't fail the play since revenue tracking is secondary
		logger.GetSugar().Warnf("failed to update RTP state for revenue tracking: %v", err)
	}

	return nil
}

func adjustCatchProbability(baseProbability int, rtpDelta float64) int {
	if rtpDelta > 0 {
		// When RTP is too high, reduce catch probability
		// Use a conservative adjustment factor
		penalty := int(float64(baseProbability) * rtpDelta * CatchPenaltyFactor) // 2% adjustment per percentage point
		return max(1, baseProbability-penalty)
	}

	if rtpDelta < 0 {
		// When RTP is too low, increase catch probability
		// Use a conservative boost factor
		boost := int(float64(baseProbability) * -rtpDelta * CatchBoostFactor) // 1% boost per percentage point
		return min(MaxProbability, baseProbability+boost)
	}

	return baseProbability
}

func AdjustSpawnWeight(
	base int,
	itemValue int64,
	rtpDelta float64,
) int {
	if rtpDelta > 0 && itemValue > 0 {
		// When RTP is too high, reduce spawn probability
		// Use a more moderate penalty factor (0.5 instead of 1.0)
		penalty := int(float64(base) * rtpDelta * RTPPenaltyFactor)
		return max(1, base-penalty)
	}

	if rtpDelta < 0 {
		// When RTP is too low, increase spawn probability
		// Use a moderate boost factor (0.3 instead of 0.5)
		boost := int(float64(base) * -rtpDelta * RTPBoostFactor)
		return base + boost
	}

	return base
}

func CalculateRTPDelta(state *domain.ClawMachineRTPState) float64 {
	if state == nil || state.TotalRevenue == 0 {
		return 0
	}
	currentRTP := (float64(state.TotalPayout) / float64(state.TotalRevenue)) * 100.0
	// Return the percentage difference (no extra division by 100)
	return currentRTP - state.TargetRTP
}
