package services

import (
	"fmt"
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
)

type UserItemService struct {
	userItemRepository repository.UserItemRespository
}

func NewUserItemService(userItemRepository repository.UserItemRespository) *UserItemService {
	return &UserItemService{userItemRepository: userItemRepository}
}

func (s *UserItemService) SetUserItem(userId uint, item *models.Item, userChestHistory *models.UserChestHistory) (*models.UserItem, error){
	userItem, err := s.userItemRepository.SetUserItem(userId, item, userChestHistory)

	if err != nil {
		return nil, fmt.Errorf("UserItemService::SetUserItem %w", err)
	}

	return userItem, nil
}
