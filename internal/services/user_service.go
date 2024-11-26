package services

import (
	"fmt"

	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/pkg/dto"
	"example.com/v2/pkg/str"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) FirstOrCreateByTgId(dto *dto.GameStartDto) (*models.User, error) {
	user, err := u.repo.FindByTgId(dto.TelegramId)

	if err != nil {
		return nil, fmt.Errorf("FirstOrCreateByTgId: %w", err)
	}

	if user != nil {
		return user, nil
	}

	user, err = u.repo.Create(repository.CreateUserDTO{TelegramID: dto.TelegramId, Username: dto.Username})

	if err != nil {
		return nil, fmt.Errorf("FirstOrCreateByTgId: err %w", err)
	}

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
