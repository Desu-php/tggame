package repository

import (
	"fmt"
	"example.com/v2/internal/models"
	"gorm.io/gorm"
)

type UserChestRepository interface {
	Create(userChest *models.UserChest) (*models.UserChest, error)
	DecrementHealth(userChest *models.UserChest, damage uint) error
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

func (r *userChestRepository) DecrementHealth(userChest *models.UserChest, damage uint) error{
	userChest.CurrentHealth = userChest.CurrentHealth - int(damage)

	result := r.db.Save(userChest)

	if result.Error != nil {
		return fmt.Errorf("UserChestRepository::DecrementHealth: err %w", result.Error)
	}

	return nil
}
