package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
	tm   *TokenManager
}

func NewService(repo *Repository, tm *TokenManager) *Service {
	return &Service{repo: repo, tm: tm}
}

func (s *Service) Register(ctx context.Context, email, password, mobile string) (AuthResponse, error) {
	// Default role for new users
	role := "user"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, err
	}

	u, err := s.repo.CreateUser(ctx, email, string(hash), mobile, role)
	if err != nil {
		return AuthResponse{}, err
	}

	return s.issueTokens(ctx, u)
}

func (s *Service) Login(ctx context.Context, email, password string) (AuthResponse, error) {
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}

	return s.issueTokens(ctx, u)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (AuthResponse, error) {
	claims, err := s.tm.ParseRefreshToken(refreshToken)
	if err != nil {
		return AuthResponse{}, err
	}

	rt, err := s.repo.GetRefreshToken(ctx, claims.TokenID)
	if err != nil {
		return AuthResponse{}, err
	}

	if rt.RevokedAt != nil {
		return AuthResponse{}, ErrRefreshTokenRevoked
	}

	// If DB says expired, reject (JWT should also be expired, but DB check is extra safety)
	if time.Now().UTC().After(rt.ExpiresAt) {
		_ = s.repo.RevokeRefreshToken(ctx, rt.ID, time.Now().UTC())
		return AuthResponse{}, errors.New("refresh token expired")
	}

	// Rotate refresh tokens: revoke old token, issue new token row + JWT
	if err := s.repo.RevokeRefreshToken(ctx, rt.ID, time.Now().UTC()); err != nil {
		return AuthResponse{}, err
	}

	u, err := s.repo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return AuthResponse{}, err
	}

	return s.issueTokens(ctx, u)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.tm.ParseRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	// Revoke if not already revoked
	_ = s.repo.RevokeRefreshToken(ctx, claims.TokenID, time.Now().UTC())
	return nil
}

func (s *Service) issueTokens(ctx context.Context, u User) (AuthResponse, error) {
	access, _, err := s.tm.IssueAccessToken(u.ID, u.Email, u.Mobile, u.Role)
	if err != nil {
		return AuthResponse{}, err
	}

	// Create a DB row ID for refresh token
	tokenID := uuid.NewString()

	refreshJWT, refreshExp, err := s.tm.IssueRefreshToken(u.ID, tokenID)
	if err != nil {
		return AuthResponse{}, err
	}

	if err := s.repo.InsertRefreshToken(ctx, tokenID, u.ID, refreshExp); err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		AccessToken:  access,
		RefreshToken: refreshJWT,
	}, nil
}
