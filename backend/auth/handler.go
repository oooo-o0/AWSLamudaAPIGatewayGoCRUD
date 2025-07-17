package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func AuthLoginHandler(ctx context.Context, req events.APIGatewayProxyRequest, jwtMgr *JWTManager) (events.APIGatewayProxyResponse, error) {
	var loginReq LoginRequest
	if err := json.Unmarshal([]byte(req.Body), &loginReq); err != nil {
		return clientError(http.StatusBadRequest, "invalid request body")
	}

	// ここで実際はDB照合を行うべき
	if loginReq.Username != "admin" || loginReq.Password != "password" {
		return clientError(http.StatusUnauthorized, "invalid credentials")
	}

	token, err := jwtMgr.GenerateToken("user-id-123", []string{"admin"}, 1*time.Hour)
	if err != nil {
		return serverError(err)
	}

	respBody, _ := json.Marshal(LoginResponse{Token: token})
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBody),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func AuthVerifyTokenHandler(ctx context.Context, req events.APIGatewayProxyRequest, jwtMgr *JWTManager) (events.APIGatewayProxyResponse, error) {
	authHeader := req.Headers["Authorization"]
	if authHeader == "" {
		return clientError(http.StatusUnauthorized, "missing Authorization header")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return clientError(http.StatusUnauthorized, "invalid Authorization header format")
	}

	claims, err := jwtMgr.VerifyToken(tokenStr)
	if err != nil {
		return clientError(http.StatusUnauthorized, "invalid token")
	}

	respBody, _ := json.Marshal(claims)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBody),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func clientError(status int, msg string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       `{"error":"` + msg + `"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       `{"error":"internal server error"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, err
}
