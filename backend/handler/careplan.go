package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"kaigo-insurance-system/backend/config"
	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/repository"
	"kaigo-insurance-system/backend/service"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	carePlanService service.CarePlanService
)

func init() {
	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	dynamoClient := dynamodb.NewFromConfig(awsCfg)
	repo := repository.NewCarePlanRepository(dynamoClient, os.Getenv("TABLE_NAME"))
	carePlanService = service.NewCarePlanService(repo)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return createCarePlan(ctx, req)
	case "GET":
		if id := req.PathParameters["care_plan_id"]; id != "" {
			return getCarePlan(ctx, req, id)
		}
		userID := req.QueryStringParameters["user_id"]
		return listCarePlansByUser(ctx, userID)
	case "PUT":
		return updateCarePlan(ctx, req)
	case "DELETE":
		id := req.PathParameters["care_plan_id"]
		return deleteCarePlan(ctx, id)
	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed}, nil
	}
}

func createCarePlan(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var plan model.CarePlan
	if err := json.Unmarshal([]byte(req.Body), &plan); err != nil {
		return responseError(http.StatusBadRequest, "Invalid JSON")
	}
	err := carePlanService.Create(ctx, &plan)
	if err != nil {
		return responseError(http.StatusInternalServerError, err.Error())
	}
	return responseJSON(http.StatusCreated, plan)
}

func getCarePlan(ctx context.Context, req events.APIGatewayProxyRequest, id string) (events.APIGatewayProxyResponse, error) {
	plan, err := carePlanService.Get(ctx, id)
	if err != nil {
		return responseError(http.StatusNotFound, "Not found")
	}
	return responseJSON(http.StatusOK, plan)
}

func updateCarePlan(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var plan model.CarePlan
	if err := json.Unmarshal([]byte(req.Body), &plan); err != nil {
		return responseError(http.StatusBadRequest, "Invalid JSON")
	}
	err := carePlanService.Update(ctx, &plan)
	if err != nil {
		return responseError(http.StatusInternalServerError, err.Error())
	}
	return responseJSON(http.StatusOK, plan)
}

func deleteCarePlan(ctx context.Context, id string) (events.APIGatewayProxyResponse, error) {
	err := carePlanService.Delete(ctx, id)
	if err != nil {
		return responseError(http.StatusInternalServerError, err.Error())
	}
	return responseJSON(http.StatusNoContent, nil)
}

func listCarePlansByUser(ctx context.Context, userID string) (events.APIGatewayProxyResponse, error) {
	plans, err := carePlanService.ListByUser(ctx, userID)
	if err != nil {
		return responseError(http.StatusInternalServerError, err.Error())
	}
	return responseJSON(http.StatusOK, plans)
}
