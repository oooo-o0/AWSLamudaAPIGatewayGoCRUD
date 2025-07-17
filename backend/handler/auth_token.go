package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/service"

	"github.com/aws/aws-lambda-go/events"
)

type AuthTokenHandler struct {
	Service service.AuthTokenService
}

func NewAuthTokenHandler(s service.AuthTokenService) *AuthTokenHandler {
	return &AuthTokenHandler{Service: s}
}

func (h *AuthTokenHandler) Create(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input model.AuthToken
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return clientError(http.StatusBadRequest, "invalid request body")
	}

	err := h.Service.Create(ctx, &input)
	if err != nil {
		return serverError(err)
	}

	body, _ := json.Marshal(input)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(body),
	}, nil
}

func (h *AuthTokenHandler) Get(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tokenID := req.PathParameters["id"]
	token, err := h.Service.GetByID(ctx, tokenID)
	if err != nil {
		return clientError(http.StatusNotFound, "token not found")
	}
	body, _ := json.Marshal(token)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func (h *AuthTokenHandler) Revoke(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tokenID := req.PathParameters["id"]
	err := h.Service.Revoke(ctx, tokenID)
	if err != nil {
		return serverError(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message":"revoked"}`,
	}, nil
}
