package gachaMachine

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

//
// ---------- TYPES ----------
//

// GachaItem is a concrete reward (character, weapon, etc.)
type GachaItem struct {
	ID     int
	Rarity string
	Weight int
}

// Result sent back to client
type RollResult struct {
	ItemID int    `json:"item_id"`
	Rarity string `json:"rarity"`
}

// Player pity state (stored server-side, DB/Redis)
type PlayerGachaState struct {
	PullsSinceSSR int
}

// Banner configuration (usually loaded from JSON/DB)
type Banner struct {
	ID          int
	RarityRates map[string]int         // rarity -> weight
	Pools       map[string][]GachaItem // rarity -> items
}

// RollHelper is the authoritative gacha machine
type RollHelper struct {
	rng   *rand.Rand
	mutex sync.Mutex
}

//
// ---------- CONSTRUCTOR ----------
//

// NewRollHelper creates a server-owned RNG instance
func NewRollHelper() *RollHelper {
	return &RollHelper{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

//
// ---------- PUBLIC API ----------
//

// Roll performs a single authoritative gacha pull
func (h *RollHelper) Roll(
	banner Banner,
	state *PlayerGachaState,
) (RollResult, error) {

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if state == nil {
		return RollResult{}, errors.New("nil gacha state")
	}

	// Step 1: roll rarity
	rarity := h.rollRarity(banner.RarityRates)

	// Step 2: apply pity
	rarity = applyPity(state, rarity)

	// Step 3: roll item
	item, err := h.rollItem(banner.Pools[rarity])
	if err != nil {
		return RollResult{}, err
	}

	// Step 4: update pity state
	updatePity(state, rarity)

	// Step 5: return client-safe result
	return RollResult{
		ItemID: item.ID,
		Rarity: item.Rarity,
	}, nil
}

//
// ---------- INTERNAL LOGIC ----------
//

func (h *RollHelper) rollRarity(rates map[string]int) string {
	total := 0
	for _, w := range rates {
		total += w
	}

	roll := h.rng.Intn(total)
	acc := 0

	for rarity, w := range rates {
		acc += w
		if roll < acc {
			return rarity
		}
	}

	panic("unreachable")
}

func (h *RollHelper) rollItem(items []GachaItem) (GachaItem, error) {
	if len(items) == 0 {
		return GachaItem{}, errors.New("empty item pool")
	}

	total := 0
	for _, i := range items {
		total += i.Weight
	}

	roll := h.rng.Intn(total)
	acc := 0

	for _, i := range items {
		acc += i.Weight
		if roll < acc {
			return i, nil
		}
	}

	panic("unreachable")
}

func applyPity(state *PlayerGachaState, rarity string) string {
	// Example: guaranteed SSR at 90 pulls
	if state.PullsSinceSSR >= 89 {
		return "SSR"
	}
	return rarity
}

func updatePity(state *PlayerGachaState, rarity string) {
	if rarity == "SSR" {
		state.PullsSinceSSR = 0
	} else {
		state.PullsSinceSSR++
	}
}
