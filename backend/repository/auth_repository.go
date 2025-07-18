package repository

import (
	"context"
	"errors"

	"kaigo-insurance-system/backend/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type AuthRepository interface {
	SaveToken(ctx context.Context, token *model.AuthToken) error
	GetTokenByID(ctx context.Context, tokenID string) (*model.AuthToken, error)
	RevokeToken(ctx context.Context, tokenID string) error
}

type dynamoAuthRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewAuthRepository(client *dynamodb.Client, tableName string) AuthRepository {
	return &dynamoAuthRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *dynamoAuthRepository) SaveToken(ctx context.Context, token *model.AuthToken) error {
	item, err := attributevalue.MarshalMap(token)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *dynamoAuthRepository) GetTokenByID(ctx context.Context, tokenID string) (*model.AuthToken, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]dynamodb.AttributeValue{
			"TokenID": &dynamodb.AttributeValueMemberS{Value: tokenID},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, errors.New("token not found")
	}

	var token model.AuthToken
	if err := attributevalue.UnmarshalMap(out.Item, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *dynamoAuthRepository) RevokeToken(ctx context.Context, tokenID string) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]dynamodb.AttributeValue{
			"TokenID": &dynamodb.AttributeValueMemberS{Value: tokenID},
		},
		UpdateExpression:          aws.String("SET IsRevoked = :r"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{":r": &dynamodb.AttributeValueMemberBOOL{Value: true}},
	})
	return err
}
