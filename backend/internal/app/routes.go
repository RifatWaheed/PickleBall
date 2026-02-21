package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"pickleBall/internal/auth"
)

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")

	if err := auth.RegisterRoutes(v1, db); err != nil {
		// if JWT secrets missing etc, crash early
		panic(err)
	}
}
