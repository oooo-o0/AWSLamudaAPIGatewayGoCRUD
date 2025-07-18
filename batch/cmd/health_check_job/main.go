package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oooo-o0/kaigo-insurance-system/batch/config"
	"github.com/oooo-o0/kaigo-insurance-system/batch/job/health"
	"github.com/oooo-o0/kaigo-insurance-system/batch/repository"
	"github.com/oooo-o0/kaigo-insurance-system/batch/service"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils"
)

func main() {
	cfg := config.Load()
	logger := utils.NewLogger("health-check-job")

	// DynamoDBリポジトリ初期化
	auditRepo := repository.NewDynamoHealthRepository(cfg)

	// サービス初期化
	healthService := service.NewHealthService(auditRepo, logger)

	// ジョブ実行
	handler := health.NewHealthCheckJob(healthService, logger)

	lambda.Start(handler.Handle)
}
