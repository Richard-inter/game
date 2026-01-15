package domain

import (
	"time"

	"gorm.io/gorm"
)

// GachaPlayer represents a gacha player with gacha-specific attributes
type GachaPlayer struct {
	gorm.Model
	PlayerID   int64  `gorm:"uniqueIndex;not null"` // Reference to base player
	Gems       int64  `gorm:"default:0"`            // Premium currency
	Tickets    int64  `gorm:"default:0"`            // Free pull currency
	OwnedItems string `gorm:"type:text"`            // JSON array of owned item IDs
	LastPullAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// GachaPool represents a gacha pool with items and settings
type GachaPool struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Cost      int64  `gorm:"not null"`  // Cost per pull
	Currency  string `gorm:"not null"`  // "gems" or "tickets"
	MaxPulls  int32  `gorm:"default:0"` // 0 = unlimited
	IsActive  bool   `gorm:"default:true"`
	Items     string `gorm:"type:text"` // JSON array of GachaItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GachaItem represents an individual gacha item
type GachaItem struct {
	gorm.Model
	ItemID      int64   `gorm:"uniqueIndex;not null"`
	Name        string  `gorm:"not null"`
	Rarity      string  `gorm:"not null"` // "common", "rare", "epic", "legendary"
	DropRate    float64 `gorm:"not null"` // 0.0 to 1.0
	Category    string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GachaPullResult represents a single gacha pull result
type GachaPullResult struct {
	gorm.Model
	PlayerID int64  `gorm:"not null"`
	PoolID   int64  `gorm:"not null"`
	ItemID   int64  `gorm:"not null"`
	PullID   int64  `gorm:"not null"` // Sequential pull ID for player
	IsNew    bool   `gorm:"default:false"`
	Cost     int64  `gorm:"not null"`
	Currency string `gorm:"not null"`
	PulledAt time.Time
}

// GachaPoolItem represents the relationship between pools and items with specific drop rates
type GachaPoolItem struct {
	gorm.Model
	PoolID   int64   `gorm:"not null"`
	ItemID   int64   `gorm:"not null"`
	DropRate float64 `gorm:"not null"`  // Override default item drop rate for this pool
	Weight   int32   `gorm:"default:1"` // Weight for random selection
}
