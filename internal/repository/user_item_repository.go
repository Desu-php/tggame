package repository

import (
	"context"
	"errors"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
	"gorm.io/gorm"
)

type UserItemRepository interface {
	SetUserItem(ctx context.Context, userID uint, item *models.Item, userChestHistory *models.UserChestHistory) (*models.UserItem, error)
	GetLast(ctx context.Context, userID uint) (*models.UserItem, error)
	GetUserItems(ctx context.Context, userID uint) ([]GroupedUserItem, error)
	Exists(ctx context.Context, userID uint, itemID uint) (bool, error)
}

type GroupedUserItem struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Count          int     `json:"count"`
	Type           string  `json:"type"`
	Rarity         string  `json:"rarity"`
	Image          string  `json:"image"`
	Damage         uint    `json:"damage"`
	CriticalDamage uint    `json:"critical_damage"`
	CriticalChance float64 `json:"critical_chance"`
	GoldMultiplier float64 `json:"gold_multiplier"`
	PassiveDamage  uint    `json:"passive_damage"`
}

type userItemRepository struct {
	db *db.DB
}

func NewUserItemRepository(db *db.DB) UserItemRepository {
	return &userItemRepository{db: db}
}

func (r *userItemRepository) SetUserItem(ctx context.Context, userID uint, item *models.Item, userChestHistory *models.UserChestHistory) (*models.UserItem, error) {
	var userChestHistoryID *uint = nil

	if userChestHistory != nil {
		userChestHistoryID = &userChestHistory.ID
	}

	userItem := &models.UserItem{UserID: userID, ItemID: item.ID, UserChestHistoryID: userChestHistoryID}

	result := r.db.WithContext(ctx).Create(userItem)

	if result.Error != nil {
		return nil, fmt.Errorf("UserItemRepository::SetUserItem %w", result.Error)
	}

	return userItem, nil
}

func (r *userItemRepository) GetLast(ctx context.Context, userID uint) (*models.UserItem, error) {
	var userItem models.UserItem

	result := r.db.WithContext(ctx).Preload("Item.Type").
		Preload("Item.Rarity").
		Last(&userItem, "user_id = ?", userID)

	if result.Error != nil {
		return nil, fmt.Errorf("UserItemRepository::GetLast %w", result.Error)
	}

	return &userItem, nil
}

func (r *userItemRepository) GetUserItems(ctx context.Context, userID uint) ([]GroupedUserItem, error) {
	var userItems []GroupedUserItem

	result := r.db.WithContext(ctx).Model(models.UserItem{}).
		Select("i.id, i.name, COUNT(i.id) as count, it.name as type, r.name as rarity, i.image, i.damage, i.critical_damage, i.critical_chance, i.gold_multiplier, i.passive_damage").
		Joins("JOIN items as i ON i.id = user_items.item_id").
		Joins("JOIN item_types as it ON it.id = i.type_id").
		Joins("JOIN rarities as r ON r.id = i.rarity_id").
		Where("user_items.user_id = ?", userID).
		Group("i.id, i.name, type, rarity, image").
		Scan(&userItems)

	if result.Error != nil {
		return nil, fmt.Errorf("UserItemRepository::GetUserItems %w", result.Error)
	}

	return userItems, nil
}

func (r *userItemRepository) Exists(ctx context.Context, userID uint, itemID uint) (bool, error) {
	err := r.db.WithContext(ctx).Model(models.UserItem{}).
		Where("user_id = ? AND item_id = ?", userID, itemID).
		First(&models.UserItem{}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("UserItemRepository::Exists %w", err)
	}

	return true, nil
}
