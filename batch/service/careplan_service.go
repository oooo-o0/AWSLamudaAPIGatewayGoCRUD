package service

import (
	"context"
	"time"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/batch/repository"
)

type CarePlanService struct {
	careplanRepo     repository.CarePlanRepository
	notificationRepo repository.NotificationRepository
}

func NewCarePlanService() *CarePlanService {
	return &CarePlanService{
		careplanRepo:     repository.NewCarePlanRepository(),
		notificationRepo: repository.NewNotificationRepository(),
	}
}

// 期限切れプランの取得
func (s *CarePlanService) GetExpiredCarePlans(ctx context.Context, date string) ([]model.CarePlan, error) {
	return s.careplanRepo.FindCarePlansBeforeDate(ctx, date)
}

// ステータス更新 & 更新時刻変更
func (s *CarePlanService) MarkAsExpired(ctx context.Context, plan *model.CarePlan) error {
	plan.Description += "（期限切れ）"
	plan.UpdatedAt = time.Now()
	return s.careplanRepo.UpdateCarePlan(ctx, plan)
}

// 通知送信
func (s *CarePlanService) SendNotification(ctx context.Context, notif *model.Notification) error {
	return s.notificationRepo.PutNotification(ctx, notif)
}
