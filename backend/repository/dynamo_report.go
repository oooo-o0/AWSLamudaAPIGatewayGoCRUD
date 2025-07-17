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

type ReportRepository interface {
	Create(ctx context.Context, report *model.Report) error
	GetByID(ctx context.Context, id string) (*model.Report, error)
	QueryByTypeAndDate(ctx context.Context, reportType model.ReportType, from, to time.Time) ([]model.Report, error)
	Update(ctx context.Context, report *model.Report) error
	Delete(ctx context.Context, id string) error
}

type reportRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewReportRepository(db *dynamodb.Client, tableName string) ReportRepository {
	return &reportRepo{db: db, tableName: tableName}
}

func (r *reportRepo) Create(ctx context.Context, report *model.Report) error {
	now := time.Now()
	report.CreatedAt = now
	report.UpdatedAt = now

	item, err := attributevalue.MarshalMap(report)
	if err != nil {
		return err
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(ReportID)"),
	})
	return err
}

func (r *reportRepo) GetByID(ctx context.Context, id string) (*model.Report, error) {
	out, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ReportID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil || out.Item == nil {
		return nil, errors.New("report not found")
	}
	var report model.Report
	err = attributevalue.UnmarshalMap(out.Item, &report)
	return &report, err
}

func (r *reportRepo) QueryByTypeAndDate(ctx context.Context, reportType model.ReportType, from, to time.Time) ([]model.Report, error) {
	exprAttrValues := map[string]types.AttributeValue{
		":rt":   &types.AttributeValueMemberS{Value: string(reportType)},
		":from": &types.AttributeValueMemberS{Value: from.Format(time.RFC3339)},
		":to":   &types.AttributeValueMemberS{Value: to.Format(time.RFC3339)},
	}

	out, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		IndexName:                 aws.String("ReportTypeIndex"),
		KeyConditionExpression:    aws.String("ReportType = :rt AND GeneratedDate BETWEEN :from AND :to"),
		ExpressionAttributeValues: exprAttrValues,
	})
	if err != nil {
		return nil, err
	}

	var reports []model.Report
	err = attributevalue.UnmarshalListOfMaps(out.Items, &reports)
	return reports, err
}

func (r *reportRepo) Update(ctx context.Context, report *model.Report) error {
	report.UpdatedAt = time.Now()
	item, err := attributevalue.MarshalMap(report)
	if err != nil {
		return err
	}
	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *reportRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ReportID": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}
