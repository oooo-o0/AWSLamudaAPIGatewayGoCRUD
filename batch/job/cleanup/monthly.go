package cleanup

import (
	"context"
	"time"

	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/batch/config"
	"github.com/oooo-o0/kaigo-insurance-system/batch/service"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils"
)

func ExecuteMonthlyCleanup(ctx context.Context) error {
	logger := utils.NewLogger("MonthlyCleanupJob")
	cfg := config.Load()

	svc := service.NewCleanupService(cfg)

	expireBefore := time.Now().AddDate(0, 0, -90)
	logger.Infof("開始: ReportType=daily かつ GeneratedDate < %s を削除", expireBefore.Format(time.RFC3339))

	// 実行
	count, err := svc.DeleteExpiredReports(ctx, model.ReportTypeDaily, expireBefore)
	if err != nil {
		logger.Errorf("削除処理失敗: %v", err)
		return err
	}

	logger.Infof("正常終了: 削除件数=%d", count)
	return nil
}
