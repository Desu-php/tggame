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

	api.Use(middleware.SessionMiddleware(sessionAdapter, logger, userRepository))
	{
		api.POST("click", clickController.Store)

		api.GET("user/item/last", userItemController.GetLast)
		api.GET("user/items", userItemController.GetUserItems)
	}
}
