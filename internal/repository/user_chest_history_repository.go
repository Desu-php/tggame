package repository

import (
	"fmt"

	"example.com/v2/internal/models"
	"gorm.io/gorm"
)

type UserChestHistoryRepository interface {
	Create(userChest *models.UserChest) (*models.UserChestHistory, error)
}

type userChestHistoryRepository struct {
	db *gorm.DB
}

func NewUserChestHistoryRepository(db *gorm.DB) UserChestHistoryRepository {
	return &userChestHistoryRepository{db: db}
}

func (r *userChestHistoryRepository) Create(userChest *models.UserChest) (*models.UserChestHistory, error) {
	userChestHistory := &models.UserChestHistory{
		UserChestID: userChest.ID,
		Health:      userChest.Health,
		Level:       userChest.Level,
	}

	if err := r.db.Create(&userChestHistory).Error; err != nil {
		return nil, fmt.Errorf("UserChestHistoryRepository::Create %w", err)
	}

	return userChestHistory, nil
}
