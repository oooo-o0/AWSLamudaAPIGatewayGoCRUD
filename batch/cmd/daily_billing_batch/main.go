package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	"kaigo-insurance-system/batch/config"
	"kaigo-insurance-system/batch/job/billing"
	"kaigo-insurance-system/batch/repository"
	"kaigo-insurance-system/batch/service"
	"kaigo-insurance-system/batch/utils"
)

func main() {
	// Lambdaのエントリポイントに登録
	lambda.Start(handler)
}

func handler(ctx context.Context) error {
	// 環境変数・設定読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("failed to load config: %v", err)
		return err
	}

	// ロガー初期化
	logger := utils.NewLogger(cfg.LogLevel)

	// DynamoDBリポジトリ初期化
	billingRepo, err := repository.NewDynamoBillingRepository(cfg.DynamoDBTableBilling)
	if err != nil {
		logger.Errorf("failed to create billing repository: %v", err)
		return err
	}

	// BillingService初期化
	billingService := service.NewBillingService(billingRepo, logger)

	// ジョブ作成
	job := billing.NewDailyBillingJob(billingService, logger)

	logger.Info("starting DailyBillingBatch")

	// バッチ処理実行
	if err := job.Execute(ctx); err != nil {
		logger.Errorf("DailyBillingBatch execution failed: %v", err)
		return err
	}

	logger.Info("DailyBillingBatch completed successfully")

	return nil
}
