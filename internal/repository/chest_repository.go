package repository

import (
	"context"
	"errors"
	"example.com/v2/pkg/db"
	"fmt"

	"example.com/v2/internal/models"
	"gorm.io/gorm"
)

type ChestRepository interface {
	GetDefault(ctx context.Context) (*models.Chest, error)
	FindByUserChest(ctx context.Context, userChest *models.UserChest) (*models.Chest, error)
	GetNextChest(ctx context.Context, currentLevel uint) (*models.Chest, error)
}

type chestRepository struct {
	db *db.DB
}

func NewChestRepository(db *db.DB) ChestRepository {
	return &chestRepository{db: db}
}

func (r *chestRepository) GetDefault(ctx context.Context) (*models.Chest, error) {
	var chest models.Chest

	result := r.db.WithContext(ctx).
		Preload("Rarity").
		Preload("MaxRarity").
		Where(&models.Chest{IsDefault: true}).First(&chest)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("ChestRepository::GetDefault: err %w", result.Error)
	}

	return &chest, nil
}

func (r *chestRepository) GetNextChest(ctx context.Context, currentLevel uint) (*models.Chest, error) {
	var chest models.Chest

	dbContext := r.db.WithContext(ctx)

	result := dbContext.Model(&models.Chest{}).
		Preload("Rarity").
		Preload("MaxRarity").
		Where("start_level <= ?", currentLevel).
		Where("end_level >= ?", currentLevel).
		First(&chest)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result = dbContext.Model(&models.Chest{}).
			Preload("Rarity").
			Preload("MaxRarity").
			Order("end_level desc").
			First(&chest)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		if result.Error != nil {
			return nil, fmt.Errorf("ChestRepository::GetNextChest: err %w", result.Error)
		}

		return &chest, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("ChestRepository::GetNextChest: err %w", result.Error)
	}

	return &chest, nil
}

func (r *chestRepository) FindByUserChest(ctx context.Context, userChest *models.UserChest) (*models.Chest, error) {
	var chest models.Chest

	result := r.db.WithContext(ctx).First(&chest, "id = ?", userChest.ChestID)

	if result.Error != nil {
		return nil, fmt.Errorf("ChestRepository::FindByUserChest %w", result.Error)
	}

	return &chest, nil
}
