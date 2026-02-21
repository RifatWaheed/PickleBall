package auth

import (
	"os"
	"strconv"
	"time"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	if d, err := time.ParseDuration(v); err == nil {
		return d
	}

	if n, err := strconv.Atoi(v); err == nil {
		return time.Duration(n) * time.Minute
	}

	return fallback
}
