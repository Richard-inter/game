package gachaMachine_runtime

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/Richard-inter/game/internal/domain"
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

func (s *GachaMachineWebsocketService) PullGachaByMachineID(
	ctx context.Context,
	pityState *domain.GachaPityState,
	resp *domain.GachaMachine,
) int64 {
	if pityState.UltraRarePityCount >= resp.UltraRarePity {
		return s.pullByRarity(resp, "ultra_rare")
	}

	if pityState.SuperRarePityCount >= resp.SuperRarePity {
		return s.pullByRarity(resp, "super_rare")
	}

	return s.pullFromAll(resp)
}

func (s *GachaMachineWebsocketService) pullByRarity(
	resp *domain.GachaMachine,
	rarity string,
) int64 {
	entries := make([]Entry, 0)
	for _, item := range resp.Items {
		if item.Item.Rarity == rarity {
			entries = append(entries, Entry{
				ID:     item.Item.ID,
				Weight: item.Item.PullWeight,
			})
		}
	}
	return PullGachaByEntries(entries)
}

func (s *GachaMachineWebsocketService) pullFromAll(
	resp *domain.GachaMachine,
) int64 {
	entries := make([]Entry, 0, len(resp.Items))
	for _, item := range resp.Items {
		entries = append(entries, Entry{
			ID:     item.Item.ID,
			Weight: item.Item.PullWeight,
		})
	}
	return PullGachaByEntries(entries)
}

func (s *GachaMachineWebsocketService) PullGachaSingle(
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

	return itemID, nil
}

func (s *GachaMachineWebsocketService) PullGachaByMachineIDMulti(
	ctx context.Context,
	machineID, playerID int64,
	count int,
) ([]int64, error) {
	results := make([]int64, 0, count)

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
		results = append(results, itemID)

		s.updatePityAfterPull(pityState, itemID, resp)

		_ = s.redis.SetGachaPityStateToRedis(ctx, machineID, playerID, pityState)
	}

	_ = s.repo.SetGachaPityState(ctx, pityState)
	_ = s.redis.DeleteGachaPityStateFromRedis(ctx, machineID, playerID)

	return results, nil
}

func (s *GachaMachineWebsocketService) updatePityAfterPull(
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

func (s *GachaMachineWebsocketService) PlayMachine(ctx context.Context, playerID, machineID int64, pullCount int32) error {
	resp, err := s.repo.GetGachaMachineInfo(ctx, machineID)
	if err != nil {
		return err
	}

	if pullCount == 1 {
		_, err = s.repo.AdjustPlayerCoin(ctx, playerID, int64(resp.Price), "minus")
		if err != nil {
			return err
		}
	}

	if pullCount == 10 {
		_, err = s.repo.AdjustPlayerCoin(ctx, playerID, int64(resp.PriceTimesTen), "minus")
		if err != nil {
			return err
		}
	}

	return nil
}
