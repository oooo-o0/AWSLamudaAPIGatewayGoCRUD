package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oooo-o0/kaigo-insurance-system/batch/config"
	"github.com/oooo-o0/kaigo-insurance-system/batch/job/notification"
	"github.com/oooo-o0/kaigo-insurance-system/batch/repository"
	"github.com/oooo-o0/kaigo-insurance-system/batch/service"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils"
)

func main() {
	// 環境変数や設定の読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Logger 初期化
	logger := utils.NewLogger(cfg.LogLevel)

	// DynamoDBリポジトリ初期化
	notificationRepo, err := repository.NewDynamoNotificationRepository(cfg.AWSRegion, cfg.NotificationTableName)
	if err != nil {
		logger.Fatalf("failed to create notification repository: %v", err)
	}

	// Service 初期化
	notificationService := service.NewNotificationService(notificationRepo, logger)

	// Job 初期化
	job := notification.NewDispatchJob(notificationService, logger)

	// Lambdaハンドラー登録
	lambda.Start(func(ctx context.Context) error {
		logger.Info("NotificationDispatch batch started")
		err := job.Execute(ctx)
		if err != nil {
			logger.Errorf("NotificationDispatch batch failed: %v", err)
			return err
		}
		logger.Info("NotificationDispatch batch completed successfully")
		return nil
	})
}
