package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	AppName  string
	AppEnv   string
	AppPort  string
	Database *Database
}

type Database struct {
	Host    string
	Port    string
	Name    string
	User    string
	Passord string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
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
	}
}
