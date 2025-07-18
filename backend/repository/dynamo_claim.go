package repository

import (
	"context"
	"errors"
	"fmt"

	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/yourorg/kaigo/backend/model"
)

type ClaimRepository interface {
	Create(ctx context.Context, claim *model.Claim) error
	GetByID(ctx context.Context, claimID string) (*model.Claim, error)
	Update(ctx context.Context, claim *model.Claim) error
	Delete(ctx context.Context, claimID string) error
	QueryByFacilityID(ctx context.Context, facilityID string) ([]*model.Claim, error)
	QueryByClaimDate(ctx context.Context, claimDate string) ([]*model.Claim, error)
}

type dynamoClaimRepo struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoClaimRepository(client *dynamodb.Client, tableName string) ClaimRepository {
	return &dynamoClaimRepo{
		client:    client,
		tableName: tableName,
	}
}

func (r *dynamoClaimRepo) Create(ctx context.Context, claim *model.Claim) error {
	now := time.Now().Format(time.RFC3339)
	claim.CreatedAt = now
	claim.UpdatedAt = now

	item := map[string]types.AttributeValue{
		"ClaimID":     &types.AttributeValueMemberS{Value: claim.ClaimID},
		"FacilityID":  &types.AttributeValueMemberS{Value: claim.FacilityID},
		"UserID":      &types.AttributeValueMemberS{Value: claim.UserID},
		"ClaimDate":   &types.AttributeValueMemberS{Value: claim.ClaimDate},
		"Amount":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", claim.Amount)},
		"Description": &types.AttributeValueMemberS{Value: claim.Description},
		"Status":      &types.AttributeValueMemberS{Value: claim.Status},
		"CreatedAt":   &types.AttributeValueMemberS{Value: claim.CreatedAt},
		"UpdatedAt":   &types.AttributeValueMemberS{Value: claim.UpdatedAt},
	}

	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(ClaimID)"), // 重複防止
	})
	if err != nil {
		return fmt.Errorf("failed to put claim item: %w", err)
	}
	return nil
}

func (r *dynamoClaimRepo) GetByID(ctx context.Context, claimID string) (*model.Claim, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ClaimID": &types.AttributeValueMemberS{Value: claimID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get claim: %w", err)
	}
	if out.Item == nil {
		return nil, errors.New("claim not found")
	}

	claim := &model.Claim{}
	if err := attributeValueToClaim(out.Item, claim); err != nil {
		return nil, err
	}
	return claim, nil
}

func (r *dynamoClaimRepo) Update(ctx context.Context, claim *model.Claim) error {
	claim.UpdatedAt = time.Now().Format(time.RFC3339)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ClaimID": &types.AttributeValueMemberS{Value: claim.ClaimID},
		},
		UpdateExpression: aws.String("SET FacilityID=:f, UserID=:u, ClaimDate=:cd, Amount=:a, Description=:d, Status=:s, UpdatedAt=:ua"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":f":  &types.AttributeValueMemberS{Value: claim.FacilityID},
			":u":  &types.AttributeValueMemberS{Value: claim.UserID},
			":cd": &types.AttributeValueMemberS{Value: claim.ClaimDate},
			":a":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", claim.Amount)},
			":d":  &types.AttributeValueMemberS{Value: claim.Description},
			":s":  &types.AttributeValueMemberS{Value: claim.Status},
			":ua": &types.AttributeValueMemberS{Value: claim.UpdatedAt},
		},
		ConditionExpression: aws.String("attribute_exists(ClaimID)"),
	}
	_, err := r.client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update claim: %w", err)
	}
	return nil
}

func (r *dynamoClaimRepo) Delete(ctx context.Context, claimID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ClaimID": &types.AttributeValueMemberS{Value: claimID},
		},
		ConditionExpression: aws.String("attribute_exists(ClaimID)"),
	})
	if err != nil {
		return fmt.Errorf("failed to delete claim: %w", err)
	}
	return nil
}

func (r *dynamoClaimRepo) QueryByFacilityID(ctx context.Context, facilityID string) ([]*model.Claim, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("FacilityIDIndex"),
		KeyConditionExpression: aws.String("FacilityID = :f"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":f": &types.AttributeValueMemberS{Value: facilityID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query claims by facility: %w", err)
	}
	return unmarshalClaims(out.Items)
}

func (r *dynamoClaimRepo) QueryByClaimDate(ctx context.Context, claimDate string) ([]*model.Claim, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("ClaimDateIndex"),
		KeyConditionExpression: aws.String("ClaimDate = :cd"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cd": &types.AttributeValueMemberS{Value: claimDate},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query claims by claim date: %w", err)
	}
	return unmarshalClaims(out.Items)
}

func attributeValueToClaim(item map[string]types.AttributeValue, claim *model.Claim) error {
	var err error
	if v, ok := item["ClaimID"].(*types.AttributeValueMemberS); ok {
		claim.ClaimID = v.Value
	} else {
		return errors.New("ClaimID missing or invalid")
	}
	if v, ok := item["FacilityID"].(*types.AttributeValueMemberS); ok {
		claim.FacilityID = v.Value
	}
	if v, ok := item["UserID"].(*types.AttributeValueMemberS); ok {
		claim.UserID = v.Value
	}
	if v, ok := item["ClaimDate"].(*types.AttributeValueMemberS); ok {
		claim.ClaimDate = v.Value
	}
	if v, ok := item["Amount"].(*types.AttributeValueMemberN); ok {
		claim.Amount, err = parseFloat(v.Value)
		if err != nil {
			return err
		}
	}
	if v, ok := item["Description"].(*types.AttributeValueMemberS); ok {
		claim.Description = v.Value
	}
	if v, ok := item["Status"].(*types.AttributeValueMemberS); ok {
		claim.Status = v.Value
	}
	if v, ok := item["CreatedAt"].(*types.AttributeValueMemberS); ok {
		claim.CreatedAt = v.Value
	}
	if v, ok := item["UpdatedAt"].(*types.AttributeValueMemberS); ok {
		claim.UpdatedAt = v.Value
	}
	return nil
}

func unmarshalClaims(items []map[string]types.AttributeValue) ([]*model.Claim, error) {
	var claims []*model.Claim
	for _, item := range items {
		c := &model.Claim{}
		if err := attributeValueToClaim(item, c); err != nil {
			return nil, err
		}
		claims = append(claims, c)
	}
	return claims, nil
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}
