package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/batch/interface/repository"
	"github.com/oooo-o0/kaigo-insurance-system/batch/repository"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils"
)

type CleanupService interface {
	DeleteExpiredReports(ctx context.Context, reportType model.ReportType, before time.Time) (int, error)
}

type cleanupService struct {
	repo repository.DynamoCleanupRepository
}

func NewCleanupService(cfg *utils.Config) CleanupService {
	return &cleanupService{
		repo: repository.NewDynamoCleanupRepository(cfg),
	}
}

func (s *cleanupService) DeleteExpiredReports(ctx context.Context, reportType model.ReportType, before time.Time) (int, error) {
	reports, err := s.repo.FindExpiredReports(ctx, reportType, before)
	if err != nil {
		return 0, fmt.Errorf("レポート取得失敗: %w", err)
	}

	count := 0
	for _, r := range reports {
		if err := s.repo.DeleteReport(ctx, r.ReportID); err != nil {
			utils.NewLogger("CleanupService").Warnf("ReportID=%s の削除失敗: %v", r.ReportID, err)
			continue
		}

		log := model.AuditLog{
			LogID:       utils.UUID(),
			UserID:      r.UserID,
			EventType:   "DELETE",
			Description: fmt.Sprintf("古い日次レポートを削除: %s", r.Title),
			Timestamp:   time.Now(),
			Metadata:    fmt.Sprintf(`{"ReportID": "%s"}`, r.ReportID),
		}
		_ = s.repo.PutAuditLog(ctx, log)

		noti := model.Notification{
			NotificationID: utils.UUID(),
			UserID:         r.UserID,
			Title:          "古いレポートが削除されました",
			Message:        fmt.Sprintf("タイトル: %s は90日以上前のため自動削除されました。", r.Title),
			Status:         model.StatusUnread,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		_ = s.repo.PutNotification(ctx, noti)

		count++
	}

	return count, nil
}
