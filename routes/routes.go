package routes

import (
	"example.com/v2/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	gameController *controllers.GameController,
	) {
		
	api := r.Group("/api")
	{
		api.POST("/start", gameController.Start)
	}
}