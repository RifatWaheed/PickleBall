package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const accessClaimsKey = "access_claims"

func RequireAuth(tm *TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims, err := tm.ParseAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(accessClaimsKey, claims)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetAccessClaims(c)
		if claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if claims.Role != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

func GetAccessClaims(c *gin.Context) *AccessClaims {
	v, ok := c.Get(accessClaimsKey)
	if !ok {
		return nil
	}
	claims, ok := v.(*AccessClaims)
	if !ok {
		return nil
	}
	return claims
}
