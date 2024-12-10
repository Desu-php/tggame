package services

import (
	itemService "example.com/v2/internal/services/item"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewUserService,
	NewGameService,
	NewUserChestService,
	NewClickService,
	itemService.NewRarityService,
	itemService.NewItemService,
)
