package repository

import (
	"example.com/v2/internal/models"
	"fmt"
	"gorm.io/gorm"
)

type UserChestRepository interface {
	Create(userChest *models.UserChest) (*models.UserChest, error)
	DecrementHealth(userChest *models.UserChest, damage uint) error
	Update(userChest *models.UserChest) error
	FindByUser(user *models.User) (*models.UserChest, error)
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

func (r *userChestRepository) DecrementHealth(userChest *models.UserChest, damage uint) error {
	userChest.CurrentHealth = userChest.CurrentHealth - int(damage)

	result := r.db.Save(userChest)

	if result.Error != nil {
		return fmt.Errorf("UserChestRepository::DecrementHealth: err %w", result.Error)
	}

	return nil
}

func (r *userChestRepository) Update(userChest *models.UserChest) error {
	result := r.db.Save(userChest)

	if result.Error != nil {
		return fmt.Errorf("UserChestRepository::Update: err %w", result.Error)
	}

	return nil
}

func (r userChestRepository) FindByUser(user *models.User) (*models.UserChest, error) {
	var userChest models.UserChest

	result := r.db.Preload("Chest").First(&userChest, "user_id = ?", user.ID)

	if result.Error != nil {
		return nil, fmt.Errorf("UserChestRepository::FindByUser: err %w", result.Error)
	}

	return &userChest, nil
}
