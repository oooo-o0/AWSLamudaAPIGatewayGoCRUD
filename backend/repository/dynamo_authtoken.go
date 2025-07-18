package repository

import (
	"context"
	"errors"
	"time"

	"kaigo-insurance-system/backend/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AuthTokenRepository interface {
	Create(ctx context.Context, token *model.AuthToken) error
	GetByID(ctx context.Context, tokenID string) (*model.AuthToken, error)
	GetByUserID(ctx context.Context, userID string) ([]model.AuthToken, error)
	Delete(ctx context.Context, tokenID string) error
	Revoke(ctx context.Context, tokenID string) error
	CleanupExpired(ctx context.Context, before time.Time) error
}

type authTokenRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewAuthTokenRepository(db *dynamodb.Client, tableName string) AuthTokenRepository {
	return &authTokenRepo{db: db, tableName: tableName}
}

func (r *authTokenRepo) Create(ctx context.Context, token *model.AuthToken) error {
	now := time.Now()
	token.CreatedAt = now
	token.Revoked = false

	item, err := attributevalue.MarshalMap(token)
	if err != nil {
		return err
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(TokenID)"),
	})
	return err
}

func (r *authTokenRepo) GetByID(ctx context.Context, tokenID string) (*model.AuthToken, error) {
	out, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"TokenID": &types.AttributeValueMemberS{Value: tokenID},
		},
	})
	if err != nil || out.Item == nil {
		return nil, errors.New("token not found")
	}
	var token model.AuthToken
	err = attributevalue.UnmarshalMap(out.Item, &token)
	return &token, err
}

func (r *authTokenRepo) GetByUserID(ctx context.Context, userID string) ([]model.AuthToken, error) {
	out, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("UserIDIndex"),
		KeyConditionExpression: aws.String("UserID = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, err
	}

	var tokens []model.AuthToken
	err = attributevalue.UnmarshalListOfMaps(out.Items, &tokens)
	return tokens, err
}

func (r *authTokenRepo) Delete(ctx context.Context, tokenID string) error {
	_, err := r.db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"TokenID": &types.AttributeValueMemberS{Value: tokenID},
		},
	})
	return err
}

func (r *authTokenRepo) Revoke(ctx context.Context, tokenID string) error {
	// トークン失効フラグを立てる更新
	_, err := r.db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"TokenID": &types.AttributeValueMemberS{Value: tokenID},
		},
		UpdateExpression:          aws.String("SET Revoked = :rev"),
		ExpressionAttributeValues: map[string]types.AttributeValue{":rev": &types.AttributeValueMemberBOOL{Value: true}},
	})
	return err
}

func (r *authTokenRepo) CleanupExpired(ctx context.Context, before time.Time) error {
	// 有効期限切れのトークンを削除（バッチなどで使用）
	// DynamoDB単体での範囲削除は難しいため、通常はScan+Deleteのループ処理をする想定

	var lastEvaluatedKey map[string]types.AttributeValue

	for {
		input := &dynamodb.ScanInput{
			TableName:        aws.String(r.tableName),
			FilterExpression: aws.String("Expiry <= :now"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":now": &types.AttributeValueMemberS{Value: before.Format(time.RFC3339)},
			},
			ExclusiveStartKey: lastEvaluatedKey,
		}

		out, err := r.db.Scan(ctx, input)
		if err != nil {
			return err
		}

		for _, item := range out.Items {
			var token model.AuthToken
			if err := attributevalue.UnmarshalMap(item, &token); err == nil {
				_ = r.Delete(ctx, token.TokenID) // 削除失敗はログに任せて継続
			}
		}

		if out.LastEvaluatedKey == nil {
			break
		}
		lastEvaluatedKey = out.LastEvaluatedKey
	}
	return nil
}
