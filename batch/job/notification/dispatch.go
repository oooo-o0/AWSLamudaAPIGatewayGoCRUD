package notification

import (
	"context"
	"time"

	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/batch/service"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils"
)

type DispatchJob struct {
	notificationService *service.NotificationService
	logger              utils.Logger
}

func NewDispatchJob(ns *service.NotificationService, logger utils.Logger) *DispatchJob {
	return &DispatchJob{
		notificationService: ns,
		logger:              logger,
	}
}

// Execute 未送信通知の取得→送信→ステータス更新を行う
func (j *DispatchJob) Execute(ctx context.Context) error {
	// 未読(UNREAD)通知を取得（上限100件など制限可）
	notifications, err := j.notificationService.FetchUnreadNotifications(ctx, 100)
	if err != nil {
		j.logger.Errorf("failed to fetch unread notifications: %v", err)
		return err
	}

	if len(notifications) == 0 {
		j.logger.Info("no unread notifications to dispatch")
		return nil
	}

	for _, n := range notifications {
		j.logger.Infof("dispatching notification ID=%s userID=%s", n.NotificationID, n.UserID)

		// 送信処理 (ここではモックとして、実際はメール・プッシュAPI呼び出しなど)
		err := j.sendNotification(ctx, n)
		if err != nil {
			j.logger.Errorf("failed to send notification ID=%s: %v", n.NotificationID, err)
			continue // エラーでもループ続行
		}

		// ステータスをREADに更新
		n.Status = model.StatusRead
		n.UpdatedAt = time.Now().UTC()

		err = j.notificationService.UpdateNotificationStatus(ctx, &n)
		if err != nil {
			j.logger.Errorf("failed to update notification status ID=%s: %v", n.NotificationID, err)
		}
	}

	return nil
}

// sendNotification 実際の通知送信処理（例: メール送信やプッシュ通知）
// 今回はモック実装。将来的に外部サービス連携等を実装
func (j *DispatchJob) sendNotification(ctx context.Context, n model.Notification) error {
	// 送信API呼び出し例:
	// err := externalAPI.SendEmail(n.UserID, n.Title, n.Message)
	// if err != nil { return err }

	j.logger.Infof("mock send notification to userID=%s title=%s", n.UserID, n.Title)

	// ここは成功モック
	return nil
}
