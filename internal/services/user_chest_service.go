package services

import (
	"fmt"

	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
)

type UserChestService struct {
	userChestRepo repository.UserChestRepository
	chestRepo     repository.ChestRepository
}

func NewUserChestService(
	userChestRepo repository.UserChestRepository,
	chestRepo repository.ChestRepository,
) *UserChestService {
	return &UserChestService{userChestRepo: userChestRepo, chestRepo: chestRepo}
}

func(s *UserChestService) Create(user *models.User) (*models.UserChest, error) {
   chest, err := s.chestRepo.GetDefault()

   if err != nil {
	return nil, fmt.Errorf("UserChestService::Create err %w", err)
   }

   if chest == nil {
	return nil, fmt.Errorf("UserChestService::Create default chest not found")	
   }

  userChest, err := s.userChestRepo.Create(&models.UserChest{
		UserID: user.ID,
		ChestID: chest.ID,
		CurrentHealth: chest.Health,
		Level: 1,
   })

   if err != nil {
	return nil, fmt.Errorf("UserChestService::Create err %w", err)
   }

   userChest.Chest = *chest

   return userChest, nil
}
