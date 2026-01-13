package domain

type ClawMachine struct {
	ID      int64  `gorm:"column:id;primaryKey" json:"machineID"`
	Name    string `gorm:"column:name" json:"name"`
	Price   int64  `gorm:"column:price" json:"price"`
	MaxItem int32  `gorm:"column:max_item" json:"maxItem"`

	Items []ClawMachineItem `gorm:"foreignKey:ClawMachineID;constraint:OnDelete:CASCADE"`
}

type ClawMachineItem struct {
	ID            int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ClawMachineID int64 `gorm:"column:claw_machine_id" json:"clawMachineID"`
	ItemID        int64 `gorm:"column:item_id" json:"itemID"`

	Item Item `gorm:"foreignKey:ItemID;references:ID"`
}

type Item struct {
	ID              int64  `gorm:"column:id;primaryKey" json:"itemID"`
	Name            string `gorm:"column:name" json:"name"`
	Rarity          string `gorm:"column:rarity" json:"rarity"`
	SpawnPercentage int64  `gorm:"column:spawn_percentage" json:"spawnPercentage"`
	CatchPercentage int64  `gorm:"column:catch_percentage" json:"catchPercentage"`
	MaxItemSpawned  int64  `gorm:"column:max_item_spawned" json:"maxItemSpawned"`
}

type ClawPlayer struct {
	Player  Player `gorm:"embedded;embeddedPrefix:player_"`
	Coin    int64  `gorm:"column:coin;not null" json:"coin"`
	Diamond int64  `gorm:"column:diamond;not null" json:"diamond"`
}

type ClawMachineGameRecord struct {
	ID            int64 `gorm:"column:id;primaryKey" json:"gameID"`
	ClawMachineID int64 `gorm:"column:claw_machine_id" json:"clawMachineID"`
	PlayerID      int64 `gorm:"column:player_id" json:"playerID"`
	TouchedItemID int64 `gorm:"column:touched_item_id" json:"touchedItemID"`
	Catched       bool  `gorm:"column:catched" json:"catched"`
}

func (ClawMachine) TableName() string {
	return "claw_machine"
}

func (ClawMachineItem) TableName() string {
	return "claw_machine_item"
}

func (Item) TableName() string {
	return "claw_item"
}

func (ClawPlayer) TableName() string {
	return "claw_player"
}

func (ClawMachineGameRecord) TableName() string {
	return "claw_machine_game_record"
}
