package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool) error {
	tm, err := NewTokenManagerFromEnv()
	if err != nil {
		return err
	}

	repo := NewRepository(db)
	svc := NewService(repo, tm)
	h := NewHandler(svc)

	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)

		// Example protected route
		auth.GET("/me", RequireAuth(tm), h.Me)
	}

	return nil
}
