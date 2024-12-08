package repository

import (
	"errors"
	"fmt"

	"example.com/v2/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindById(id uint64) (*models.User, error)
	Create(dto CreateUserDTO) (*models.User, error)
	GetAll() ([]models.User, error)
	FindByTgId(id uint64) (*models.User, error)
	UpdateSession(user *models.User, session string) error
	FindWithoutPreloadingByTgId(id uint64) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

type CreateUserDTO struct {
	Username   string
	TelegramID uint64
}

func NewUserRepository(baseRepo *BaseRepository) UserRepository {
	return &userRepository{db: baseRepo.DB}
}

func (r *userRepository) FindById(id uint64) (*models.User, error) {
	var user models.User

	result := r.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("FindById: err %w", result.Error)
	}

	return &user, nil
}

func (r *userRepository) Create(dto CreateUserDTO) (*models.User, error) {
	user := &models.User{
		Username:   dto.Username,
		TelegramID: dto.TelegramID,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("Create: err %w", err)
	}

	return user, nil
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByTgId(id uint64) (*models.User, error) {
	var user models.User

	result := r.db.Preload("UserChest.Chest").First(&user, "telegram_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("FindByTgId: err %w", result.Error)
	}

	return &user, nil
}

func (r *userRepository) FindWithoutPreloadingByTgId(id uint64) (*models.User, error) {
	var user models.User

	result := r.db.First(&user, "telegram_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("FindWithoutPreloadingByTgId: err %w", result.Error)
	}

	return &user, nil
}

func (r *userRepository) UpdateSession(user *models.User, session string) error {
	result := r.db.Model(&user).Update("session", session)
	
	if result.Error != nil {
		return fmt.Errorf("UpdateSession: err %w", result.Error)
	}
	
	return nil
}