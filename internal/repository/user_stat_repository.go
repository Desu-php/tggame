package repository

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
)

type UserStatRepository interface {
	GetStat(ctx context.Context, user *models.User) (*models.UserStat, error)
	Create(ctx context.Context, user *models.User) (*models.UserStat, error)
}

type userStatRepository struct {
	db *db.DB
}

func NewUserStatRepository(db *db.DB) UserStatRepository {
	return &userStatRepository{db: db}
}

func (r *userStatRepository) GetStat(ctx context.Context, user *models.User) (*models.UserStat, error) {
	var userStat models.UserStat

	err := r.db.WithContext(ctx).Model(models.UserStat{}).
		Where("user_id = ?", user.ID).
		Find(&userStat).Error

	if err != nil {
		return nil, fmt.Errorf("UserStatRepository::GetStat %w", err)
	}

	return &userStat, nil
}

func (r *userStatRepository) Create(ctx context.Context, user *models.User) (*models.UserStat, error) {
	userStat := &models.UserStat{UserID: user.ID}

	err := r.db.WithContext(ctx).Create(userStat).Error

	if err != nil {
		return nil, fmt.Errorf("UserStatRepository::Create %w", err)
	}

	return userStat, nil
}
