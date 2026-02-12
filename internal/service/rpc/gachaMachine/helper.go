package gachaMachine

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/1nterdigital/game/internal/domain"
)

var (
	globalRand *rand.Rand
	randInit   sync.Once
)

func getGlobalRand() *rand.Rand {
	randInit.Do(func() {
		globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
	return globalRand
}

type Entry struct {
	ID     int64
	Weight int32
}

// RTP helper functions
func GetGachaRarityValue(rarity string, price int64) int64 {
	// Gacha rarity values
	rarityMultipliers := map[string]float64{
		"normal":     0.10,
		"rare":       0.40,
		"super_rare": 1.20,
		"ultra_rare": 2.50,
	}

	if multiplier, exists := rarityMultipliers[rarity]; exists {
		return int64(multiplier * float64(price))
	}

	// default to normal
	return int64(0.10 * float64(price))
}

func CalculateGachaRTPDelta(state *domain.GachaMachineRTPState) float64 {
	if state == nil || state.TotalRevenue == 0 {
		return 0
	}
	currentRTP := (float64(state.TotalPayout) / float64(state.TotalRevenue)) * 100.0
	// Return the percentage difference (no extra division by 100)
	return currentRTP - state.TargetRTP
}

func AdjustGachaPullWeight(base int32, itemValue int64, rtpDelta float64) int32 {
	if rtpDelta > 0 && itemValue > 0 {
		// When RTP is too high, reduce pull probability
		// Use a more moderate penalty factor (0.5 instead of 1.0)
		penalty := int32(float64(base) * rtpDelta * 0.5)
		if base-penalty < 1 {
			return 1
		}
		return base - penalty
	}

	if rtpDelta < 0 {
		// When RTP is too low, increase pull probability
		// Use a moderate boost factor (0.3 instead of 0.5)
		boost := int32(float64(base) * -rtpDelta * 0.3)
		return base + boost
	}

	return base
}

func PullGachaByEntries(entries []Entry) int64 {
	if len(entries) == 0 {
		return 0
	}

	var total int32
	safeEntries := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if e.Weight > 0 {
			safeEntries = append(safeEntries, e)
			total += e.Weight
		}
	}

	if total == 0 {
		return 0
	}

	r := getGlobalRand().Int31n(total)
	var sum int32

	for _, e := range safeEntries {
		sum += e.Weight
		if r < sum {
			return e.ID
		}
	}

	return 0
}

func (s *GachaMachineGRPCService) PullGachaByMachineID(
	ctx context.Context,
	pityState *domain.GachaPityState,
	resp *domain.GachaMachine,
) int64 {
	// Get RTP state for this machine
	rtpState, err := s.repo.GetGachaMachineRTPState(ctx, resp.ID)
	if err != nil {
		// If RTP state doesn't exist, use nil (no adjustment)
		rtpState = nil
	}

	if pityState.UltraRarePityCount >= resp.UltraRarePity {
		return s.pullByRarity(resp, "ultra_rare", rtpState)
	}

	if pityState.SuperRarePityCount >= resp.SuperRarePity {
		return s.pullByRarity(resp, "super_rare", rtpState)
	}

	return s.pullFromAll(resp, rtpState)
}

func (s *GachaMachineGRPCService) pullByRarity(
	resp *domain.GachaMachine,
	rarity string,
	rtpState *domain.GachaMachineRTPState,
) int64 {
	entries := make([]Entry, 0)
	rtpDelta := CalculateGachaRTPDelta(rtpState)

	for _, item := range resp.Items {
		if item.Item.Rarity == rarity {
			itemValue := GetGachaRarityValue(item.Item.Rarity, resp.Price)
			adjustedWeight := AdjustGachaPullWeight(item.Item.PullWeight, itemValue, rtpDelta)
			entries = append(entries, Entry{
				ID:     item.Item.ID,
				Weight: adjustedWeight,
			})
		}
	}
	return PullGachaByEntries(entries)
}

func (s *GachaMachineGRPCService) pullFromAll(
	resp *domain.GachaMachine,
	rtpState *domain.GachaMachineRTPState,
) int64 {
	entries := make([]Entry, 0, len(resp.Items))
	rtpDelta := CalculateGachaRTPDelta(rtpState)

	for _, item := range resp.Items {
		itemValue := GetGachaRarityValue(item.Item.Rarity, resp.Price)
		adjustedWeight := AdjustGachaPullWeight(item.Item.PullWeight, itemValue, rtpDelta)
		entries = append(entries, Entry{
			ID:     item.Item.ID,
			Weight: adjustedWeight,
		})
	}
	return PullGachaByEntries(entries)
}

func (s *GachaMachineGRPCService) PullGachaSingle(
	ctx context.Context,
	machineID, playerID int64,
) (int64, error) {
	resp, err := s.repo.GetGachaMachineInfo(ctx, machineID)
	if err != nil || resp == nil {
		return 0, err
	}

	pityState, err := s.repo.GetGachaPityState(ctx, machineID, playerID)
	if err != nil || pityState == nil {
		return 0, err
	}

	itemID := s.PullGachaByMachineID(ctx, pityState, resp)

	s.updatePityAfterPull(pityState, itemID, resp)

	if err := s.repo.SetGachaPityState(ctx, pityState); err != nil {
		return 0, err
	}

	err = s.repo.AddItemInventory(ctx, playerID, itemID, 1)
	if err != nil {
		return 0, err
	}

	return itemID, nil
}

func (s *GachaMachineGRPCService) PullGachaByMachineIDMulti(
	ctx context.Context,
	machineID, playerID int64,
	count int,
) ([]int64, error) {
	resultMap := make(map[int64]int32)
	returnResults := make([]int64, 0, count)

	resp, err := s.repo.GetGachaMachineInfo(ctx, machineID)
	if err != nil {
		return nil, err
	}

	pityState, err := s.repo.GetGachaPityState(ctx, machineID, playerID)
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		itemID := s.PullGachaByMachineID(ctx, pityState, resp)
		returnResults = append(returnResults, itemID)
		resultMap[itemID]++

		s.updatePityAfterPull(pityState, itemID, resp)

		_ = s.redis.SetGachaPityStateToRedis(ctx, machineID, playerID, pityState)
	}

	_ = s.repo.SetGachaPityState(ctx, pityState)
	_ = s.redis.DeleteGachaPityStateFromRedis(ctx, machineID, playerID)

	for itemID, qty := range resultMap {
		err := s.repo.AddItemInventory(ctx, playerID, itemID, qty)
		if err != nil {
			return nil, err
		}
	}

	return returnResults, nil
}

func (s *GachaMachineGRPCService) updatePityAfterPull(
	pity *domain.GachaPityState,
	itemID int64,
	resp *domain.GachaMachine,
) {
	rarity := getItemRarity(itemID, resp)

	switch rarity {
	case "ultra_rare":
		pity.UltraRarePityCount = 0
		pity.SuperRarePityCount++

	case "super_rare":
		pity.SuperRarePityCount = 0
		pity.UltraRarePityCount++

	default:
		pity.SuperRarePityCount++
		pity.UltraRarePityCount++
	}
}

func getItemRarity(itemID int64, resp *domain.GachaMachine) string {
	for _, item := range resp.Items {
		if item.Item.ID == itemID {
			return item.Item.Rarity
		}
	}
	return ""
}

func (s *GachaMachineGRPCService) PlayMachine(ctx context.Context, playerID, machineID int64, pullCount int32) error {
	resp, err := s.repo.GetGachaMachineInfo(ctx, machineID)
	if err != nil {
		return err
	}

	if pullCount == 1 {
		_, err = s.repo.AdjustPlayerCoin(ctx, playerID, int64(resp.Price), "minus")
		if err != nil {
			return err
		}

		// Update RTP state to track revenue for this pull (regardless of outcome)
		err = s.repo.UpdateGachaMachineRTP(ctx, machineID, resp.Price, 0)
		if err != nil {
			// Log error but don't fail the play since revenue tracking is secondary
			fmt.Printf("failed to update gacha machine RTP state: %v", err)
		}
	}

	if pullCount == 10 {
		_, err = s.repo.AdjustPlayerCoin(ctx, playerID, int64(resp.PriceTimesTen), "minus")
		if err != nil {
			return err
		}

		// Update RTP state for 10-pull revenue
		err = s.repo.UpdateGachaMachineRTP(ctx, machineID, resp.PriceTimesTen, 0)
		if err != nil {
			fmt.Printf("failed to update gacha machine RTP state: %v", err)
		}
	}

	return nil
}
