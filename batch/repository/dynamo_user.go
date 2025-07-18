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

type DynamoUserRepository struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func NewDynamoUserRepository(tableName string) *DynamoUserRepository {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)
	return &DynamoUserRepository{
		db:        db,
		tableName: tableName,
	}
}

// QueryUsersByLastLoginBefore GSIを使ってlast_login_atがcutoff以前のユーザーを取得
func (r *DynamoUserRepository) QueryUsersByLastLoginBefore(ctx context.Context, cutoff time.Time) ([]*model.User, error) {
	// 例: GSI名 "LastLoginIndex", PartitionKey: 固定値 or user_type, SortKey: LastLoginAt
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("LastLoginIndex"),
		KeyConditionExpression: aws.String("LastLoginAt <= :cutoff"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":cutoff": {
				S: aws.String(cutoff.Format(time.RFC3339)),
			},
		},
	}

	result, err := r.db.QueryWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("dynamo query failed: %w", err)
	}

	var users []*model.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}

	return users, nil
}

// Delete ユーザーの物理削除
func (r *DynamoUserRepository) Delete(ctx context.Context, userID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(userID),
			},
		},
	}
	_, err := r.db.DeleteItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete user from dynamo: %w", err)
	}
	return nil
}
