package services

import (
	"fmt"

	"example.com/v2/internal/models"
	"example.com/v2/pkg/dto"
	"example.com/v2/pkg/transaction"
)

type GameService struct{
	userService *UserService
	transactionManager transaction.TransactionManager
}

func NewGameService(userService *UserService, txManager transaction.TransactionManager) *GameService{
	return &GameService{userService: userService, transactionManager: txManager}
}

func (g *GameService) Start(dto *dto.GameStartDto) (*models.User, error) {
	var result *models.User

	err := g.transactionManager.RunInTransaction(func() error {
		user, err := g.userService.FirstOrCreateByTgId(dto)
		if err != nil {
			return fmt.Errorf("Start: err %w", err)
		}

		user, err = g.userService.GenerateSession(user)
		if err != nil {
			return fmt.Errorf("Start: err %w", err)
		}

		result = user
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}