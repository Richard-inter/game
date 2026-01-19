package repository

import (
	"fmt"

	"github.com/Richard-inter/game/internal/domain"
	"gorm.io/gorm"
)

type gachaMachineRepository struct {
	db *gorm.DB
}

type GachaMachineRepository interface {
	// player
	CreateGachaPlayer(gachaPlayer *domain.GachaPlayer) (*domain.GachaPlayer, error)
	GetGachaPlayerInfo(playerID int64) (*domain.GachaPlayer, error)
	AdjustPlayerCoin(playerID int64, amount int64, adjustmentType string) (*domain.GachaPlayer, error)
	AdjustPlayerDiamond(playerID int64, amount int64, adjustmentType string) (*domain.GachaPlayer, error)

	// machine
	CreateGachaMachine(gachaMachine *domain.GachaMachine) (*domain.GachaMachine, error)
	UpdateGachaMachineItems(gachaMachineID int64, items []domain.GachaMachineItem) error
	GetGachaMachineInfo(machineID int64) (*domain.GachaMachine, error)
	GetAllGachaMachines() ([]*domain.GachaMachine, error)

	// items
	CreateGachaItems(items *[]domain.GachaItem) (*[]domain.GachaItem, error)
}

func NewGachaMachineRepository(db *gorm.DB) GachaMachineRepository {
	return &gachaMachineRepository{
		db: db,
	}
}

func (r *gachaMachineRepository) CreateGachaPlayer(gachaPlayer *domain.GachaPlayer) (*domain.GachaPlayer, error) {
	err := r.db.Create(gachaPlayer).Error
	if err != nil {
		return nil, err
	}
	return gachaPlayer, nil
}

func (r *gachaMachineRepository) GetGachaPlayerInfo(playerID int64) (*domain.GachaPlayer, error) {
	var gachaPlayer domain.GachaPlayer
	err := r.db.Where("player_id = ?", playerID).First(&gachaPlayer).Error
	if err != nil {
		return nil, err
	}
	return &gachaPlayer, nil
}

func (r *gachaMachineRepository) adjustPlayerBalance(playerID int64, amount int64, adjustmentType, field string) (*domain.GachaPlayer, error) {
	if adjustmentType != "plus" && adjustmentType != "minus" {
		return nil, fmt.Errorf("invalid adjustment type: %s", adjustmentType)
	}

	if adjustmentType == "minus" {
		amount = -amount
	}

	tx := r.db.Model(&domain.GachaPlayer{}).
		Where("player_id = ?", playerID).
		Where(fmt.Sprintf("%s + ? >= 0", field), amount).
		UpdateColumn(field, gorm.Expr(fmt.Sprintf("%s + ?", field), amount))

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		var exists bool
		if err := r.db.Model(&domain.GachaPlayer{}).
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

	var updatedPlayer domain.GachaPlayer
	if err := r.db.First(&updatedPlayer, "player_id = ?", playerID).Error; err != nil {
		return nil, err
	}

	return &updatedPlayer, nil
}

func (r *gachaMachineRepository) AdjustPlayerCoin(playerID int64, amount int64, adjustmentType string) (*domain.GachaPlayer, error) {
	return r.adjustPlayerBalance(playerID, amount, adjustmentType, "coin")
}

func (r *gachaMachineRepository) AdjustPlayerDiamond(playerID int64, amount int64, adjustmentType string) (*domain.GachaPlayer, error) {
	return r.adjustPlayerBalance(playerID, amount, adjustmentType, "diamond")
}

func (r *gachaMachineRepository) CreateGachaMachine(
	gachaMachine *domain.GachaMachine,
) (*domain.GachaMachine, error) {
	tx := r.db.Begin()
	if err := tx.Omit("Items").Create(gachaMachine).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	fmt.Println(gachaMachine.ID)
	fmt.Println(gachaMachine.Items)
	for i := range gachaMachine.Items {
		gachaMachine.Items[i].ID = 0
		gachaMachine.Items[i].GachaMachineID = gachaMachine.ID

		if err := tx.Create(&gachaMachine.Items[i]).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.db.
		Preload("Items.Item").
		First(gachaMachine, gachaMachine.ID).Error; err != nil {
		return gachaMachine, nil
	}

	return gachaMachine, nil
}

func (r *gachaMachineRepository) UpdateGachaMachineItems(gachaMachineID int64, items []domain.GachaMachineItem) error {
	// Start a transaction
	tx := r.db.Begin()

	// Delete existing items
	if err := tx.Where("gacha_machine_id = ?", gachaMachineID).Delete(&domain.GachaMachineItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insert new items
	for _, item := range items {
		item.GachaMachineID = gachaMachineID
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *gachaMachineRepository) GetGachaMachineInfo(machineID int64) (*domain.GachaMachine, error) {
	var gachaMachine domain.GachaMachine
	err := r.db.Preload("Items.Item").Where("id = ?", machineID).First(&gachaMachine).Error
	if err != nil {
		return nil, err
	}
	return &gachaMachine, nil
}

func (r *gachaMachineRepository) GetAllGachaMachines() ([]*domain.GachaMachine, error) {
	var gachaMachines []*domain.GachaMachine
	err := r.db.Preload("Items.Item").Find(&gachaMachines).Error
	if err != nil {
		return nil, err
	}
	return gachaMachines, nil
}

func (r *gachaMachineRepository) CreateGachaItems(items *[]domain.GachaItem) (*[]domain.GachaItem, error) {
	err := r.db.Create(items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
