package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"kaigo-insurance-system/backend/auth"
	"kaigo-insurance-system/backend/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func AuthVerifyTokenHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	authHeader := req.Headers["Authorization"]
	if authHeader == "" {
		return middleware.NewErrorResponse(http.StatusUnauthorized, "Authorization header missing"), nil
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return middleware.NewErrorResponse(http.StatusUnauthorized, "Invalid authorization header"), nil
	}

	tokenStr := parts[1]
	claims, err := auth.VerifyToken(tokenStr)
	if err != nil {
		return middleware.NewErrorResponse(http.StatusUnauthorized, "Invalid or expired token"), nil
	}

	body, _ := json.Marshal(map[string]interface{}{
		"user_id": claims.UserID,
		"role":    claims.Role,
		"exp":     claims.ExpiresAt.Time.Unix(),
	})

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
		Headers:    middleware.DefaultHeaders(),
	}, nil
}

func main() {
	lambda.Start(AuthVerifyTokenHandler)
}
