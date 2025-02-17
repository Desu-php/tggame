package repository

import (
	"context"
	"example.com/v2/pkg/db"
	"fmt"

	"example.com/v2/internal/models"
)

type UserChestHistoryRepository interface {
	Create(ctx context.Context, userChest *models.UserChest) (*models.UserChestHistory, error)
}

type userChestHistoryRepository struct {
	db *db.DB
}

func NewUserChestHistoryRepository(db *db.DB) UserChestHistoryRepository {
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
