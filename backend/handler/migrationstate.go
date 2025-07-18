package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/service"

	"github.com/aws/aws-lambda-go/events"
)

type MigrationStateHandler struct {
	Service service.MigrationStateService
}

func NewMigrationStateHandler(s service.MigrationStateService) *MigrationStateHandler {
	return &MigrationStateHandler{Service: s}
}

// マイグレーション履歴登録API
func (h *MigrationStateHandler) Create(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input model.MigrationState
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

// 最新マイグレーション取得API
func (h *MigrationStateHandler) GetLatest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	state, err := h.Service.GetLatest(ctx)
	if err != nil {
		return serverError(err)
	}
	if state == nil {
		return clientError(http.StatusNotFound, "no migration state found")
	}
	body, _ := json.Marshal(state)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}
