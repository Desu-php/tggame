package repository

import (
	"example.com/v2/config"
	balance "example.com/v2/internal/repository/balance"
	itemService "example.com/v2/internal/repository/item"
	"example.com/v2/pkg/db"
	"fmt"
	"go.uber.org/fx"
)

func NewGORMInstance(cfg *config.Config) (*db.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	return db.NewDB(dsn)
}

var Module = fx.Provide(
	NewGORMInstance,
	NewUserRepository,
	NewChestRepository,
	NewUserChestRepository,
	NewUserChestHistoryRepository,
	itemService.NewItemRepository,
	itemService.NewRarityRepository,
	NewUserItemRepository,
	NewReferralUserRepository,
	NewUserStatRepository,
	NewBalanceRepository,
	balance.NewTransactionRepository,
	balance.NewBalanceRepository,
	NewAspectRepository,
)
