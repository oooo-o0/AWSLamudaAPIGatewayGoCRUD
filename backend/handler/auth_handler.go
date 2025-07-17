package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/service"

	"github.com/aws/aws-lambda-go/events"
)

type AuthHandler struct {
	Service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{Service: s}
}

func (h *AuthHandler) HandleLogin(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body model.LoginRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return clientError(http.StatusBadRequest, "invalid body")
	}

	authResp, err := h.Service.Authenticate(ctx, body.Email, body.Password)
	if err != nil {
		return clientError(http.StatusUnauthorized, err.Error())
	}

	jsonBody, _ := json.Marshal(authResp)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonBody),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func (h *AuthHandler) HandleSignup(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body model.SignupRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return clientError(http.StatusBadRequest, "invalid body")
	}

	if err := h.Service.Signup(ctx, &body); err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       `{"message":"user registered"}`,
	}, nil
}

func (h *AuthHandler) HandleTokenVerify(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	token := req.Headers["Authorization"]
	if token == "" {
		return clientError(http.StatusUnauthorized, "missing token")
	}

	_, err := h.Service.ValidateToken(ctx, token)
	if err != nil {
		return clientError(http.StatusUnauthorized, "invalid token")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message":"valid token"}`,
	}, nil
}
