package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
)

type DynamoNotificationRepository struct {
	client    *dynamodb.DynamoDB
	tableName string
}

func NewDynamoNotificationRepository(region, tableName string) (*DynamoNotificationRepository, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}
	return &DynamoNotificationRepository{
		client:    dynamodb.New(sess),
		tableName: tableName,
	}, nil
}

// QueryByStatus 指定ステータスの通知を上限limit件取得
func (r *DynamoNotificationRepository) QueryByStatus(ctx context.Context, status model.NotificationStatus, limit int) ([]model.Notification, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("StatusIndex"), // ステータスでのGSIを想定
		KeyConditionExpression: aws.String("#status = :status"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("Status"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(string(status)),
			},
		},
		Limit:            aws.Int64(int64(limit)),
		ScanIndexForward: aws.Bool(true), // 昇順（古い順）で処理
	}

	resp, err := r.client.QueryWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("dynamodb query failed: %w", err)
	}

	notifications := make([]model.Notification, 0, len(resp.Items))
	err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, &notifications)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal notifications: %w", err)
	}
	return notifications, nil
}

// UpdateStatus ステータスと更新日時のみを更新
func (r *DynamoNotificationRepository) UpdateStatus(ctx context.Context, notificationID string, status model.NotificationStatus, updatedAt time.Time) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"NotificationID": {S: aws.String(notificationID)},
		},
		UpdateExpression: aws.String("SET #status = :status, UpdatedAt = :updatedAt"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("Status"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status":    {S: aws.String(string(status))},
			":updatedAt": {S: aws.String(updatedAt.Format(time.RFC3339))},
		},
		ConditionExpression: aws.String("attribute_exists(NotificationID)"),
	}

	_, err := r.client.UpdateItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}
	return nil
}
