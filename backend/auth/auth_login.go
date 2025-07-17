package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"kaigo-insurance-system/backend/auth"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// LoginRequest はログインAPIのリクエストボディ
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse はログイン成功時のレスポンス
type LoginResponse struct {
	Token string `json:"token"`
}

var authSvc auth.AuthService

func init() {
	jwtSecret := os.Getenv("JWT_SECRET")
	authSvc = auth.NewJWTService(jwtSecret, 60) // 有効期限60分
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginReq LoginRequest
	err := json.Unmarshal([]byte(req.Body), &loginReq)
	if err != nil || loginReq.Username == "" || loginReq.Password == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error":"invalid request"}`,
		}, nil
	}

	token, err := authSvc.Login(ctx, loginReq.Username, loginReq.Password)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       `{"error":"authentication failed"}`,
		}, nil
	}

	resp := LoginResponse{Token: token}
	respBody, _ := json.Marshal(resp)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
