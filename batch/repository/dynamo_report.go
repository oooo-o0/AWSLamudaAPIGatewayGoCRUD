package repository

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/batch/config"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils/logger"
)

type ReportRepository interface {
	Put(ctx context.Context, report *model.Report) error
}

type reportRepository struct {
	tableName string
	client    *dynamodb.DynamoDB
}

func NewReportRepository() ReportRepository {
	return &reportRepository{
		tableName: config.GetEnv("REPORT_TABLE", "Report"),
		client:    config.GetDynamoDBClient(),
	}
}

func (r *reportRepository) Put(ctx context.Context, report *model.Report) error {
	log := logger.WithContext(ctx).Named("ReportRepository")

	item, err := dynamodbattribute.MarshalMap(report)
	if err != nil {
		log.Error("Marshal失敗", "error", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItemWithContext(ctx, input)
	if err != nil {
		log.Error("DynamoDB PutItem失敗", "error", err)
	}
	return err
}
