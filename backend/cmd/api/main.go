package main

import (
	"context"
	"log"
	"time"

	"pickleBall/configs"
	"pickleBall/internal/app"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := configs.Load()

	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := app.NewDBPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// ✅ Pass db into routes (we will use it in auth)
	app.RegisterRoutes(router, db)

	addr := ":" + cfg.Port
	log.Printf("server running on %s\n", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
