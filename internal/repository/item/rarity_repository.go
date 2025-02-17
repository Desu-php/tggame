package repository

import (
	"context"
	"example.com/v2/pkg/db"
	"fmt"

	"example.com/v2/internal/models"
)

type RarityRepository interface {
	GetAll(ctx context.Context) ([]models.Rarity, error)
}

type rarityRepository struct {
	db *db.DB
}

func NewRarityRepository(db *db.DB) RarityRepository {
	return &rarityRepository{db: db}
}

func (r *rarityRepository) GetAll(ctx context.Context) ([]models.Rarity, error) {
	var rarities []models.Rarity

	if err := r.db.WithContext(ctx).Order("sort").Find(&rarities).Error; err != nil {
		return nil, fmt.Errorf("RarityRepository::GetAll %v", err)
	}

	return rarities, nil
}
