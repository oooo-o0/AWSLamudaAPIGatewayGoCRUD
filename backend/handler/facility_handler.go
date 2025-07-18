package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type FacilityHandler struct {
	service service.FacilityService
}

func NewFacilityHandler(svc service.FacilityService) *FacilityHandler {
	return &FacilityHandler{service: svc}
}

func (h *FacilityHandler) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := req.Path
	method := req.HTTPMethod

	switch {
	case method == http.MethodPost && path == "/facilities":
		return h.createFacility(ctx, req)
	case method == http.MethodGet && strings.HasPrefix(path, "/facilities/"):
		id := strings.TrimPrefix(path, "/facilities/")
		return h.getFacility(ctx, id)
	case method == http.MethodPut && strings.HasPrefix(path, "/facilities/"):
		id := strings.TrimPrefix(path, "/facilities/")
		return h.updateFacility(ctx, id, req)
	case method == http.MethodDelete && strings.HasPrefix(path, "/facilities/"):
		id := strings.TrimPrefix(path, "/facilities/")
		return h.deleteFacility(ctx, id)
	case method == http.MethodGet && path == "/facilities": // ?name=xxx で検索
		name := req.QueryStringParameters["name"]
		return h.listFacilitiesByName(ctx, name)
	default:
		return clientError(http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *FacilityHandler) createFacility(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var f model.Facility
	if err := json.Unmarshal([]byte(req.Body), &f); err != nil {
		return clientError(http.StatusBadRequest, "invalid request body")
	}

	created, err := h.service.CreateFacility(ctx, &f)
	if err != nil {
		return serverError(err)
	}

	return responseJSON(http.StatusCreated, created)
}

func (h *FacilityHandler) getFacility(ctx context.Context, id string) (events.APIGatewayProxyResponse, error) {
	f, err := h.service.GetFacility(ctx, id)
	if err != nil {
		return clientError(http.StatusNotFound, "facility not found")
	}
	return responseJSON(http.StatusOK, f)
}

func (h *FacilityHandler) updateFacility(ctx context.Context, id string, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var f model.Facility
	if err := json.Unmarshal([]byte(req.Body), &f); err != nil {
		return clientError(http.StatusBadRequest, "invalid request body")
	}
	f.FacilityID = id

	updated, err := h.service.UpdateFacility(ctx, &f)
	if err != nil {
		return serverError(err)
	}
	return responseJSON(http.StatusOK, updated)
}

func (h *FacilityHandler) deleteFacility(ctx context.Context, id string) (events.APIGatewayProxyResponse, error) {
	err := h.service.DeleteFacility(ctx, id)
	if err != nil {
		return clientError(http.StatusNotFound, "facility not found")
	}
	return responseJSON(http.StatusNoContent, nil)
}

func (h *FacilityHandler) listFacilitiesByName(ctx context.Context, name string) (events.APIGatewayProxyResponse, error) {
	if name == "" {
		return clientError(http.StatusBadRequest, "query parameter 'name' is required")
	}
	list, err := h.service.FindFacilitiesByName(ctx, name)
	if err != nil {
		return serverError(err)
	}
	return responseJSON(http.StatusOK, list)
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
		Body:       `{"error":"` + err.Error() + `"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func responseJSON(status int, body interface{}) (events.APIGatewayProxyResponse, error) {
	respBody, err := json.Marshal(body)
	if err != nil {
		return serverError(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(respBody),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

// Lambdaのエントリポイント
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// DynamoDBクライアント初期化や環境変数読み込みは別途設定
	// ここでは例示のため省略

	// TODO: DynamoDBクライアント生成と環境変数TABLE_NAME取得
	// svc := dynamodb.New(session.New())
	// repo := repository.NewDynamoFacilityRepository(svc, tableName)
	// service := service.NewFacilityService(repo)
	// handler := NewFacilityHandler(service)
	// return handler.HandleRequest(ctx, req)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotImplemented,
		Body:       `{"error":"handler not implemented"}`,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
