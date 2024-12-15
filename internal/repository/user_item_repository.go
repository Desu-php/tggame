package repository

import (
	"example.com/v2/internal/models"
	"fmt"
	"gorm.io/gorm"
)

type UserItemRespository interface {
	SetUserItem(userID uint, item *models.Item, userChestHistory *models.UserChestHistory) (*models.UserItem, error)
	GetLast(userID uint) (*models.UserItem, error)
	GetUserItems(userID uint) ([]GroupedUserItem, error)
}

type GroupedUserItem struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Count  int    `json:"count"`
	Type   string `json:"type"`
	Rarity string `json:"rarity"`
	Image  string `json:"image"`
}

type userItemRespository struct {
	db *gorm.DB
}

func NewUserItemRepository(db *gorm.DB) UserItemRespository {
	return &userItemRespository{db: db}
}

func (r *userItemRespository) SetUserItem(userID uint, item *models.Item, userChestHistory *models.UserChestHistory) (*models.UserItem, error) {
	userItem := &models.UserItem{UserID: userID, ItemID: item.ID, UserChestHistoryID: userChestHistory.ID}

	result := r.db.Create(userItem)

	if result.Error != nil {
		return nil, fmt.Errorf("UserItemRespository::SetUserItem %w", result.Error)
	}

	return userItem, nil
}

func (r *userItemRespository) GetLast(userID uint) (*models.UserItem, error) {
	var userItem models.UserItem

	result := r.db.Preload("Item.Type").
		Preload("Item.Rarity").
		Last(&userItem, "user_id = ?", userID)

	if result.Error != nil {
		return nil, fmt.Errorf("UserItemRespository::GetLast %w", result.Error)
	}

	return &userItem, nil
}

func (r *userItemRespository) GetUserItems(userID uint) ([]GroupedUserItem, error) {
	var userItems []GroupedUserItem

	result := r.db.Model(models.UserItem{}).
		Select("i.id, i.name, COUNT(i.id) as count, it.name as type, r.name as rarity, i.image").
		Joins("JOIN items as i ON i.id = user_items.item_id").
		Joins("JOIN item_types as it ON it.id = i.type_id").
		Joins("JOIN rarities as r ON r.id = i.rarity_id").
		Where("user_items.user_id = ?", userID).
		Group("i.id, i.name, type, rarity, image").
		Scan(&userItems)

	if result.Error != nil {
		return nil, fmt.Errorf("UserItemRespository::GetUserItems %w", result.Error)
	}

	return userItems, nil
}
