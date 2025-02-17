package repository

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
)

type UserItemRepository interface {
	SetUserItem(ctx context.Context, userID uint, item *models.Item, userChestHistory *models.UserChestHistory) (*models.UserItem, error)
	GetLast(ctx context.Context, userID uint) (*models.UserItem, error)
	GetUserItems(ctx context.Context, userID uint) ([]GroupedUserItem, error)
}

type GroupedUserItem struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Count  int    `json:"count"`
	Type   string `json:"type"`
	Rarity string `json:"rarity"`
	Image  string `json:"image"`
}

type userItemRepository struct {
	db *db.DB
}

func NewUserItemRepository(db *db.DB) UserItemRepository {
	return &userItemRepository{db: db}
}

func (r *userItemRepository) SetUserItem(ctx context.Context, userID uint, item *models.Item, userChestHistory *models.UserChestHistory) (*models.UserItem, error) {
	userItem := &models.UserItem{UserID: userID, ItemID: item.ID, UserChestHistoryID: userChestHistory.ID}

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
		Select("i.id, i.name, COUNT(i.id) as count, it.name as type, r.name as rarity, i.image").
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
