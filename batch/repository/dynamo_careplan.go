package repository

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/batch/config"
	"kaigo-insurance-system/batch/utils/dynamo"
)

type CarePlanRepository interface {
	FindCarePlansBeforeDate(ctx context.Context, date string) ([]model.CarePlan, error)
	UpdateCarePlan(ctx context.Context, plan *model.CarePlan) error
}

type carePlanRepositoryImpl struct {
	tableName string
	client    *dynamodb.DynamoDB
}

func NewCarePlanRepository() CarePlanRepository {
	return &carePlanRepositoryImpl{
		tableName: config.GetCarePlanTableName(),
		client:    dynamo.GetClient(),
	}
}

func (r *carePlanRepositoryImpl) FindCarePlansBeforeDate(ctx context.Context, date string) ([]model.CarePlan, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(r.tableName),
		FilterExpression: aws.String("EndDate < :now"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":now": {S: aws.String(date)},
		},
	}

	result, err := r.client.ScanWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	var careplans []model.CarePlan
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &careplans)
	return careplans, err
}

func (r *carePlanRepositoryImpl) UpdateCarePlan(ctx context.Context, plan *model.CarePlan) error {
	av, err := dynamodbattribute.MarshalMap(plan)
	if err != nil {
		return err
	}

	_, err = r.client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	return err
}
