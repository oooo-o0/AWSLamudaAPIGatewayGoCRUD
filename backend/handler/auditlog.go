package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/service"

	"github.com/aws/aws-lambda-go/events"
)

type AuditLogHandler struct {
	Service service.AuditLogService
}

func NewAuditLogHandler(s service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{Service: s}
}

func (h *AuditLogHandler) Create(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input model.AuditLog
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return clientError(http.StatusBadRequest, "invalid request body")
	}

	// Timestampを現在時刻に補正する例
	input.Timestamp = time.Now().UTC()

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

func (h *AuditLogHandler) QueryByEventType(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	eventType := req.QueryStringParameters["eventType"]
	fromStr := req.QueryStringParameters["from"]
	toStr := req.QueryStringParameters["to"]

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return clientError(http.StatusBadRequest, "invalid from datetime")
	}
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return clientError(http.StatusBadRequest, "invalid to datetime")
	}

	logs, err := h.Service.QueryByEventType(ctx, eventType, from, to)
	if err != nil {
		return serverError(err)
	}

	body, _ := json.Marshal(logs)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}
