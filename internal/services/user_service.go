package services

import (
	"context"
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

func (u *UserService) FirstOrCreateByTgId(ctx context.Context, dto *dto.GameStartDto) (*models.User, error) {
	user, err := u.repo.FindByTgId(ctx, dto.TelegramId)

	if err != nil {
		return nil, fmt.Errorf("FirstOrCreateByTgId: %w", err)
	}

	if user != nil {
		return user, nil
	}

	err = u.transaction.RunInTransaction(ctx, func(ctx context.Context) error {
		user, err = u.repo.Create(ctx, repository.CreateUserDTO{TelegramID: dto.TelegramId, Username: dto.Username})

		if err != nil {
			return fmt.Errorf("FirstOrCreateByTgId: err %w", err)
		}

		if dto.ReferrerId != nil {
			referrer, err := u.repo.FindById(ctx, uint64(*dto.ReferrerId))

			if err != nil {
				return fmt.Errorf("FirstOrCreateByTgId: err %w", err)
			}

			if referrer != nil {
				err = u.referralUserRepository.Create(ctx, user.ID, referrer.ID)

				if err != nil {
					return fmt.Errorf("FirstOrCreateByTgId: err %w", err)
				}
			}
		}

		userChest, err := u.userChestService.Create(ctx, user)

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

func (u *UserService) GenerateSession(ctx context.Context, user *models.User) (*models.User, error) {
	session, err := str.GenerateSessionKey()

	if err != nil {
		return nil, fmt.Errorf("GenerateSession: err %w", err)
	}

	err = u.repo.UpdateSession(ctx, user, session)

	if err != nil {
		return nil, fmt.Errorf("GenerateSession: err %w", err)
	}

	user.Session = session

	return user, nil
}
