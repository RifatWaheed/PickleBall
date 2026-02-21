package auth

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=320"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Mobile   string `json:"mobile"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=320"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type MeResponse struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
	Role   string `json:"role"`
}
