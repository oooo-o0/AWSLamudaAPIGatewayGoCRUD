package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"kaigo-insurance-system/backend/auth"
	"kaigo-insurance-system/backend/middleware"
	"kaigo-insurance-system/backend/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var userService service.UserService

func init() {
	userService = service.NewUserService()
}

type RoleUpdateRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func UserRoleHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 認証情報の取得(仮にヘッダーのAuthorizationから)
	tokenStr := req.Headers["Authorization"]
	claims, err := auth.VerifyToken(strings.TrimPrefix(tokenStr, "Bearer "))
	if err != nil {
		return middleware.NewErrorResponse(http.StatusUnauthorized, "Unauthorized"), nil
	}

	// admin権限のみ許可
	if !auth.HasRole(claims.Role, auth.RoleAdmin) {
		return middleware.NewErrorResponse(http.StatusForbidden, "Forbidden"), nil
	}

	var reqBody RoleUpdateRequest
	if err := json.Unmarshal([]byte(req.Body), &reqBody); err != nil {
		return middleware.NewErrorResponse(http.StatusBadRequest, "Invalid request body"), nil
	}

	if reqBody.UserID == "" || reqBody.Role == "" {
		return middleware.NewErrorResponse(http.StatusBadRequest, "Missing user_id or role"), nil
	}

	err = userService.UpdateUserRole(ctx, reqBody.UserID, reqBody.Role)
	if err != nil {
		return middleware.NewErrorResponse(http.StatusInternalServerError, "Failed to update role"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message":"Role updated successfully"}`,
		Headers:    middleware.DefaultHeaders(),
	}, nil
}

func main() {
	lambda.Start(UserRoleHandler)
}
