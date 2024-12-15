package services

import (
	"fmt"
	"log"

	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/pkg/dto"
	"example.com/v2/pkg/str"
	"example.com/v2/pkg/transaction"
)

type UserService struct {
	repo             repository.UserRepository
	userChestService *UserChestService
	transaction      transaction.TransactionManager
}

func NewUserService(
	repo repository.UserRepository,
	userChestService *UserChestService,
	transaction transaction.TransactionManager,
) *UserService {
	return &UserService{
		repo:             repo,
		userChestService: userChestService,
		transaction:      transaction,
	}
}

func (u *UserService) FirstOrCreateByTgId(dto *dto.GameStartDto) (*models.User, error) {
	user, err := u.repo.FindByTgId(dto.TelegramId)

	if err != nil {
		return nil, fmt.Errorf("FirstOrCreateByTgId: %w", err)
	}

	if user != nil {
		return user, nil
	}

	u.transaction.RunInTransaction(func() error {
		user, err = u.repo.Create(repository.CreateUserDTO{TelegramID: dto.TelegramId, Username: dto.Username})

		if err != nil {
			return fmt.Errorf("FirstOrCreateByTgId: err %w", err)
		}

		userChest, err := u.userChestService.Create(user)

		log.Println(userChest)

		if err != nil {
			return fmt.Errorf("FirstOrCreateByTgId: err %w", err)
		}

		user.UserChest = *userChest

		return  nil
	})

	return user, nil
}

func (u UserService) GenerateSession(user *models.User) (*models.User, error) {
	session, err := str.GenerateSessionKey()

	if err != nil {
		return nil, fmt.Errorf("GenerateSession: err %w", err)
	}

	u.repo.UpdateSession(user, session)
	user.Session = session

	return user, nil
}
