package services

import (
	"fmt"
	"log"

	"example.com/v2/internal/adapter"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/dto"
	"example.com/v2/pkg/transaction"
)

type GameService struct {
	userService        *UserService
	transactionManager transaction.TransactionManager
	cacheAdapter       adapter.UserSessionAdapter
}

func NewGameService(userService *UserService, txManager transaction.TransactionManager, adapter adapter.UserSessionAdapter) *GameService {
	log.Println("NewGameService")
	return &GameService{userService: userService, transactionManager: txManager, cacheAdapter: adapter}
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

	err = g.cacheAdapter.Set(result.TelegramID, result.Session)

	if err != nil {
		return nil, fmt.Errorf("Start: err %w", err)
	}

	return result, nil
}
