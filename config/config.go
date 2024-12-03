package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName  string
	AppEnv   string
	AppPort  string
	Database *Database
	Redis *Redis
	Telegram *Telegram
}

type Database struct {
	Host    string
	Port    string
	Name    string
	User    string
	Passord string
}

type Redis struct {
	Host string
	Port string
	Password string
}

type Telegram struct {
	Token string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := &Config{
		AppName: os.Getenv("APP_NAME"),
		AppEnv:  os.Getenv("APP_ENV"),
		AppPort: os.Getenv("APP_PORT"),
		Database: &Database{
			Host:    os.Getenv("DB_HOST"),
			Port:    os.Getenv("DB_PORT"),
			Name:    os.Getenv("DB_DATABASE"),
			User:    os.Getenv("DB_USERNAME"),
			Passord: os.Getenv("DB_PASSWORD"),
		},
		Redis: &Redis{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		Telegram: &Telegram{
			Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
		},
	}

	fmt.Println(cfg)

	return cfg
}
