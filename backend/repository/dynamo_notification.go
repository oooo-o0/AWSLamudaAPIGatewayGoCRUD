package repository

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"kaigo-insurance-system/backend/model"
)

type NotificationRepository interface {
	Put(ctx context.Context, notification *model.Notification) error
	GetByID(ctx context.Context, notificationID string) (*model.Notification, error)
	UpdateStatus(ctx context.Context, notificationID string, status model.NotificationStatus) error
	Delete(ctx context.Context, notificationID string) error
	QueryByUserID(ctx context.Context, userID string) ([]*model.Notification, error)
	QueryByStatus(ctx context.Context, status model.NotificationStatus) ([]*model.Notification, error)
}

type dynamoNotificationRepository struct {
	db        *dynamodb.Client
	tableName string
}

func NewDynamoNotificationRepository(db *dynamodb.Client, tableName string) NotificationRepository {
	return &dynamoNotificationRepository{
		db:        db,
		tableName: tableName,
	}
}

func (r *dynamoNotificationRepository) Put(ctx context.Context, notification *model.Notification) error {
	notification.UpdatedAt = time.Now().UTC()
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = notification.UpdatedAt
	}

	item, err := attributevalue.MarshalMap(notification)
	if err != nil {
		return err
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	return err
}

func (r *dynamoNotificationRepository) GetByID(ctx context.Context, notificationID string) (*model.Notification, error) {
	out, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"NotificationID": &types.AttributeValueMemberS{Value: notificationID},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, errors.New("notification not found")
	}

	var notification model.Notification
	if err := attributevalue.UnmarshalMap(out.Item, &notification); err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *dynamoNotificationRepository) UpdateStatus(ctx context.Context, notificationID string, status model.NotificationStatus) error {
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := r.db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"NotificationID": &types.AttributeValueMemberS{Value: notificationID},
		},
		UpdateExpression: aws.String("SET #S = :s, UpdatedAt = :u"),
		ExpressionAttributeNames: map[string]string{
			"#S": "Status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":s": &types.AttributeValueMemberS{Value: string(status)},
			":u": &types.AttributeValueMemberS{Value: now},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})
	return err
}

func (r *dynamoNotificationRepository) Delete(ctx context.Context, notificationID string) error {
	_, err := r.db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"NotificationID": &types.AttributeValueMemberS{Value: notificationID},
		},
	})
	return err
}

func (r *dynamoNotificationRepository) QueryByUserID(ctx context.Context, userID string) ([]*model.Notification, error) {
	var notifications []*model.Notification

	out, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              aws.String("UserIDIndex"),
		KeyConditionExpression: aws.String("UserID = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
		ScanIndexForward: aws.Bool(false), // 最新順に
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *dynamoNotificationRepository) QueryByStatus(ctx context.Context, status model.NotificationStatus) ([]*model.Notification, error) {
	var notifications []*model.Notification

	out, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              aws.String("StatusIndex"),
		KeyConditionExpression: aws.String("Status = :status"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: string(status)},
		},
		ScanIndexForward: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
