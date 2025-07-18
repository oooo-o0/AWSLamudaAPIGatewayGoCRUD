package repository

import (
	"context"
	"time"

	"batch/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
)

type FacilityRepository interface {
	FetchClaimsForDate(ctx context.Context, date string) ([]model.Claim, error)
	StoreFacilityUsageSummary(ctx context.Context, date string, summaries map[string]model.FacilityUsageSummary) error
}

type facilityRepository struct {
	dynamo  *dynamodb.DynamoDB
	tblName string
}

func NewFacilityRepository() FacilityRepository {
	return &facilityRepository{
		dynamo:  config.DynamoClient(),
		tblName: "Claim",
	}
}

func (r *facilityRepository) FetchClaimsForDate(ctx context.Context, date string) ([]model.Claim, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tblName),
		IndexName:              aws.String("ClaimDateIndex"),
		KeyConditionExpression: aws.String("claim_date = :v1"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {S: aws.String(date)},
		},
	}

	out, err := r.dynamo.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	var claims []model.Claim
	if err := dynamodbattribute.UnmarshalListOfMaps(out.Items, &claims); err != nil {
		return nil, err
	}
	return claims, nil
}

func (r *facilityRepository) StoreFacilityUsageSummary(ctx context.Context, date string, summaries map[string]model.FacilityUsageSummary) error {
	for _, summary := range summaries {
		item := map[string]interface{}{
			"pk":           "FACILITY_USAGE#" + summary.FacilityID,
			"sk":           date,
			"total_claims": summary.TotalClaims,
			"total_amount": summary.TotalAmount,
			"created_at":   time.Now().Format(time.RFC3339),
		}

		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			return err
		}

		_, err = r.dynamo.PutItemWithContext(ctx, &dynamodb.PutItemInput{
			TableName: aws.String("FacilityUsageSummary"),
			Item:      av,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
