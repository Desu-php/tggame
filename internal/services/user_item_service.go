package services

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"fmt"
)

type UserItemService struct {
	userItemRepository repository.UserItemRepository
}

func NewUserItemService(userItemRepository repository.UserItemRepository) *UserItemService {
	return &UserItemService{userItemRepository: userItemRepository}
}

func (s *UserItemService) SetUserItem(ctx context.Context, userId uint, item *models.Item, userChestHistory *models.UserChestHistory) (*models.UserItem, error) {
	userItem, err := s.userItemRepository.SetUserItem(ctx, userId, item, userChestHistory)

	if err != nil {
		return nil, fmt.Errorf("UserItemService::SetUserItem %w", err)
	}

	return userItem, nil
}
