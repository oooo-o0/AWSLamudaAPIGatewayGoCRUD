package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"kaigo-insurance-system/backend/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type FacilityRepository interface {
	Create(ctx context.Context, f *model.Facility) error
	GetByID(ctx context.Context, id string) (*model.Facility, error)
	Update(ctx context.Context, f *model.Facility) error
	Delete(ctx context.Context, id string) error
	ListByName(ctx context.Context, name string) ([]*model.Facility, error)
}

type dynamoFacilityRepository struct {
	Client    dynamodbiface.DynamoDBAPI
	TableName string
}

func NewDynamoFacilityRepository(client dynamodbiface.DynamoDBAPI, tableName string) FacilityRepository {
	return &dynamoFacilityRepository{
		Client:    client,
		TableName: tableName,
	}
}

func (r *dynamoFacilityRepository) Create(ctx context.Context, f *model.Facility) error {
	now := time.Now().UTC()
	f.CreatedAt = now
	f.UpdatedAt = now

	av, err := dynamodbattribute.MarshalMap(f)
	if err != nil {
		return fmt.Errorf("marshal facility struct failed: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(r.TableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(FacilityID)"), // 重複禁止
	}

	_, err = r.Client.PutItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("PutItem failed: %w", err)
	}
	return nil
}

func (r *dynamoFacilityRepository) GetByID(ctx context.Context, id string) (*model.Facility, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"FacilityID": {S: aws.String(id)},
		},
	}

	result, err := r.Client.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("GetItem failed: %w", err)
	}
	if result.Item == nil {
		return nil, errors.New("facility not found")
	}

	var f model.Facility
	err = dynamodbattribute.UnmarshalMap(result.Item, &f)
	if err != nil {
		return nil, fmt.Errorf("unmarshal facility failed: %w", err)
	}
	return &f, nil
}

func (r *dynamoFacilityRepository) Update(ctx context.Context, f *model.Facility) error {
	// UpdatedAt更新のみ強制
	f.UpdatedAt = time.Now().UTC()

	// UpdateExpressionでName, Address, Phone, UpdatedAtを更新
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"FacilityID": {S: aws.String(f.FacilityID)},
		},
		UpdateExpression: aws.String("SET #N = :name, Address = :address, Phone = :phone, UpdatedAt = :updatedAt"),
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("Name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name":      {S: aws.String(f.Name)},
			":address":   {S: aws.String(f.Address)},
			":phone":     {S: aws.String(f.Phone)},
			":updatedAt": {S: aws.String(f.UpdatedAt.Format(time.RFC3339))},
		},
		ConditionExpression: aws.String("attribute_exists(FacilityID)"),
		ReturnValues:        aws.String("ALL_NEW"),
	}

	_, err := r.Client.UpdateItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("UpdateItem failed: %w", err)
	}
	return nil
}

func (r *dynamoFacilityRepository) Delete(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"FacilityID": {S: aws.String(id)},
		},
		ConditionExpression: aws.String("attribute_exists(FacilityID)"),
	}

	_, err := r.Client.DeleteItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("DeleteItem failed: %w", err)
	}
	return nil
}

func (r *dynamoFacilityRepository) ListByName(ctx context.Context, name string) ([]*model.Facility, error) {
	// GSI: NameIndexを使い名前で検索
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		IndexName:              aws.String("NameIndex"),
		KeyConditionExpression: aws.String("#N = :name"),
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("Name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {S: aws.String(name)},
		},
	}

	result, err := r.Client.QueryWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("Query failed: %w", err)
	}

	facilities := make([]*model.Facility, 0, len(result.Items))
	for _, item := range result.Items {
		var f model.Facility
		if err := dynamodbattribute.UnmarshalMap(item, &f); err != nil {
			return nil, fmt.Errorf("unmarshal facility failed: %w", err)
		}
		facilities = append(facilities, &f)
	}
	return facilities, nil
}
