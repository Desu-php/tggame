package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/v2/config"
	"example.com/v2/internal/adapter"
	"example.com/v2/internal/controllers"
	"example.com/v2/internal/redis"
	"example.com/v2/internal/repository"
	"example.com/v2/internal/services"
	"example.com/v2/pkg"
	"example.com/v2/routes"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.LoadConfig,
			NewGinEngine,
			redis.NewRedisClient,
		),
		repository.Module,
		adapter.Module,
		services.Module,
		controllers.Module,
		pkg.Module,
		fx.Invoke(routes.RegisterRoutes),
		fx.Invoke(StartServer),
		fx.Invoke(func ()  {
			go processClicks()
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down application...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		log.Fatalf("Application failed to stop gracefully: %v", err)
	}

	log.Println("Application stopped")
}

func NewGinEngine() *gin.Engine {
	return gin.Default()
}

func StartServer(router *gin.Engine, cfg *config.Config) {
	go func() {
		log.Printf("Starting server on port %s", cfg.AppPort)
		if err := router.Run("localhost:" + cfg.AppPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
}

func processClicks() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
}
