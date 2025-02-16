package repository

import (
	"context"
	"fmt"

	"example.com/v2/internal/models"
	"gorm.io/gorm"
)

type UserChestHistoryRepository interface {
	Create(ctx context.Context, userChest *models.UserChest) (*models.UserChestHistory, error)
}

type userChestHistoryRepository struct {
	db *gorm.DB
}

func NewUserChestHistoryRepository(db *gorm.DB) UserChestHistoryRepository {
	return &userChestHistoryRepository{db: db}
}

func (r *userChestHistoryRepository) Create(ctx context.Context, userChest *models.UserChest) (*models.UserChestHistory, error) {
	userChestHistory := &models.UserChestHistory{
		UserChestID: userChest.ID,
		Health:      userChest.Health,
		Level:       userChest.Level,
	}

	if err := r.db.WithContext(ctx).Create(&userChestHistory).Error; err != nil {
		return nil, fmt.Errorf("UserChestHistoryRepository::Create %w", err)
	}

	return userChestHistory, nil
}
