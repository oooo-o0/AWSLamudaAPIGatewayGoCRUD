package billing

import (
	"context"
	"fmt"

	"kaigo-insurance-system/batch/service"
	"kaigo-insurance-system/batch/utils"
)

type DailyBillingJob struct {
	service service.BillingServiceInterface
	logger  utils.LoggerInterface
}

func NewDailyBillingJob(srv service.BillingServiceInterface, logger utils.LoggerInterface) *DailyBillingJob {
	return &DailyBillingJob{
		service: srv,
		logger:  logger,
	}
}

// Execute は日次請求バッチのメイン処理
func (j *DailyBillingJob) Execute(ctx context.Context) error {
	j.logger.Info("DailyBillingJob started")

	// 1. 対象日の請求対象ユーザー取得
	users, err := j.service.GetUsersForBilling(ctx)
	if err != nil {
		return fmt.Errorf("failed to get users for billing: %w", err)
	}
	j.logger.Infof("found %d users to bill", len(users))

	// 2. ユーザーごとに請求データ作成・DynamoDB保存
	for _, user := range users {
		j.logger.Infof("processing billing for user_id=%s", user.UserID)

		err := j.service.ProcessBillingForUser(ctx, user)
		if err != nil {
			j.logger.Errorf("failed to process billing for user_id=%s: %v", user.UserID, err)
			// 続行：一部ユーザー失敗でもバッチ停止しない
			continue
		}
		j.logger.Infof("billing processed for user_id=%s", user.UserID)
	}

	j.logger.Info("DailyBillingJob completed")

	return nil
}
