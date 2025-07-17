package service

import (
	"context"
	"errors"
	"os"
	"time"

	"kaigo-insurance-system/backend/model"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Authenticate(ctx context.Context, email, password string) (*model.AuthResponse, error)
	Signup(ctx context.Context, req *model.SignupRequest) error
	ValidateToken(ctx context.Context, tokenString string) (*jwt.RegisteredClaims, error)
}

type jwtAuthService struct {
	secretKey []byte
}

func NewAuthService() AuthService {
	return &jwtAuthService{
		secretKey: []byte(os.Getenv("JWT_SECRET")),
	}
}

func (a *jwtAuthService) Authenticate(ctx context.Context, email, password string) (*model.AuthResponse, error) {
	// 仮実装：実際はDynamoDBからユーザーを検索してPWハッシュを検証
	if email != "admin@example.com" || password != "securepass" {
		return nil, errors.New("invalid credentials")
	}

	claims := &jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token:   tokenString,
		Expires: claims.ExpiresAt.Time.Format(time.RFC3339),
	}, nil
}

func (a *jwtAuthService) Signup(ctx context.Context, req *model.SignupRequest) error {
	// 仮登録: 実際はユーザーをDynamoDBに保存し、重複チェック
	if req.Email == "" || req.Password == "" {
		return errors.New("missing required fields")
	}
	return nil
}

func (a *jwtAuthService) ValidateToken(ctx context.Context, tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return claims, nil
	}

	return nil, errors.New("could not extract claims")
}
