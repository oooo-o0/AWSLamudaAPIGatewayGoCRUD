package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// AuthorizationヘッダーからBearerトークンを取り出す
func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("authorization header missing")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil
}

// AuthMiddleware は認証ミドルウェア。LambdaのAPIGatewayProxyRequestを受けてユーザーIDを返す
func AuthMiddleware(ctx context.Context, req events.APIGatewayProxyRequest, svc AuthService) (string, error) {
	authHeader := req.Headers["Authorization"]
	token, err := extractBearerToken(authHeader)
	if err != nil {
		return "", err
	}

	userID, err := svc.VerifyToken(ctx, token)
	if err != nil {
		return "", err
	}
	return userID, nil
}
