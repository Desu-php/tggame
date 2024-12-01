package routes

import (
	"net/http"

	"example.com/v2/internal/adapter"
	"example.com/v2/internal/controllers"
	"example.com/v2/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RegisterRoutes(
	r *gin.Engine,
	gameController *controllers.GameController,
	sessionAdapter adapter.UserSessionAdapter,
	logger *logrus.Logger,
	clickController *controllers.ClickController,
	) {
		
	game := r.Group("/api/game")
	{
		game.POST("/start", gameController.Start)
	}

	api := r.Group("/api")
	api.Use(middleware.SessionMiddleware(sessionAdapter, logger))
	{
		api.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK,  gin.H{"Success": true})
		})

		api.POST("click", clickController.Store)
		api.GET("clicks", clickController.Get)
	}
}