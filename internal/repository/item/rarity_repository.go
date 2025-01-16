package repository

import (
	"fmt"

	"example.com/v2/internal/models"
	"gorm.io/gorm"
)

type RarityRepository interface {
	GetAll() ([]models.Rarity, error)
}

type rarityRepository struct {
	db *gorm.DB
}

func NewRarityRepository(db *gorm.DB) RarityRepository {
	return &rarityRepository{db: db}
}

func (r *rarityRepository) GetAll() ([]models.Rarity, error) {
	var rarities []models.Rarity

	if err := r.db.Order("sort").Find(&rarities).Error; err != nil {
		return nil, fmt.Errorf("RarityRepository::GetAll %v", err)
	}

	return rarities, nil
}
