package auth

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type cognitoService struct {
	client      *cognito.Client
	userPoolID  string
	appClientID string
}

func NewCognitoService(cfg aws.Config, userPoolID, appClientID string) AuthService {
	return &cognitoService{
		client:      cognito.NewFromConfig(cfg),
		userPoolID:  userPoolID,
		appClientID: appClientID,
	}
}

func (c *cognitoService) Login(ctx context.Context, username, password string) (string, error) {
	input := &cognito.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
		ClientId: &c.appClientID,
	}

	resp, err := c.client.InitiateAuth(ctx, input)
	if err != nil {
		return "", err
	}

	if resp.AuthenticationResult == nil {
		return "", errors.New("authentication failed")
	}

	return *resp.AuthenticationResult.IdToken, nil
}

func (c *cognitoService) VerifyToken(ctx context.Context, token string) (string, error) {
	// Cognitoトークン検証は別途実装が必要（JWKS取得や署名検証）
	return "", errors.New("not implemented")
}

func (c *cognitoService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return "", errors.New("not implemented")
}

func (c *cognitoService) GetRoles(ctx context.Context, userID string) ([]string, error) {
	// Cognitoユーザーグループなどからロール取得
	return []string{}, errors.New("not implemented")
}
