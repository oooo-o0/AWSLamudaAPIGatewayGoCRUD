package repository

import (
	"context"
	"kaigo-insurance-system/backend/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CarePlanRepository interface {
	Create(ctx context.Context, plan *model.CarePlan) error
	GetByID(ctx context.Context, id string) (*model.CarePlan, error)
	Update(ctx context.Context, plan *model.CarePlan) error
	Delete(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, userID string) ([]model.CarePlan, error)
}

type carePlanRepository struct {
	db        *dynamodb.Client
	tableName string
}

func NewCarePlanRepository(db *dynamodb.Client, tableName string) CarePlanRepository {
	return &carePlanRepository{db: db, tableName: tableName}
}

func (r *carePlanRepository) Create(ctx context.Context, plan *model.CarePlan) error {
	item, err := attributevalue.MarshalMap(plan)
	if err != nil {
		return err
	}
	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	return err
}

func (r *carePlanRepository) GetByID(ctx context.Context, id string) (*model.CarePlan, error) {
	resp, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"CarePlanID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil || resp.Item == nil {
		return nil, err
	}
	var plan model.CarePlan
	err = attributevalue.UnmarshalMap(resp.Item, &plan)
	return &plan, err
}

func (r *carePlanRepository) Update(ctx context.Context, plan *model.CarePlan) error {
	return r.Create(ctx, plan) // overwrite
}

func (r *carePlanRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"CarePlanID": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}

func (r *carePlanRepository) ListByUserID(ctx context.Context, userID string) ([]model.CarePlan, error) {
	indexName := "UserIDIndex"
	resp, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              &indexName,
		KeyConditionExpression: aws.String("UserID = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, err
	}
	var plans []model.CarePlan
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &plans)
	return plans, err
}
