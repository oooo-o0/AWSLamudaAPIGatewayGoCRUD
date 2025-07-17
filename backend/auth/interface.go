package auth

import (
	"context"
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrUnauthorized       = errors.New("unauthorized")
)

// AuthService は認証機能の抽象インターフェース
type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)  // JWTトークンを返す
	VerifyToken(ctx context.Context, token string) (string, error)         // トークン検証しユーザーID返す
	RefreshToken(ctx context.Context, refreshToken string) (string, error) // トークン再発行
	GetRoles(ctx context.Context, userID string) ([]string, error)         // ユーザーロール取得
}
