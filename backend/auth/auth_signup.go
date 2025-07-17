package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"kaigo-insurance-system/backend/common"
	"kaigo-insurance-system/backend/middleware"
	"kaigo-insurance-system/backend/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type SignupRequest struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}

type SignupResponse struct {
	Message string `json:"message"`
}

var userService service.UserService

func init() {
	userService = service.NewUserService()
}

func SignupHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var signupReq SignupRequest
	if err := json.Unmarshal([]byte(req.Body), &signupReq); err != nil {
		return middleware.NewErrorResponse(http.StatusBadRequest, "Invalid request body"), nil
	}

	// 入力チェック
	if signupReq.UserID == "" || signupReq.Email == "" || signupReq.Password == "" || signupReq.Name == "" {
		return middleware.NewErrorResponse(http.StatusBadRequest, "Missing required fields"), nil
	}

	// ユーザー登録
	err := userService.Signup(ctx, signupReq.UserID, signupReq.Email, signupReq.Password, signupReq.Role, signupReq.Name)
	if err != nil {
		if common.IsConflictError(err) {
			return middleware.NewErrorResponse(http.StatusConflict, "User already exists"), nil
		}
		return middleware.NewErrorResponse(http.StatusInternalServerError, "Internal server error"), nil
	}

	resp := SignupResponse{
		Message: "User registered successfully",
	}
	body, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(body),
		Headers:    middleware.DefaultHeaders(),
	}, nil
}

func main() {
	lambda.Start(SignupHandler)
}
