package repository

import (
	"context"
	"example.com/v2/internal/models"
	"fmt"
	"gorm.io/gorm"
)

type BalanceRepository interface {
	Create(ctx context.Context, userId uint) (*models.Balance, error)
	Update(ctx context.Context, userId uint, amount int64) (*models.Balance, error)
	FindByUserId(ctx context.Context, userId uint) (*models.Balance, error)
}

type balanceRepository struct {
	db *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) BalanceRepository {
	return &balanceRepository{db: db}
}

func (r *balanceRepository) Create(ctx context.Context, userId uint) (*models.Balance, error) {
	balance := &models.Balance{
		UserID:  userId,
		Balance: 0,
	}

	err := r.db.WithContext(ctx).Create(balance)

	if err != nil {
		return nil, fmt.Errorf("BalanceRepository::Create :%v", err)
	}

	return balance, nil
}

func (r *balanceRepository) Update(ctx context.Context, userId uint, amount int64) (*models.Balance, error) {
	dbCtx := r.db.WithContext(ctx)

	balance, err := r.FindByUserId(ctx, userId)

	if err != nil {
		return nil, fmt.Errorf("BalanceRepository::Update :%v", err)
	}

	balance.Balance += amount

	err = dbCtx.Save(&balance).Error

	if err != nil {
		return nil, fmt.Errorf("BalanceRepository::Update :%v", err)
	}

	return balance, nil
}

func (r *balanceRepository) FindByUserId(ctx context.Context, userId uint) (*models.Balance, error) {
	balance := models.Balance{}

	err := r.db.WithContext(ctx).Where("user_id = ?", userId).First(&balance).Error

	if err != nil {
		return nil, fmt.Errorf("BalanceRepository::Update :%v", err)
	}

	return &balance, nil
}
