package repository

import (
	"example.com/v2/internal/models"
	"fmt"
	"gorm.io/gorm"
)

type UserStatRepository interface {
	GetStat(user *models.User) (*models.UserStat, error)
	Create(user *models.User) (*models.UserStat, error)
}

type userStatRepository struct {
	db *gorm.DB
}

func NewUserStatRepository(db *gorm.DB) UserStatRepository {
	return &userStatRepository{db: db}
}

func (r *userStatRepository) GetStat(user *models.User) (*models.UserStat, error) {
	var userStat models.UserStat

	err := r.db.Model(models.UserStat{}).
		Where("user_id = ?", user.ID).
		Find(&userStat).Error

	if err != nil {
		return nil, fmt.Errorf("UserStatRepository::GetStat %w", err)
	}

	return &userStat, nil
}

func (r *userStatRepository) Create(user *models.User) (*models.UserStat, error) {
	userStat := &models.UserStat{UserID: user.ID}

	err := r.db.Create(userStat).Error

	if err != nil {
		return nil, fmt.Errorf("UserStatRepository::Create %w", err)
	}

	return userStat, nil
}
