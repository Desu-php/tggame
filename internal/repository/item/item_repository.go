package repository

import (
	"fmt"

	"example.com/v2/internal/models"
	"gorm.io/gorm"
)

type ItemRepository interface {
	GetAllByRarity(rarity *models.Rarity) ([]models.Item, error)
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) GetAllByRarity(rarity *models.Rarity) ([]models.Item, error){
	var items []models.Item

	if err := r.db.Where("rarity_id = ?", rarity.ID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch items: %v", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no items found for rarity: %v", rarity.Name)
	}

	return items, nil
}