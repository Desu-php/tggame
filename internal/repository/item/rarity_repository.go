package repository

import (
	"context"
	"fmt"

	"example.com/v2/internal/models"
	"gorm.io/gorm"
)

type RarityRepository interface {
	GetAll(ctx context.Context) ([]models.Rarity, error)
}

type rarityRepository struct {
	db *gorm.DB
}

func NewRarityRepository(db *gorm.DB) RarityRepository {
	return &rarityRepository{db: db}
}

func (r *rarityRepository) GetAll(ctx context.Context) ([]models.Rarity, error) {
	var rarities []models.Rarity

	if err := r.db.WithContext(ctx).Order("sort").Find(&rarities).Error; err != nil {
		return nil, fmt.Errorf("RarityRepository::GetAll %v", err)
	}

	return rarities, nil
}
