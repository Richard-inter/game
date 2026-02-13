package repository

import (
	"context"
	"fmt"

	"github.com/1nterdigital/game/internal/domain"
	"gorm.io/gorm"
)

type gachaMachineRepository struct {
	db         *gorm.DB
	defaultRTP float64
}

type GachaMachineRepository interface {
	// player
	CreateGachaPlayer(ctx context.Context, gachaPlayer *domain.GachaPlayer) (*domain.GachaPlayer, error)
	GetGachaPlayerInfo(ctx context.Context, playerID int64) (*domain.GachaPlayer, error)
	AdjustPlayerCoin(ctx context.Context, playerID int64, amount int64, adjustmentType string) (*domain.GachaPlayer, error)
	AdjustPlayerDiamond(ctx context.Context, playerID int64, amount int64, adjustmentType string) (*domain.GachaPlayer, error)
	GetPlayerInventory(ctx context.Context, playerID int64) ([]domain.GachaPlayerInventory, error)
	AddItemInventory(ctx context.Context, playerID int64, itemID int64, quantity int32) error

	// machine
	CreateGachaMachine(ctx context.Context, gachaMachine *domain.GachaMachine) (*domain.GachaMachine, error)
	UpdateGachaMachineItems(ctx context.Context, gachaMachineID int64, items []domain.GachaMachineItem) error
	GetGachaMachineInfo(ctx context.Context, machineID int64) (*domain.GachaMachine, error)
	GetAllGachaMachines(ctx context.Context) ([]*domain.GachaMachine, error)

	// items
	CreateGachaItems(ctx context.Context, items *[]domain.GachaItem) (*[]domain.GachaItem, error)

	// game
	CreateGachaPullSession(ctx context.Context, session *domain.GachaPullSession) (*domain.GachaPullSession, error)
	CreateGachaPullHistories(ctx context.Context, histories *[]domain.GachaPullHistory) (*[]domain.GachaPullHistory, error)
	GetGachaPullSession(ctx context.Context, sessionID int64) (*domain.GachaPullSession, error)
	GetGachaPullHistoriesBySessionID(ctx context.Context, sessionID int64) ([]*domain.GachaPullHistory, error)
	GetPlayerPullHistory(ctx context.Context, playerID int64) ([]*domain.GachaPullSession, error)

	// pity state
	GetGachaPityState(ctx context.Context, playerID int64, machineID int64) (*domain.GachaPityState, error)
	GetAllGachaPityStatesForPlayer(ctx context.Context, playerID int64) ([]*domain.GachaPityState, error)
	SetGachaPityState(ctx context.Context, pityState *domain.GachaPityState) error

	// rtp
	GetGachaMachineRTPState(ctx context.Context, machineID int64) (*domain.GachaMachineRTPState, error)
	InitGachaMachineRTPState(ctx context.Context, machineID int64, targetRTP float64) error
	UpdateGachaMachineRTP(ctx context.Context, machineID int64, price int64, payout int64) error
	UpdateGachaMachineTargetRTP(ctx context.Context, machineID int64, targetRTP float64) error
}

func NewGachaMachineRepository(db *gorm.DB, defaultRTP float64) GachaMachineRepository {
	return &gachaMachineRepository{
		db:         db,
		defaultRTP: defaultRTP,
	}
}

func (r *gachaMachineRepository) CreateGachaPlayer(ctx context.Context, gachaPlayer *domain.GachaPlayer) (*domain.GachaPlayer, error) {
	err := r.db.WithContext(ctx).Create(gachaPlayer).Error
	if err != nil {
		return nil, err
	}
	return gachaPlayer, nil
}

func (r *gachaMachineRepository) GetGachaPlayerInfo(ctx context.Context, playerID int64) (*domain.GachaPlayer, error) {
	var gachaPlayer domain.GachaPlayer
	err := r.db.WithContext(ctx).Where("player_id = ?", playerID).First(&gachaPlayer).Error
	if err != nil {
		return nil, err
	}
	return &gachaPlayer, nil
}

func (r *gachaMachineRepository) adjustPlayerBalance(ctx context.Context, playerID int64, amount int64, adjustmentType, field string) (*domain.GachaPlayer, error) {
	if adjustmentType != "plus" && adjustmentType != "minus" {
		return nil, fmt.Errorf("invalid adjustment type: %s", adjustmentType)
	}

	if adjustmentType == "minus" {
		amount = -amount
	}

	tx := r.db.WithContext(ctx).Model(&domain.GachaPlayer{}).
		Where("player_id = ?", playerID).
		Where(fmt.Sprintf("%s + ? >= 0", field), amount).
		UpdateColumn(field, gorm.Expr(fmt.Sprintf("%s + ?", field), amount))

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		var exists bool
		if err := r.db.WithContext(ctx).Model(&domain.GachaPlayer{}).
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
	if err := r.db.WithContext(ctx).First(&updatedPlayer, "player_id = ?", playerID).Error; err != nil {
		return nil, err
	}

	return &updatedPlayer, nil
}

func (r *gachaMachineRepository) AdjustPlayerCoin(ctx context.Context, playerID int64, amount int64, adjustmentType string) (*domain.GachaPlayer, error) {
	return r.adjustPlayerBalance(ctx, playerID, amount, adjustmentType, "coin")
}

func (r *gachaMachineRepository) AdjustPlayerDiamond(ctx context.Context, playerID int64, amount int64, adjustmentType string) (*domain.GachaPlayer, error) {
	return r.adjustPlayerBalance(ctx, playerID, amount, adjustmentType, "diamond")
}

func (r *gachaMachineRepository) GetPlayerInventory(ctx context.Context, playerID int64) ([]domain.GachaPlayerInventory, error) {
	var inventory []domain.GachaPlayerInventory
	err := activeQuery(r.db.WithContext(ctx)).
		Where("player_id = ? AND is_active = ?", playerID, true).
		Find(&inventory).Error
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *gachaMachineRepository) AddItemInventory(ctx context.Context, playerID int64, itemID int64, quantity int32) error {
	var inventory domain.GachaPlayerInventory
	err := activeQuery(r.db.WithContext(ctx)).
		Where("player_id = ? AND item_id = ?", playerID, itemID).
		First(&inventory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			inventory = domain.GachaPlayerInventory{
				PlayerID: playerID,
				ItemID:   itemID,
				Quantity: quantity,
				IsActive: true,
			}
			return r.db.WithContext(ctx).Create(&inventory).Error
		}
		return err
	}

	inventory.Quantity += quantity
	return r.db.WithContext(ctx).Save(&inventory).Error
}

func (r *gachaMachineRepository) CreateGachaMachine(
	ctx context.Context,
	gachaMachine *domain.GachaMachine,
) (*domain.GachaMachine, error) {
	tx := r.db.WithContext(ctx).Begin()
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

	if err := r.db.WithContext(ctx).
		Preload("Items.Item").
		First(gachaMachine, gachaMachine.ID).Error; err != nil {
		return gachaMachine, nil
	}

	return gachaMachine, nil
}

func (r *gachaMachineRepository) UpdateGachaMachineItems(ctx context.Context, gachaMachineID int64, items []domain.GachaMachineItem) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()

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

func (r *gachaMachineRepository) GetGachaMachineInfo(ctx context.Context, machineID int64) (*domain.GachaMachine, error) {
	var gachaMachine domain.GachaMachine
	err := r.db.WithContext(ctx).Preload("Items.Item").Where("id = ?", machineID).First(&gachaMachine).Error
	if err != nil {
		return nil, err
	}
	return &gachaMachine, nil
}

func (r *gachaMachineRepository) GetAllGachaMachines(ctx context.Context) ([]*domain.GachaMachine, error) {
	var gachaMachines []*domain.GachaMachine
	err := r.db.WithContext(ctx).Preload("Items.Item").Find(&gachaMachines).Error
	if err != nil {
		return nil, err
	}
	return gachaMachines, nil
}

func (r *gachaMachineRepository) CreateGachaItems(ctx context.Context, items *[]domain.GachaItem) (*[]domain.GachaItem, error) {
	err := r.db.WithContext(ctx).Create(items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *gachaMachineRepository) CreateGachaPullSession(ctx context.Context, session *domain.GachaPullSession) (*domain.GachaPullSession, error) {
	err := r.db.WithContext(ctx).Create(session).Error
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *gachaMachineRepository) CreateGachaPullHistories(ctx context.Context, histories *[]domain.GachaPullHistory) (*[]domain.GachaPullHistory, error) {
	err := r.db.WithContext(ctx).Create(histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *gachaMachineRepository) GetGachaPullSession(ctx context.Context, sessionID int64) (*domain.GachaPullSession, error) {
	var session domain.GachaPullSession
	err := r.db.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *gachaMachineRepository) GetPlayerPullHistory(ctx context.Context, playerID int64) ([]*domain.GachaPullSession, error) {
	var sessions []*domain.GachaPullSession
	err := activeQuery(r.db.WithContext(ctx)).
		Preload("GachaPullHistories.Item").
		Where("player_id = ?", playerID).
		Order("created_at DESC").
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *gachaMachineRepository) GetGachaPullHistoriesBySessionID(ctx context.Context, sessionID int64) ([]*domain.GachaPullHistory, error) {
	var histories []*domain.GachaPullHistory
	err := activeQuery(r.db.WithContext(ctx)).
		Preload("Item").
		Where("gacha_pull_session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *gachaMachineRepository) GetGachaPityState(
	ctx context.Context,
	playerID int64,
	machineID int64,
) (*domain.GachaPityState, error) {

	pityState := domain.GachaPityState{
		PlayerID:       playerID,
		GachaMachineID: machineID,
	}

	err := r.db.WithContext(ctx).
		Where("player_id = ? AND gacha_machine_id = ?", playerID, machineID).
		FirstOrCreate(&pityState).Error

	if err != nil {
		return nil, err
	}

	return &pityState, nil
}

func (r *gachaMachineRepository) GetAllGachaPityStatesForPlayer(ctx context.Context, playerID int64) ([]*domain.GachaPityState, error) {
	var pityStates []*domain.GachaPityState
	err := r.db.WithContext(ctx).
		Where("player_id = ?", playerID).
		Find(&pityStates).Error

	if err != nil {
		return nil, err
	}
	return pityStates, nil
}

func (r *gachaMachineRepository) GetGachaMachineRTPState(
	ctx context.Context,
	machineID int64,
) (*domain.GachaMachineRTPState, error) {
	var state domain.GachaMachineRTPState
	err := activeQuery(r.db.WithContext(ctx)).
		Where("gacha_machine_id = ?", machineID).
		First(&state).Error

	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *gachaMachineRepository) InitGachaMachineRTPState(
	ctx context.Context,
	machineID int64,
	targetRTP float64,
) error {
	return r.db.WithContext(ctx).
		Create(&domain.GachaMachineRTPState{
			GachaMachineID: machineID,
			TargetRTP:      targetRTP,
		}).Error
}

func (r *gachaMachineRepository) UpdateGachaMachineRTP(
	ctx context.Context,
	machineID int64,
	price int64,
	payout int64,
) error {
	return activeQuery(r.db.WithContext(ctx)).
		Model(&domain.GachaMachineRTPState{}).
		Where("gacha_machine_id = ?", machineID).
		Updates(map[string]any{
			"total_plays":   gorm.Expr("total_plays + 1"),
			"total_revenue": gorm.Expr("total_revenue + ?", price),
			"total_payout":  gorm.Expr("total_payout + ?", payout),
		}).Error
}

func (r *gachaMachineRepository) UpdateGachaMachineTargetRTP(
	ctx context.Context,
	machineID int64,
	targetRTP float64,
) error {
	return activeQuery(r.db.WithContext(ctx)).
		Model(&domain.GachaMachineRTPState{}).
		Where("gacha_machine_id = ?", machineID).
		Update("target_rtp", targetRTP).Error
}

func (r *gachaMachineRepository) SetGachaPityState(ctx context.Context, pityState *domain.GachaPityState) error {
	return r.db.WithContext(ctx).Model(&domain.GachaPityState{}).
		Where("player_id = ? AND gacha_machine_id = ?", pityState.PlayerID, pityState.GachaMachineID).
		Updates(map[string]interface{}{
			"ultra_rare_pity_count": pityState.UltraRarePityCount,
			"super_rare_pity_count": pityState.SuperRarePityCount,
		}).Error
}
