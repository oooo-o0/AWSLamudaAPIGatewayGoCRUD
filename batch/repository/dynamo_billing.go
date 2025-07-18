package repository

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"kaigo-insurance-system/backend/model"
)

// BillingRepositoryInterface は請求バッチで使うリポジトリインターフェース
type BillingRepositoryInterface interface {
	FetchUsersForBilling(ctx context.Context) ([]model.User, error)
	FetchApprovedClaimsByUser(ctx context.Context, userID string) ([]model.Claim, error)
	SaveInvoice(ctx context.Context, invoice *model.Invoice) error
}

type DynamoBillingRepository struct {
	client       *dynamodb.DynamoDB
	usersTable   string
	claimsTable  string
	invoiceTable string
}

func NewDynamoBillingRepository(usersTable string) (*DynamoBillingRepository, error) {
	if usersTable == "" {
		return nil, errors.New("users table name required")
	}
	// テーブル名は環境変数やconfigで管理。ここは簡略化しusersTableは必須パラメータ。
	sess := session.Must(session.NewSession())
	client := dynamodb.New(sess)
	return &DynamoBillingRepository{
		client:       client,
		usersTable:   usersTable,
		claimsTable:  "Claim",   // Claimテーブル名は固定またはconfig化可能
		invoiceTable: "Invoice", // Invoiceテーブル名も同様
	}, nil
}

func (r *DynamoBillingRepository) FetchUsersForBilling(ctx context.Context) ([]model.User, error) {
	// Role="user"のユーザー全件取得例。フィルターで実装可能。
	input := &dynamodb.ScanInput{
		TableName:        aws.String(r.usersTable),
		FilterExpression: aws.String("#role = :role"),
		ExpressionAttributeNames: map[string]*string{
			"#role": aws.String("Role"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":role": {S: aws.String("user")},
		},
	}

	result, err := r.client.ScanWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0, len(result.Items))
	for _, item := range result.Items {
		var u model.User
		if err := dynamodbattribute.UnmarshalMap(item, &u); err != nil {
			continue // アンマーシャル失敗はスキップ
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *DynamoBillingRepository) FetchApprovedClaimsByUser(ctx context.Context, userID string) ([]model.Claim, error) {
	// GSI: UserIDIndex でクエリ。Status="Approved" のみ取得
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.claimsTable),
		IndexName: aws.String("UserIDIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"UserID": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{S: aws.String(userID)},
				},
			},
		},
		FilterExpression: aws.String("#status = :status"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("Status"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {S: aws.String("Approved")},
		},
	}

	result, err := r.client.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	claims := make([]model.Claim, 0, len(result.Items))
	for _, item := range result.Items {
		var c model.Claim
		if err := dynamodbattribute.UnmarshalMap(item, &c); err != nil {
			continue
		}
		claims = append(claims, c)
	}
	return claims, nil
}

func (r *DynamoBillingRepository) SaveInvoice(ctx context.Context, invoice *model.Invoice) error {
	item, err := dynamodbattribute.MarshalMap(invoice)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.invoiceTable),
		Item:      item,
	}

	_, err = r.client.PutItemWithContext(ctx, input)
	return err
}
