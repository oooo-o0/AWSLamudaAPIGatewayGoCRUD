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

type ReportHandler struct {
	Service service.ReportService
}

func NewReportHandler(s service.ReportService) *ReportHandler {
	return &ReportHandler{Service: s}
}

func (h *ReportHandler) Create(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input model.Report
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

func (h *ReportHandler) Get(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	report, err := h.Service.Get(ctx, id)
	if err != nil {
		return clientError(http.StatusNotFound, "report not found")
	}

	body, _ := json.Marshal(report)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func (h *ReportHandler) QueryByTypeAndDate(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reportType := req.QueryStringParameters["type"]
	fromStr := req.QueryStringParameters["from"]
	toStr := req.QueryStringParameters["to"]

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return clientError(http.StatusBadRequest, "invalid from date")
	}
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return clientError(http.StatusBadRequest, "invalid to date")
	}

	reports, err := h.Service.QueryByTypeAndDate(ctx, model.ReportType(reportType), from, to)
	if err != nil {
		return serverError(err)
	}

	body, _ := json.Marshal(reports)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func (h *ReportHandler) Update(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input model.Report
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return clientError(http.StatusBadRequest, "invalid request body")
	}

	err := h.Service.Update(ctx, &input)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message":"updated"}`,
	}, nil
}

func (h *ReportHandler) Delete(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	err := h.Service.Delete(ctx, id)
	if err != nil {
		return clientError(http.StatusNotFound, "report not found")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message":"deleted"}`,
	}, nil
}
