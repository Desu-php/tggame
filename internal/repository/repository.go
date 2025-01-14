package repository

import (
	"example.com/v2/config"
	itemService "example.com/v2/internal/repository/item"
	"example.com/v2/pkg/db"
	"fmt"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewGORMInstance(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	return db.SetupGorm(dsn)
}

var Module = fx.Provide(
	NewGORMInstance,
	NewBaseRepository,
	NewUserRepository,
	NewChestRepository,
	NewUserChestRepository,
	NewUserChestHistoryRepository,
	itemService.NewItemRepository,
	itemService.NewRarityRepository,
	NewUserItemRepository,
)
