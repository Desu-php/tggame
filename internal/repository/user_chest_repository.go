package repository

import (
	"fmt"
	"example.com/v2/internal/models"
	"gorm.io/gorm"
)

type UserChestRepository interface {
	Create(userChest *models.UserChest) (*models.UserChest, error)
}

type userChestRepository struct {
	db *gorm.DB
}

func NewUserChestRepository(db *gorm.DB) UserChestRepository {
	return &userChestRepository{db: db}
}

func (r *userChestRepository) Create(userChest *models.UserChest) (*models.UserChest, error) {
	result := r.db.Create(userChest)

	if result.Error != nil {
		return nil, fmt.Errorf("UserChestRepository::Create: err %w", result.Error)
	}

	return userChest, nil
}
