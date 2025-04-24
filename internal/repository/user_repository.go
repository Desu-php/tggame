package repository

import (
	"context"
	"errors"
	"example.com/v2/pkg/db"
	"example.com/v2/pkg/transaction"
	"fmt"

	"example.com/v2/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindById(ctx context.Context, id uint64) (*models.User, error)
	Create(ctx context.Context, dto CreateUserDTO) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	FindByTgId(ctx context.Context, id uint64) (*models.User, error)
	UpdateSession(ctx context.Context, user *models.User, session string) error
	FindWithoutPreloadingByTgId(ctx context.Context, id uint64) (*models.User, error)
	GetTop(ctx context.Context) ([]TopUserDTO, error)
}

type userRepository struct {
	db                 *db.DB
	userStatRepository UserStatRepository
	transaction        transaction.TransactionManager
}

type CreateUserDTO struct {
	Username   string
	TelegramID uint64
}

type TopUserDTO struct {
	Username string `json:"username"`
	Level    uint   `json:"level"`
}

func NewUserRepository(
	db *db.DB,
	userStatRepository UserStatRepository,
	transaction transaction.TransactionManager,
) UserRepository {
	return &userRepository{
		db:                 db,
		userStatRepository: userStatRepository,
		transaction:        transaction,
	}
}

func (r *userRepository) FindById(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User

	result := r.db.WithContext(ctx).First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("FindById: err %w", result.Error)
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, dto CreateUserDTO) (*models.User, error) {
	user := &models.User{
		Username:   dto.Username,
		TelegramID: dto.TelegramID,
	}

	err := r.transaction.RunInTransaction(ctx,
		func(ctx context.Context) error {
			if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
				return fmt.Errorf("UserRepository::Create err %w", err)
			}

			_, err := r.userStatRepository.Create(ctx, user)

			if err != nil {
				return fmt.Errorf("UserRepository::Create: err %w", err)
			}

			return nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("UserRepository::Create: err %w", err)
	}

	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByTgId(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User

	result := r.db.WithContext(ctx).Preload("UserChest.Chest.Rarity").First(&user, "telegram_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("FindByTgId: err %w", result.Error)
	}

	return &user, nil
}

func (r *userRepository) FindWithoutPreloadingByTgId(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User

	result := r.db.WithContext(ctx).First(&user, "telegram_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("FindWithoutPreloadingByTgId: err %w", result.Error)
	}

	return &user, nil
}

func (r *userRepository) UpdateSession(ctx context.Context, user *models.User, session string) error {
	result := r.db.WithContext(ctx).Model(&user).Update("session", session)

	if result.Error != nil {
		return fmt.Errorf("UpdateSession: err %w", result.Error)
	}

	return nil
}

func (r *userRepository) GetTop(ctx context.Context) ([]TopUserDTO, error) {
	var users []TopUserDTO

	err := r.db.WithContext(ctx).
		Model(models.UserChest{}).
		Select("u.username, user_chests.level").
		Joins("inner join users as u on u.id = user_chests.user_id").
		Order("user_chests.level desc").
		Limit(100).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("GetTop: err %w", err)
	}

	return users, nil
}
