package main

import (
	"log"

	"github.com/oooo-o0/kaigo-insurance-system/batch/config"
	"github.com/oooo-o0/kaigo-insurance-system/batch/interface/job_handler"
	"github.com/oooo-o0/kaigo-insurance-system/batch/job/user"
	"github.com/oooo-o0/kaigo-insurance-system/batch/repository"
	"github.com/oooo-o0/kaigo-insurance-system/batch/service"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// 設定読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// DynamoDBリポジトリ初期化
	userRepo := repository.NewDynamoUserRepository(cfg.DynamoDBTableUser)

	// サービス初期化
	userSvc := service.NewUserService(userRepo)

	// ジョブ初期化
	inactiveCleanupJob := user.NewInactiveUserCleanupJob(userSvc)

	// JobHandlerに登録し、Lambdaハンドラとして起動
	handler := job_handler.NewJobHandler()
	handler.Register("InactiveUserCleanup", inactiveCleanupJob)

	lambda.Start(handler.Execute)
}
