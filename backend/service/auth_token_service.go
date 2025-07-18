package service

import (
	"context"
	"errors"
	"time"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/repository"
)

type AuthTokenService interface {
	Create(ctx context.Context, token *model.AuthToken) error
	GetByID(ctx context.Context, tokenID string) (*model.AuthToken, error)
	GetByUserID(ctx context.Context, userID string) ([]model.AuthToken, error)
	Revoke(ctx context.Context, tokenID string) error
	CleanupExpired(ctx context.Context, before time.Time) error
}

type authTokenService struct {
	repo repository.AuthTokenRepository
}

func NewAuthTokenService(repo repository.AuthTokenRepository) AuthTokenService {
	return &authTokenService{repo: repo}
}

func (s *authTokenService) Create(ctx context.Context, token *model.AuthToken) error {
	if token.TokenID == "" || token.UserID == "" || token.Token == "" {
		return errors.New("missing required fields")
	}
	if token.Expiry.Before(time.Now()) {
		return errors.New("expiry must be future date")
	}
	return s.repo.Create(ctx, token)
}

func (s *authTokenService) GetByID(ctx context.Context, tokenID string) (*model.AuthToken, error) {
	return s.repo.GetByID(ctx, tokenID)
}

func (s *authTokenService) GetByUserID(ctx context.Context, userID string) ([]model.AuthToken, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *authTokenService) Revoke(ctx context.Context, tokenID string) error {
	return s.repo.Revoke(ctx, tokenID)
}

func (s *authTokenService) CleanupExpired(ctx context.Context, before time.Time) error {
	return s.repo.CleanupExpired(ctx, before)
}
