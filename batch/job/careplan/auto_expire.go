package careplan

import (
	"context"
	"fmt"
	"time"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/batch/service"
	"kaigo-insurance-system/batch/utils/logger"
)

func AutoExpireCarePlansJob(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Info("Starting CarePlan auto-expiration job")

	careplanSvc := service.NewCarePlanService()
	now := time.Now().Format("2006-01-02")

	expiredPlans, err := careplanSvc.GetExpiredCarePlans(ctx, now)
	if err != nil {
		log.Error("Failed to fetch expired care plans", err)
		return err
	}

	log.Infof("Found %d expired care plans", len(expiredPlans))

	for _, plan := range expiredPlans {
		err := careplanSvc.MarkAsExpired(ctx, &plan)
		if err != nil {
			log.Errorf("Failed to mark CarePlan %s as expired: %v", plan.CarePlanID, err)
			continue
		}

		// 通知送信
		notification := &model.Notification{
			NotificationID: fmt.Sprintf("notif-%s", plan.CarePlanID),
			UserID:         plan.UserID,
			Title:          "介護プランの期限切れ",
			Message:        fmt.Sprintf("プラン『%s』は期限（%s）を過ぎました。", plan.Title, plan.EndDate),
			Status:         model.NotificationStatusUnread,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := careplanSvc.SendNotification(ctx, notification); err != nil {
			log.Errorf("Failed to send notification for CarePlan %s: %v", plan.CarePlanID, err)
		}
	}

	log.Info("CarePlan auto-expiration job completed")
	return nil
}
