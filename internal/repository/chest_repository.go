package repository

import (
	"errors"
	"fmt"

	"example.com/v2/internal/models"
	"gorm.io/gorm"
)

type ChestRepository interface {
	GetDefault() (*models.Chest, error)
	FindByUserChest(userChest *models.UserChest) (*models.Chest, error)
	GetNextChest(currentLevel uint) (*models.Chest, error)
}

type chestRepository struct {
	db *gorm.DB
}

func NewChestRepository(db *gorm.DB) ChestRepository {
	return &chestRepository{db: db}
}

func (r *chestRepository) GetDefault() (*models.Chest, error) {
	var chest models.Chest

	result := r.db.Where(&models.Chest{IsDefault: true}).First(&chest)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("ChestRepository::GetDefault: err %w", result.Error)
	}

	return &chest, nil
}

func (r *chestRepository) GetNextChest(currentLevel uint) (*models.Chest, error) {
	var chest models.Chest

	result := r.db.Model(&models.Chest{}).
		Where("start_level <= ?", currentLevel).
		Where("end_level >= ?", currentLevel).
		First(&chest)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result = r.db.Model(&models.Chest{}).Order("end_level desc").First(&chest)

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

func (r *chestRepository) FindByUserChest(userChest *models.UserChest) (*models.Chest, error) {
	var chest models.Chest

	result := r.db.First(&chest, "id = ?", userChest.ChestID)

	if result.Error != nil {
		return nil, fmt.Errorf("ChestRepository::FindByUserChest %w", result.Error)
	}

	return &chest, nil
}
