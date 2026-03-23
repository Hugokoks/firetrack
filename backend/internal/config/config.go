package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	FilesRoot   string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		FilesRoot:   getEnv("FILES_ROOT", "./uploads"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
