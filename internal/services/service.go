package services

import "go.uber.org/fx"

var Module = fx.Provide(
	NewUserService,
	NewGameService,
	NewUserChestService,
	NewClickService,
)