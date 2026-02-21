package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	issuer        string
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration //TTL = Time To Live or Expiration duration
	refreshTTL    time.Duration
}

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"uid"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UserID  string `json:"uid"`
	TokenID string `json:"tid"` // refresh token row id (uuid string)
}

func NewTokenManagerFromEnv() (*TokenManager, error) {
	issuer := getEnv("JWT_ISSUER", "pickleball-api")

	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")

	if accessSecret == "" || refreshSecret == "" {
		return nil, errors.New("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be set")
	}

	accessTTL := getEnvDuration("JWT_ACCESS_TTL", 15*time.Minute)
	refreshTTL := getEnvDuration("JWT_REFRESH_TTL", 30*24*time.Hour)

	return &TokenManager{
		issuer:        issuer,
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}, nil
}

func (tm *TokenManager) IssueAccessToken(userID, email, mobile, role string) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(tm.accessTTL)

	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tm.issuer,
			Subject:   userID,
			Audience:  []string{"access"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		UserID: userID,
		Email:  email,
		Mobile: mobile,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(tm.accessSecret)
	return s, exp, err
}

func (tm *TokenManager) IssueRefreshToken(userID, tokenID string) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(tm.refreshTTL)

	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tm.issuer,
			Subject:   userID,
			Audience:  []string{"refresh"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		UserID:  userID,
		TokenID: tokenID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(tm.refreshSecret)
	return s, exp, err
}

func (tm *TokenManager) ParseAccessToken(tokenString string) (*AccessClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (any, error) {
		return tm.accessSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(*AccessClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

func (tm *TokenManager) ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (any, error) {
		return tm.refreshSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(*RefreshClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}
