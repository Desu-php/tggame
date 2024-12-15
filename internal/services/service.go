package services

import (
	itemService "example.com/v2/internal/services/item"
	"go.uber.org/fx"
	authService "example.com/v2/internal/services/auth"
)

var Module = fx.Provide(
	NewUserService,
	NewGameService,
	NewUserChestService,
	NewClickService,
	itemService.NewRarityService,
	itemService.NewItemService,
	NewUserItemService,
	authService.NewAuthService,
)
