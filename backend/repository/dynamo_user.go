package repository

import (
	"context"
	"errors"
	"kaigo-insurance-system/backend/model"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type UserRepository interface {
	Get(ctx context.Context, id string) (*model.User, error)
	Put(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	ListByRole(ctx context.Context, role string) ([]*model.User, error)
}

type userRepository struct {
	table string
	db    *dynamodb.DynamoDB
}

func NewUserRepository(sess *session.Session, tableName string) UserRepository {
	return &userRepository{
		table: tableName,
		db:    dynamodb.New(sess),
	}
}

func (r *userRepository) Get(ctx context.Context, id string) (*model.User, error) {
	out, err := r.db.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {S: aws.String(id)},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}

	var user model.User
	if err := dynamodbattribute.UnmarshalMap(out.Item, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Put(ctx context.Context, user *model.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = r.db.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      item,
	})
	return err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = r.db.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      item,
	})
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {S: aws.String(id)},
		},
	})
	return err
}

func (r *userRepository) ListByRole(ctx context.Context, role string) ([]*model.User, error) {
	out, err := r.db.QueryWithContext(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		IndexName:              aws.String("RoleIndex"),
		KeyConditionExpression: aws.String("Role = :r"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {S: aws.String(role)},
		},
	})
	if err != nil {
		return nil, err
	}

	var users []*model.User
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &users)
	return users, err
}
func (r *UserRepository) GetByID(ctx context.Context, userID string) (*model.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {S: aws.String(userID)},
		},
	}
	result, err := r.db.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, errors.New("user not found")
	}

	var user model.User
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
