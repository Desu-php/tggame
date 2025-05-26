package routes

import (
	"net/http"

	"example.com/v2/config"
	"example.com/v2/internal/adapter"
	"example.com/v2/internal/controllers"
	"example.com/v2/internal/middleware"
	"example.com/v2/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RegisterRoutes(
	r *gin.Engine,
	gameController *controllers.GameController,
	sessionAdapter adapter.UserSessionAdapter,
	logger *logrus.Logger,
	clickController *controllers.ClickController,
	userRepository repository.UserRepository,
	cfg *config.Config,
	userItemController *controllers.UserItemController,
	rarityController *controllers.RarityController,
	referralController *controllers.ReferralController,
	userController *controllers.UserController,
	boosterController *controllers.BoosterController,
	userAspectController *controllers.UserAspectController,
	aspectController *controllers.AspectController,
	craftController *controllers.CraftController,
	itemController *controllers.ItemController,
	taskController *controllers.TaskController,
) {

	game := r.Group("/api/game")
	{
		game.Use(middleware.TelegramHashCheck(cfg))
		game.POST("/start", gameController.Start)
	}

	r.GET("api/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"Success": true})
	})

	api := r.Group("/api")

	api.GET("rarities", rarityController.GetRarities)

	api.GET("top/users", userController.GetTop)

	api.GET("items", itemController.GetItems)

	api.Use(middleware.SessionMiddleware(sessionAdapter, logger, userRepository))
	{
		api.POST("click", clickController.Store)

		api.GET("user/item/last", userItemController.GetLast)
		api.GET("user/items", userItemController.GetUserItems)
		api.GET("user/referrals", referralController.GetReferrals)
		api.GET("user/referrals/count", referralController.GetReferralCount)
		api.GET("user/info", userController.Info)
		api.GET("user/boosters", userAspectController.GetBoosters)

		api.GET("boosters/:type", boosterController.Index)
		api.POST("booster/:id/buy", boosterController.Buy)
		api.PUT("booster/:id/upgrade", boosterController.Upgrade)
		api.POST("aspect/:id/active", aspectController.Store)

		api.POST("craft", craftController.Craft)

		api.GET("tasks", taskController.GetAll)
		api.POST("task/:id/click/link", taskController.ClickLink)
		api.GET("task/:id/reward", taskController.ReceiveReward)
	}
}
