package repository

import (
	"context"
	"time"

	"kaigo-insurance-system/backend/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	QueryByEventType(ctx context.Context, eventType string, from, to time.Time) ([]model.AuditLog, error)
	QueryByUserID(ctx context.Context, userID string, from, to time.Time) ([]model.AuditLog, error)
}

type auditLogRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewAuditLogRepository(db *dynamodb.Client, tableName string) AuditLogRepository {
	return &auditLogRepo{db: db, tableName: tableName}
}

func (r *auditLogRepo) Create(ctx context.Context, log *model.AuditLog) error {
	item, err := attributevalue.MarshalMap(log)
	if err != nil {
		return err
	}
	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(LogID)"),
	})
	return err
}

// EventTypeIndex GSIを使ったクエリ
func (r *auditLogRepo) QueryByEventType(ctx context.Context, eventType string, from, to time.Time) ([]model.AuditLog, error) {
	out, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("EventTypeIndex"),
		KeyConditionExpression: aws.String("EventType = :et AND Timestamp BETWEEN :from AND :to"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":et":   &types.AttributeValueMemberS{Value: eventType},
			":from": &types.AttributeValueMemberS{Value: from.Format(time.RFC3339)},
			":to":   &types.AttributeValueMemberS{Value: to.Format(time.RFC3339)},
		},
	})
	if err != nil {
		return nil, err
	}
	var logs []model.AuditLog
	err = attributevalue.UnmarshalListOfMaps(out.Items, &logs)
	return logs, err
}

// UserIDでのログ抽出例（必要に応じてGSIを追加してください）
func (r *auditLogRepo) QueryByUserID(ctx context.Context, userID string, from, to time.Time) ([]model.AuditLog, error) {
	// UserIDIndex GSIを利用する想定
	out, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("UserIDIndex"),
		KeyConditionExpression: aws.String("UserID = :uid AND Timestamp BETWEEN :from AND :to"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid":  &types.AttributeValueMemberS{Value: userID},
			":from": &types.AttributeValueMemberS{Value: from.Format(time.RFC3339)},
			":to":   &types.AttributeValueMemberS{Value: to.Format(time.RFC3339)},
		},
	})
	if err != nil {
		return nil, err
	}
	var logs []model.AuditLog
	err = attributevalue.UnmarshalListOfMaps(out.Items, &logs)
	return logs, err
}
