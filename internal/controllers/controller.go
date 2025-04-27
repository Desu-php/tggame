package controllers

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewGameController,
	NewClickController,
	NewUserItemController,
	NewRarityController,
	NewReferralController,
	NewUserController,
	NewBoosterController,
	NewUserAspectController,
	NewAspectController,
	NewCraftController,
	NewItemController,
)
