package report

import (
	"context"
	"time"

	"github.com/oooo-o0/kaigo-insurance-system/batch/service"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils/logger"
)

func ExecuteWeeklyReportJob(ctx context.Context, event map[string]interface{}) error {
	log := logger.WithContext(ctx).Named("WeeklyReportJob")

	// 今週の対象期間を算出
	now := time.Now().UTC()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	log.Info("週次レポート生成処理開始", "start", startOfWeek, "end", endOfWeek)

	reportService := service.NewReportService()

	// 全ユーザー分の週間レポートを生成
	users, err := reportService.ListTargetUsers(ctx)
	if err != nil {
		log.Error("ユーザー取得失敗", "error", err)
		return err
	}

	for _, user := range users {
		if err := reportService.GenerateWeeklyReport(ctx, user, startOfWeek, endOfWeek); err != nil {
			log.Warn("ユーザーレポート生成失敗", "userID", user.UserID, "error", err)
		}
	}

	log.Info("週次レポート生成処理完了")
	return nil
}
