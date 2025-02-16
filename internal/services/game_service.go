package services

import (
	"context"
	"example.com/v2/internal/adapter"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/dto"
	"example.com/v2/pkg/transaction"
	"fmt"
)

type GameService struct {
	userService        *UserService
	transactionManager transaction.TransactionManager
	cacheAdapter       adapter.UserSessionAdapter
}

func NewGameService(userService *UserService, txManager transaction.TransactionManager, adapter adapter.UserSessionAdapter) *GameService {
	return &GameService{userService: userService, transactionManager: txManager, cacheAdapter: adapter}
}

func (g *GameService) Start(ctx context.Context, dto *dto.GameStartDto) (*models.User, error) {
	var result *models.User

	err := g.transactionManager.RunInTransaction(ctx, func(ctx context.Context) error {
		user, err := g.userService.FirstOrCreateByTgId(ctx, dto)
		if err != nil {
			return fmt.Errorf("start: err %w", err)
		}

		user, err = g.userService.GenerateSession(ctx, user)
		if err != nil {
			return fmt.Errorf("start: err %w", err)
		}

		result = user
		return nil
	})

	if err != nil {
		return nil, err
	}

	err = g.cacheAdapter.Set(result.TelegramID, result.Session)

	if err != nil {
		return nil, fmt.Errorf("start: err %w", err)
	}

	return result, nil
}
