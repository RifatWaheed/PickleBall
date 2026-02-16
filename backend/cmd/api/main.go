package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	// 1️⃣ Load configuration
	cfg := configs.Load()

	// 2️⃣ Set Gin mode (debug/release)
	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	// 3️⃣ Create DB connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := app.NewDBPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	// 4️⃣ Create Gin engine
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 5️⃣ Register routes (we will pass db soon for auth module)
	app.RegisterRoutes(router)

	// 6️⃣ Start server
	address := ":" + cfg.Port
	log.Printf("server running on %s\n", address)

	if err := router.Run(address); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
