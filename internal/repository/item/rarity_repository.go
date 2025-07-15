package repository

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
)

type RarityRepository interface {
	GetAll(ctx context.Context, minRarity *models.Rarity, maxRarity *models.Rarity) ([]models.Rarity, error)
}

type rarityRepository struct {
	db *db.DB
}

func NewRarityRepository(db *db.DB) RarityRepository {
	return &rarityRepository{db: db}
}

func (r *rarityRepository) GetAll(ctx context.Context, minRarity *models.Rarity, maxRarity *models.Rarity) ([]models.Rarity, error) {
	var rarities []models.Rarity

	query := r.db.WithContext(ctx).Order("sort")

	if minRarity != nil {
		query = query.Where("sort >= ?", minRarity.Sort)
	}

	if maxRarity != nil {
		query = query.Where("sort <= ?", maxRarity.Sort)
	}

	if err := query.Find(&rarities).Error; err != nil {
		return nil, fmt.Errorf("RarityRepository::GetAll %v", err)
	}

	return rarities, nil
}
