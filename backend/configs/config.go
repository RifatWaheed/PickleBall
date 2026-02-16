package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	GinMode     string
}

func Load() Config {
	_ = godotenv.Load() // ok if missing in prod

	cfg := Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: mustGetEnv("DATABASE_URL"),
		GinMode:     getEnv("GIN_MODE", "debug"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is not set", key)
	}
	return v
}
