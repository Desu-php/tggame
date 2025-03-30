package repository

import (
	"context"
	"errors"
	"example.com/v2/pkg/db"
	"fmt"
	"gorm.io/gorm"

	"example.com/v2/internal/models"
)

type UserChestHistoryRepository interface {
	Create(ctx context.Context, userChest *models.UserChest) (*models.UserChestHistory, error)
	LastAmount(ctx context.Context, user *models.User) (uint32, error)
	Last(ctx context.Context, user *models.User) (*models.UserChestHistory, error)
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
		Amount:      userChest.Amount,
	}

	if err := r.db.WithContext(ctx).Create(&userChestHistory).Error; err != nil {
		return nil, fmt.Errorf("UserChestHistoryRepository::Create %w", err)
	}

	return userChestHistory, nil
}

func (r *userChestHistoryRepository) LastAmount(ctx context.Context, user *models.User) (uint32, error) {
	userHistory, err := r.Last(ctx, user)

	if err != nil {
		return 0, fmt.Errorf("UserChestHistoryRepository::LastAmount %w", err)
	}

	if userHistory == nil {
		return 0, nil
	}

	return userHistory.Amount, nil
}

func (r *userChestHistoryRepository) Last(ctx context.Context, user *models.User) (*models.UserChestHistory, error) {
	var userChest models.UserChestHistory

	result := r.db.WithContext(ctx).
		Select("user_chest_histories.*").
		Joins("inner join user_chests as uc on uc.id = user_chest_histories.user_chest_id").
		Where("uc.user_id = ?", user.ID).
		Order("user_chest_histories.id desc").
		First(&userChest)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("UserChestHistoryRepository::Last %w", result.Error)
	}

	return &userChest, nil
}
