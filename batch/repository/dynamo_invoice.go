package repository

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
)

type InvoiceRepository struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func NewInvoiceRepository(db *dynamodb.DynamoDB, tableName string) *InvoiceRepository {
	return &InvoiceRepository{db: db, tableName: tableName}
}

func (r *InvoiceRepository) FetchPendingInvoices(ctx context.Context) ([]model.Invoice, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(r.tableName),
		FilterExpression: aws.String("Sent = :sent"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sent": {BOOL: aws.Bool(false)},
		},
	}

	result, err := r.db.ScanWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	var invoices []model.Invoice
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &invoices)
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (r *InvoiceRepository) MarkAsSent(ctx context.Context, invoiceID string) error {
	now := time.Now()

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"InvoiceID": {S: aws.String(invoiceID)},
		},
		UpdateExpression: aws.String("SET Sent = :sent, UpdatedAt = :updated"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sent":    {BOOL: aws.Bool(true)},
			":updated": {S: aws.String(now.Format(time.RFC3339))},
		},
	}

	_, err := r.db.UpdateItemWithContext(ctx, input)
	return err
}
