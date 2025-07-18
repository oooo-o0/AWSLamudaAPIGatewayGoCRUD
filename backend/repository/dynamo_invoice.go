package repository

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
)

type InvoiceRepository interface {
	Put(ctx context.Context, invoice *model.Invoice) error
	GetByID(ctx context.Context, invoiceID string) (*model.Invoice, error)
	ListByUserID(ctx context.Context, userID string) ([]*model.Invoice, error)
	UpdateStatus(ctx context.Context, invoiceID, status string) error
}

type DynamoInvoiceRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func NewDynamoInvoiceRepository(client *dynamodb.Client, tableName string) InvoiceRepository {
	return &DynamoInvoiceRepository{
		Client:    client,
		TableName: tableName,
	}
}

func (r *DynamoInvoiceRepository) Put(ctx context.Context, invoice *model.Invoice) error {
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	item, err := attributevalue.MarshalMap(invoice)
	if err != nil {
		return err
	}

	_, err = r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.TableName,
		Item:      item,
	})
	return err
}

func (r *DynamoInvoiceRepository) GetByID(ctx context.Context, invoiceID string) (*model.Invoice, error) {
	key, _ := attributevalue.MarshalMap(map[string]string{"InvoiceID": invoiceID})

	out, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.TableName,
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}

	var invoice model.Invoice
	err = attributevalue.UnmarshalMap(out.Item, &invoice)
	return &invoice, err
}

func (r *DynamoInvoiceRepository) ListByUserID(ctx context.Context, userID string) ([]*model.Invoice, error) {
	exprAttr := map[string]types.AttributeValue{
		":uid": &types.AttributeValueMemberS{Value: userID},
	}

	out, err := r.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &r.TableName,
		IndexName:                 awsString("UserIDIndex"),
		KeyConditionExpression:    awsString("UserID = :uid"),
		ExpressionAttributeValues: exprAttr,
	})
	if err != nil {
		return nil, err
	}

	var invoices []*model.Invoice
	err = attributevalue.UnmarshalListOfMaps(out.Items, &invoices)
	return invoices, err
}

func (r *DynamoInvoiceRepository) UpdateStatus(ctx context.Context, invoiceID, status string) error {
	now := time.Now().Format(time.RFC3339)

	_, err := r.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.TableName,
		Key: map[string]types.AttributeValue{
			"InvoiceID": &types.AttributeValueMemberS{Value: invoiceID},
		},
		UpdateExpression: awsString("SET #s = :s, UpdatedAt = :u"),
		ExpressionAttributeNames: map[string]string{
			"#s": "Status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":s": &types.AttributeValueMemberS{Value: status},
			":u": &types.AttributeValueMemberS{Value: now},
		},
	})
	return err
}

func awsString(s string) *string {
	return &s
}
