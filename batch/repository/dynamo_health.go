package repository

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/batch/config"
)

type AuditLogRepository interface {
	PutAuditLog(ctx context.Context, log model.AuditLog) error
}

type dynamoHealthRepo struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func NewDynamoHealthRepository(cfg config.Config) AuditLogRepository {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWSRegion),
	}))
	return &dynamoHealthRepo{
		db:        dynamodb.New(sess),
		tableName: cfg.AuditLogTableName,
	}
}

func (r *dynamoHealthRepo) PutAuditLog(ctx context.Context, log model.AuditLog) error {
	item, err := dynamodbattribute.MarshalMap(log)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.db.PutItemWithContext(ctx, input)
	return err
}
