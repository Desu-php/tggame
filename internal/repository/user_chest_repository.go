package repository

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
	"log"
)

type UserChestRepository interface {
	Create(ctx context.Context, userChest *models.UserChest) (*models.UserChest, error)
	DecrementHealth(ctx context.Context, userChest *models.UserChest, damage uint) error
	Update(ctx context.Context, userChest *models.UserChest) error
	FindByUser(ctx context.Context, user *models.User) (*models.UserChest, error)
}

type userChestRepository struct {
	db *db.DB
}

func NewUserChestRepository(db *db.DB) UserChestRepository {
	return &userChestRepository{db: db}
}

func (r *userChestRepository) Create(ctx context.Context, userChest *models.UserChest) (*models.UserChest, error) {
	result := r.db.WithContext(ctx).Create(userChest)

	if result.Error != nil {
		return nil, fmt.Errorf("UserChestRepository::Create: err %w", result.Error)
	}

	return userChest, nil
}

func (r *userChestRepository) DecrementHealth(ctx context.Context, userChest *models.UserChest, damage uint) error {
	userChest.CurrentHealth = userChest.CurrentHealth - int64(damage)

	result := r.db.WithContext(ctx).Save(userChest)

	if result.Error != nil {
		return fmt.Errorf("UserChestRepository::DecrementHealth: err %w", result.Error)
	}

	return nil
}

func (r *userChestRepository) Update(ctx context.Context, userChest *models.UserChest) error {
	result := r.db.WithContext(ctx).Updates(&models.UserChest{
		ID:            userChest.ID,
		ChestID:       userChest.ChestID,
		Level:         userChest.Level,
		CurrentHealth: userChest.CurrentHealth,
		Health:        userChest.Health,
		Amount:        userChest.Amount,
	})

	if result.Error != nil {
		log.Printf("Update error: %s", result.Error)
		return fmt.Errorf("UserChestRepository::Update: err %w", result.Error)
	}

	return nil
}

func (r *userChestRepository) FindByUser(ctx context.Context, user *models.User) (*models.UserChest, error) {
	var userChest models.UserChest

	result := r.db.WithContext(ctx).Preload("Chest.Rarity").First(&userChest, "user_id = ?", user.ID)

	if result.Error != nil {
		return nil, fmt.Errorf("UserChestRepository::FindByUser: err %w", result.Error)
	}

	return &userChest, nil
}
