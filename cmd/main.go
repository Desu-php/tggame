package main

import (
	"context"
	"example.com/v2/internal/commands"
	"example.com/v2/internal/cron"
	"example.com/v2/pkg"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"example.com/v2/config"
	"example.com/v2/internal/adapter"
	"example.com/v2/internal/controllers"
	"example.com/v2/internal/redis"
	"example.com/v2/internal/repository"
	"example.com/v2/internal/services"
	"example.com/v2/routes"
	"github.com/caddyserver/certmagic"
	"github.com/gin-contrib/cors"
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
		pkg.Module,
		repository.Module,
		adapter.Module,
		services.Module,
		controllers.Module,
		fx.Invoke(routes.RegisterRoutes),
		fx.Invoke(StartServer),
		cron.Module,
		commands.Module,
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

func NewGinEngine(cfg *config.Config) *gin.Engine {
	engine := gin.Default()

	c := cors.Config{
		AllowOrigins:     strings.Split(cfg.AppFrontUrl, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},  // Разрешённые заголовки
		ExposeHeaders:    []string{"*"},  // Заголовки, доступные клиенту
		AllowCredentials: true,           // Для запросов с куками
		MaxAge:           12 * time.Hour, // Время кеширования CORS-политики
	}

	engine.Use(cors.New(c))

	return engine
}

func StartServer(router *gin.Engine, cfg *config.Config) {
	go func() {
		if cfg.IsProduction() || cfg.IsStage() {
			certmagic.DefaultACME.Agreed = true
			certmagic.DefaultACME.Email = cfg.AppEmail

			if cfg.IsProduction() {
				certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
			} else {
				certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
			}

			certmagic.Default.Storage = &certmagic.FileStorage{Path: "./certmagic-storage"}

			domains := []string{cfg.AppDomain}

			err := certmagic.HTTPS(domains, router)
			if err != nil {
				log.Fatalf("Failed to start HTTPS server: %v", err)
			}
		} else {
			if err := router.Run("0.0.0.0:" + cfg.AppPort); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		}
	}()
}
