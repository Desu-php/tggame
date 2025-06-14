package repository

import (
	"context"
	"example.com/v2/pkg/db"
	"fmt"

	"example.com/v2/internal/models"
)

type ItemRepository interface {
	GetAllByRarity(ctx context.Context, rarity *models.Rarity) ([]models.Item, error)
}

type itemRepository struct {
	db *db.DB
}

func NewItemRepository(db *db.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) GetAllByRarity(ctx context.Context, rarity *models.Rarity) ([]models.Item, error) {
	var items []models.Item

	if err := r.db.WithContext(ctx).
		Where("quantity > 0").
		Where("rarity_id = ?", rarity.ID).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch items: %v", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no items found for rarity: %v", rarity.Name)
	}

	return items, nil
}
