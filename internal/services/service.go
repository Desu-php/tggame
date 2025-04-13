package services

import (
	authService "example.com/v2/internal/services/auth"
	itemService "example.com/v2/internal/services/item"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewBalanceService,
	NewUserService,
	NewGameService,
	NewUserChestService,
	NewClickService,
	itemService.NewRarityService,
	itemService.NewItemService,
	NewUserItemService,
	authService.NewAuthService,
	NewUserStatService,
	NewUserAspectService,
)
