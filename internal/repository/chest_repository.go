package repository

import (
	"errors"
	"fmt"

	"example.com/v2/internal/models"
	"gorm.io/gorm"
)


type ChestRepository interface {
	GetDefault() (*models.Chest, error)
}

type chestRepository struct {
	db *gorm.DB
}

func NewChestRespository(db *gorm.DB) ChestRepository {
	return &chestRepository{db: db}
}

func (r *chestRepository) GetDefault() (*models.Chest, error){
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