package repository

import (
	"gorm.io/gorm"

	"github.com/Richard-inter/game/internal/domain"
)

type clawMachineRepository struct {
	db *gorm.DB
}

type ClawMachineRepository interface {
	CreateClawPlayer(clawPlayer *domain.ClawPlayer) (*domain.ClawPlayer, error)
	GetClawPlayerInfo(playerID int64) (*domain.ClawPlayer, error)

	CreateClawMachine(clawMachine *domain.ClawMachine) (*domain.ClawMachine, error)
	UpdateClawMachineItems(clawMachineID int64, items []domain.ClawMachineItem) error

	GetClawMachineInfo(machineID int64) (*domain.ClawMachine, error)

	CreateClawItems(items *[]domain.Item) (*[]domain.Item, error)
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

func (r *clawMachineRepository) CreateClawItems(items *[]domain.Item) (*[]domain.Item, error) {
	err := r.db.Create(items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
