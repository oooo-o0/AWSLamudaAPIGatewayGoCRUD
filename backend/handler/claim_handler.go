package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/yourorg/kaigo/backend/model"
	"github.com/yourorg/kaigo/backend/repository"
	"github.com/yourorg/kaigo/backend/service"
)

var claimService service.ClaimService

func init() {
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		panic("TABLE_NAME environment variable is required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	repo := repository.NewDynamoClaimRepository(dynamoClient, tableName)
	claimService = service.NewClaimService(repo)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := strings.TrimSuffix(req.Path, "/")
	method := req.HTTPMethod

	switch {
	case method == http.MethodPost && path == "/claims":
		return handleCreateClaim(ctx, req)
	case method == http.MethodGet && strings.HasPrefix(path, "/claims/"):
		claimID := strings.TrimPrefix(path, "/claims/")
		return handleGetClaim(ctx, claimID)
	case method == http.MethodPut && strings.HasPrefix(path, "/claims/"):
		claimID := strings.TrimPrefix(path, "/claims/")
		return handleUpdateClaim(ctx, claimID, req)
	case method == http.MethodDelete && strings.HasPrefix(path, "/claims/"):
		claimID := strings.TrimPrefix(path, "/claims/")
		return handleDeleteClaim(ctx, claimID)
	case method == http.MethodGet && path == "/claims":
		facilityID := req.QueryStringParameters["facility_id"]
		claimDate := req.QueryStringParameters["claim_date"]
		if facilityID != "" {
			return handleQueryByFacilityID(ctx, facilityID)
		} else if claimDate != "" {
			return handleQueryByClaimDate(ctx, claimDate)
		} else {
			return clientError(http.StatusBadRequest, "facility_id or claim_date query parameter required")
		}
	default:
		return clientError(http.StatusNotFound, "not found")
	}
}

func handleCreateClaim(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var claim model.Claim
	if err := json.Unmarshal([]byte(req.Body), &claim); err != nil {
		return clientError(http.StatusBadRequest, "invalid request body")
	}
	err := claimService.CreateClaim(ctx, &claim)
	if err != nil {
		return serverError(err)
	}
	return respondJSON(http.StatusCreated, claim)
}

func handleGetClaim(ctx context.Context, claimID string) (events.APIGatewayProxyResponse, error) {
	claim, err := claimService.GetClaim(ctx, claimID)
	if err != nil {
		return clientError(http.StatusNotFound, "claim not found")
	}
	return respondJSON(http.StatusOK, claim)
}

func handleUpdateClaim(ctx context.Context, claimID string, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var claim model.Claim
	if err := json.Unmarshal([]byte(req.Body), &claim); err != nil {
		return clientError(http.StatusBadRequest, "invalid request body")
	}
	if claimID != claim.ClaimID {
		return clientError(http.StatusBadRequest, "path ClaimID and body ClaimID mismatch")
	}
	err := claimService.UpdateClaim(ctx, &claim)
	if err != nil {
		return serverError(err)
	}
	return respondJSON(http.StatusOK, claim)
}

func handleDeleteClaim(ctx context.Context, claimID string) (events.APIGatewayProxyResponse, error) {
	err := claimService.DeleteClaim(ctx, claimID)
	if err != nil {
		return serverError(err)
	}
	return respondJSON(http.StatusNoContent, nil)
}

func handleQueryByFacilityID(ctx context.Context, facilityID string) (events.APIGatewayProxyResponse, error) {
	claims, err := claimService.QueryClaimsByFacility(ctx, facilityID)
	if err != nil {
		return serverError(err)
	}
	return respondJSON(http.StatusOK, claims)
}

func handleQueryByClaimDate(ctx context.Context, claimDate string) (events.APIGatewayProxyResponse, error) {
	claims, err := claimService.QueryClaimsByDate(ctx, claimDate)
	if err != nil {
		return serverError(err)
	}
	return respondJSON(http.StatusOK, claims)
}

func clientError(status int, message string) (events.APIGatewayProxyResponse, error) {
	return respondJSON(status, map[string]string{"error": message})
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Internal server error:", err)
	return respondJSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

func respondJSON(status int, body interface{}) (events.APIGatewayProxyResponse, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error":"failed to marshal response"}`,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode:      status,
		Body:            string(jsonBytes),
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
