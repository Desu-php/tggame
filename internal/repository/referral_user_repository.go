package repository

import (
	"example.com/v2/internal/models"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type ReferralUserRepository interface {
	GetByUserID(userID uint) ([]models.ReferralUser, error)
	Create(referralUserID uint, userID uint) error
}

type referralUserRepository struct {
	db *gorm.DB
}

func NewReferralUserRepository(db *gorm.DB) ReferralUserRepository {
	return &referralUserRepository{db: db}
}

func (r *referralUserRepository) GetByUserID(userID uint) ([]models.ReferralUser, error) {
	var referralUsers []models.ReferralUser

	err := r.db.Model(&models.ReferralUser{}).
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

func (r *referralUserRepository) Create(referralUserID uint, userID uint) error {
	err := r.db.Create(&models.ReferralUser{
		UserID:         userID,
		ReferredUserID: referralUserID,
	}).Error

	if err != nil {
		return fmt.Errorf("ReferralUserRepository::Create %v", err)
	}

	return nil
}
