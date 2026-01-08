package domain

type ClawMachine struct {
	ID    int64  `gorm:"column:id;primaryKey" json:"machineID"`
	Name  string `gorm:"column:name" json:"name"`
	Price int64  `gorm:"column:price" json:"price"`

	Items []ClawMachineItem `gorm:"foreignKey:ClawMachineID;constraint:OnDelete:CASCADE"`
}

type ClawMachineItem struct {
	ID            int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ClawMachineID int64 `gorm:"column:claw_machine_id" json:"clawMachineID"`
	ItemID        int64 `gorm:"column:item_id" json:"itemID"`

	Item Item `gorm:"foreignKey:ItemID;references:ID"`
}

type Item struct {
	ID     int64  `gorm:"column:id;primaryKey" json:"itemID"`
	Name   string `gorm:"column:name" json:"name"`
	Rarity string `gorm:"column:rarity" json:"rarity"`
	Weight int64  `gorm:"column:weight" json:"weight"`
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
	Success       bool  `gorm:"column:success" json:"success"`

	ItemRecords []ClawMachineItemGameRecord `gorm:"foreignKey:GameID"`
}

type ClawMachineItemGameRecord struct {
	ID      int64 `gorm:"column:id;primaryKey" json:"id"`
	GameID  int64 `gorm:"column:game_id" json:"gameID"`
	ItemID  int64 `gorm:"column:item_id" json:"itemID"`
	Success bool  `gorm:"column:success" json:"success"`

	Game ClawMachineGameRecord `gorm:"foreignKey:GameID"`
	Item Item                  `gorm:"foreignKey:ItemID"`
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

func (ClawMachineItemGameRecord) TableName() string {
	return "claw_machine_item_game_record"
}
