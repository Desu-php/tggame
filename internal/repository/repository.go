package repository

import (
	"fmt"

	"example.com/v2/config"
	"example.com/v2/pkg/db"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewGORMInstance(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
	cfg.Database.User,
	cfg.Database.Passord,
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
	NewChestRespository,
	NewUserChestRepository,
	NewUserChestHistoryRepository,
)
