package repository

import (
	"context"
	"kaigo-insurance-system/backend/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type MigrationStateRepository interface {
	Create(ctx context.Context, state *model.MigrationState) error
	GetLatest(ctx context.Context) (*model.MigrationState, error)
}

type migrationStateRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewMigrationStateRepository(db *dynamodb.Client, tableName string) MigrationStateRepository {
	return &migrationStateRepo{db: db, tableName: tableName}
}

func (r *migrationStateRepo) Create(ctx context.Context, state *model.MigrationState) error {
	item, err := attributevalue.MarshalMap(state)
	if err != nil {
		return err
	}
	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(MigrationID)"),
	})
	return err
}

// 最新マイグレーションを取得（スキャンして最新AppliedAtを選ぶ簡易実装）
func (r *migrationStateRepo) GetLatest(ctx context.Context) (*model.MigrationState, error) {
	out, err := r.db.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}
	var states []model.MigrationState
	err = attributevalue.UnmarshalListOfMaps(out.Items, &states)
	if err != nil {
		return nil, err
	}
	if len(states) == 0 {
		return nil, nil
	}

	latest := states[0]
	for _, s := range states[1:] {
		if s.AppliedAt.After(latest.AppliedAt) {
			latest = s
		}
	}
	return &latest, nil
}
