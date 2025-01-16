package controllers

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewGameController,
	NewClickController,
	NewUserItemController,
	NewRarityController,
)
