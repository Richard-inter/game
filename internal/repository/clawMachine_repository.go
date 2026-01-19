package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Richard-inter/game/internal/domain"
)

type clawMachineRepository struct {
	db *gorm.DB
}

type ClawMachineRepository interface {
	// player
	CreateClawPlayer(clawPlayer *domain.ClawPlayer) (*domain.ClawPlayer, error)
	GetClawPlayerInfo(playerID int64) (*domain.ClawPlayer, error)
	AdjustPlayerCoin(playerID int64, amount int64, adjustmentType string) (*domain.ClawPlayer, error)
	AdjustPlayerDiamond(playerID int64, amount int64, adjustmentType string) (*domain.ClawPlayer, error)

	// game
	AddGameHistory(playerID int64, gameRecord *domain.ClawMachineGameRecord) (int64, error)
	AddTouchedItemRecord(gameID int64, itemID int64, catched bool) error

	// machine
	CreateClawMachine(clawMachine *domain.ClawMachine) (*domain.ClawMachine, error)
	UpdateClawMachineItems(clawMachineID int64, items []domain.ClawMachineItem) error
	GetClawMachineInfo(machineID int64) (*domain.ClawMachine, error)
	GetAllClawMachines() ([]*domain.ClawMachine, error)

	// items
	CreateClawItems(items *[]domain.ClawItem) (*[]domain.ClawItem, error)
}

func NewClawMachineRepository(db *gorm.DB) ClawMachineRepository {
	return &clawMachineRepository{db: db}
}

func (r *clawMachineRepository) CreateClawPlayer(clawPlayer *domain.ClawPlayer) (*domain.ClawPlayer, error) {
	err := r.db.Create(clawPlayer).Error
	if err != nil {
		return nil, err
	}
	return clawPlayer, nil
}

func (r *clawMachineRepository) GetClawPlayerInfo(playerID int64) (*domain.ClawPlayer, error) {
	var clawPlayer domain.ClawPlayer
	err := r.db.Where("player_id = ?", playerID).First(&clawPlayer).Error
	if err != nil {
		return nil, err
	}
	return &clawPlayer, nil
}

func (r *clawMachineRepository) adjustPlayerBalance(playerID int64, amount int64, adjustmentType, field string) (*domain.ClawPlayer, error) {
	if adjustmentType != "plus" && adjustmentType != "minus" {
		return nil, fmt.Errorf("invalid adjustment type: %s", adjustmentType)
	}

	if adjustmentType == "minus" {
		amount = -amount
	}

	tx := r.db.Model(&domain.ClawPlayer{}).
		Where("player_id = ?", playerID).
		Where(fmt.Sprintf("%s + ? >= 0", field), amount).
		UpdateColumn(field, gorm.Expr(fmt.Sprintf("%s + ?", field), amount))

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		var exists bool
		if err := r.db.Model(&domain.ClawPlayer{}).
			Select("1").
			Where("player_id = ?", playerID).
			Limit(1).
			Scan(&exists).Error; err != nil {
			return nil, err
		}

		if !exists {
			return nil, fmt.Errorf("player not found")
		}
		return nil, fmt.Errorf("not enough %s", field)
	}

	var updatedPlayer domain.ClawPlayer
	if err := r.db.First(&updatedPlayer, "player_id = ?", playerID).Error; err != nil {
		return nil, err
	}

	return &updatedPlayer, nil
}

func (r *clawMachineRepository) AdjustPlayerCoin(playerID int64, amount int64, adjustmentType string) (*domain.ClawPlayer, error) {
	return r.adjustPlayerBalance(playerID, amount, adjustmentType, "coin")
}

func (r *clawMachineRepository) AdjustPlayerDiamond(playerID int64, amount int64, adjustmentType string) (*domain.ClawPlayer, error) {
	return r.adjustPlayerBalance(playerID, amount, adjustmentType, "diamond")
}

func (r *clawMachineRepository) AddGameHistory(playerID int64, gameRecord *domain.ClawMachineGameRecord) (int64, error) {
	err := r.db.Create(gameRecord).Error
	if err != nil {
		return 0, err
	}

	// Get the created record with the generated ID
	var createdRecord domain.ClawMachineGameRecord
	err = r.db.First(&createdRecord, gameRecord.ID).Error
	if err != nil {
		return 0, err
	}

	return createdRecord.ID, nil
}

func (r *clawMachineRepository) AddTouchedItemRecord(gameID int64, itemID int64, catched bool) error {
	err := r.db.Model(&domain.ClawMachineGameRecord{}).
		Where("id = ?", gameID).
		Updates(map[string]any{
			"touched_item_id": itemID,
			"catched":         catched,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *clawMachineRepository) CreateClawMachine(
	clawMachine *domain.ClawMachine,
) (*domain.ClawMachine, error) {
	tx := r.db.Begin()
	if err := tx.Omit("Items").Create(clawMachine).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range clawMachine.Items {
		clawMachine.Items[i].ID = 0
		clawMachine.Items[i].ClawMachineID = clawMachine.ID

		if err := tx.Create(&clawMachine.Items[i]).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.db.
		Preload("Items.Item").
		First(clawMachine, clawMachine.ID).Error; err != nil {
		return clawMachine, nil
	}

	return clawMachine, nil
}

func (r *clawMachineRepository) UpdateClawMachineItems(clawMachineID int64, items []domain.ClawMachineItem) error {
	// Start a transaction
	tx := r.db.Begin()

	// Delete existing items
	if err := tx.Where("claw_machine_id = ?", clawMachineID).Delete(&domain.ClawMachineItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insert new items
	for _, item := range items {
		item.ClawMachineID = clawMachineID
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *clawMachineRepository) GetClawMachineInfo(machineID int64) (*domain.ClawMachine, error) {
	var clawMachine domain.ClawMachine
	err := r.db.Preload("Items.Item").Where("id = ?", machineID).First(&clawMachine).Error
	if err != nil {
		return nil, err
	}
	return &clawMachine, nil
}

func (r *clawMachineRepository) GetAllClawMachines() ([]*domain.ClawMachine, error) {
	var clawMachines []*domain.ClawMachine
	err := r.db.Preload("Items.Item").Find(&clawMachines).Error
	if err != nil {
		return nil, err
	}
	return clawMachines, nil
}

func (r *clawMachineRepository) CreateClawItems(items *[]domain.ClawItem) (*[]domain.ClawItem, error) {
	err := r.db.Create(items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
