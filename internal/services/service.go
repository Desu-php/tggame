package services

import (
	authService "example.com/v2/internal/services/auth"
	itemService "example.com/v2/internal/services/item"
	"example.com/v2/internal/services/notification"
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
	NewTaskService,
	NewReferralService,
	notification.NewTelegramChannel,
)
