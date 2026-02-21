package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	fmt.Println("REGISTER DEBUG -> email:", req.Email, "mobile:", req.Mobile)

	res, err := h.svc.Register(c.Request.Context(), req.Email, req.Password, req.Mobile)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrEmailAlreadyExists {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	res, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	_ = h.svc.Logout(c.Request.Context(), req.RefreshToken)
	// Always return OK to avoid leaking token validity info
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) Me(c *gin.Context) {
	claims := GetAccessClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, MeResponse{
		UserID: claims.UserID,
		Email:  claims.Email,
		Mobile: "", // not in claims (by choice); fetch from DB later if needed
		Role:   claims.Role,
	})
}
