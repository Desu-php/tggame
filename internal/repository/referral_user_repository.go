package repository

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
	"log"
)

type ReferralUserRepository interface {
	GetByUserID(ctx context.Context, userID uint) ([]models.ReferralUser, error)
	Create(ctx context.Context, referralUserID uint, userID uint) error
	Count(ctx context.Context, userID uint) (uint, error)
}

type referralUserRepository struct {
	db *db.DB
}

func NewReferralUserRepository(db *db.DB) ReferralUserRepository {
	return &referralUserRepository{db: db}
}

func (r *referralUserRepository) GetByUserID(ctx context.Context, userID uint) ([]models.ReferralUser, error) {
	var referralUsers []models.ReferralUser

	err := r.db.WithContext(ctx).Model(&models.ReferralUser{}).
		Preload("ReferredUser").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&referralUsers).
		Error

	if err != nil {
		return nil, fmt.Errorf("ReferralUserRepository::GetByUserID %v", err)
	}

	log.Println(referralUsers)

	return referralUsers, nil
}

func (r *referralUserRepository) Create(ctx context.Context, referralUserID uint, userID uint) error {
	err := r.db.WithContext(ctx).Create(&models.ReferralUser{
		UserID:         userID,
		ReferredUserID: referralUserID,
	}).Error

	if err != nil {
		return fmt.Errorf("ReferralUserRepository::Create %v", err)
	}

	return nil
}

func (r *referralUserRepository) Count(ctx context.Context, userID uint) (uint, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.ReferralUser{}).
		Where("user_id = ?", userID).
		Count(&count).
		Error

	if err != nil {
		return 0, fmt.Errorf("ReferralUserRepository::Create %v", err)
	}

	return uint(count), nil
}
