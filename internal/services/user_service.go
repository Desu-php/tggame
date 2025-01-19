package services

import (
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/pkg/dto"
	"example.com/v2/pkg/str"
	"example.com/v2/pkg/transaction"
	"fmt"
)

type UserService struct {
	repo                   repository.UserRepository
	userChestService       *UserChestService
	transaction            transaction.TransactionManager
	referralUserRepository repository.ReferralUserRepository
}

func NewUserService(
	repo repository.UserRepository,
	userChestService *UserChestService,
	transaction transaction.TransactionManager,
	referralUserRepository repository.ReferralUserRepository,
) *UserService {
	return &UserService{
		repo:                   repo,
		userChestService:       userChestService,
		transaction:            transaction,
		referralUserRepository: referralUserRepository,
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

	err = u.transaction.RunInTransaction(func() error {
		user, err = u.repo.Create(repository.CreateUserDTO{TelegramID: dto.TelegramId, Username: dto.Username})

		if err != nil {
			return fmt.Errorf("FirstOrCreateByTgId: err %w", err)
		}

		if dto.ReferrerId != nil {
			referrer, err := u.repo.FindById(uint64(*dto.ReferrerId))

			if err != nil {
				return fmt.Errorf("FirstOrCreateByTgId: err %w", err)
			}

			if referrer != nil {
				err = u.referralUserRepository.Create(user.ID, referrer.ID)

				if err != nil {
					return fmt.Errorf("FirstOrCreateByTgId: err %w", err)
				}
			}
		}

		userChest, err := u.userChestService.Create(user)

		if err != nil {
			return fmt.Errorf("FirstOrCreateByTgId: err %w", err)
		}

		user.UserChest = *userChest

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) GenerateSession(user *models.User) (*models.User, error) {
	session, err := str.GenerateSessionKey()

	if err != nil {
		return nil, fmt.Errorf("GenerateSession: err %w", err)
	}

	err = u.repo.UpdateSession(user, session)

	if err != nil {
		return nil, fmt.Errorf("GenerateSession: err %w", err)
	}

	user.Session = session

	return user, nil
}
