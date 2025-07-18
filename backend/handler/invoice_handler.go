package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/backend/service"
)

type InvoiceHandler struct {
	Service *service.InvoiceService
}

func NewInvoiceHandler(service *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{Service: service}
}

func (h *InvoiceHandler) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return h.handleGet(ctx, req)
	case "POST":
		return h.handlePost(ctx, req)
	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed}, nil
	}
}

func (h *InvoiceHandler) handleGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	invoiceID := req.PathParameters["invoice_id"]
	invoice, err := h.Service.GetInvoice(ctx, invoiceID)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	if invoice == nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
	}
	body, _ := json.Marshal(invoice)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(body)}, nil
}

func (h *InvoiceHandler) handlePost(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var invoice model.Invoice
	if err := json.Unmarshal([]byte(req.Body), &invoice); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}
	if err := h.Service.CreateInvoice(ctx, &invoice); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: http.StatusCreated}, nil
}
