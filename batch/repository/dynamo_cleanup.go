package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils"
)

type DynamoCleanupRepository interface {
	FindExpiredReports(ctx context.Context, reportType model.ReportType, before time.Time) ([]model.Report, error)
	DeleteReport(ctx context.Context, reportID string) error
	PutAuditLog(ctx context.Context, log model.AuditLog) error
	PutNotification(ctx context.Context, n model.Notification) error
}

type dynamoCleanupRepository struct {
	db *dynamodb.DynamoDB
}

func NewDynamoCleanupRepository(cfg *utils.Config) DynamoCleanupRepository {
	return &dynamoCleanupRepository{
		db: utils.NewDynamoClient(cfg),
	}
}

func (r *dynamoCleanupRepository) FindExpiredReports(ctx context.Context, reportType model.ReportType, before time.Time) ([]model.Report, error) {
	input := &dynamodb.ScanInput{
		TableName:        utils.StrPtr("Report"),
		FilterExpression: utils.StrPtr("ReportType = :type AND GeneratedDate < :before"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":type":   {S: utils.StrPtr(string(reportType))},
			":before": {S: utils.StrPtr(before.Format(time.RFC3339))},
		},
	}

	result, err := r.db.ScanWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("Scan失敗: %w", err)
	}

	var reports []model.Report
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &reports); err != nil {
		return nil, fmt.Errorf("Unmarshal失敗: %w", err)
	}
	return reports, nil
}

func (r *dynamoCleanupRepository) DeleteReport(ctx context.Context, reportID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: utils.StrPtr("Report"),
		Key: map[string]*dynamodb.AttributeValue{
			"ReportID": {S: utils.StrPtr(reportID)},
		},
	}
	_, err := r.db.DeleteItemWithContext(ctx, input)
	return err
}

func (r *dynamoCleanupRepository) PutAuditLog(ctx context.Context, log model.AuditLog) error {
	item, _ := dynamodbattribute.MarshalMap(log)
	input := &dynamodb.PutItemInput{
		TableName: utils.StrPtr("AuditLog"),
		Item:      item,
	}
	_, err := r.db.PutItemWithContext(ctx, input)
	return err
}

func (r *dynamoCleanupRepository) PutNotification(ctx context.Context, n model.Notification) error {
	item, _ := dynamodbattribute.MarshalMap(n)
	input := &dynamodb.PutItemInput{
		TableName: utils.StrPtr("Notification"),
		Item:      item,
	}
	_, err := r.db.PutItemWithContext(ctx, input)
	return err
}
