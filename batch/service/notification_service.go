package service

import (
	"context"
	"errors"

	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/batch/repository"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils"
)

type NotificationService struct {
	repo   repository.NotificationRepository
	logger utils.Logger
}

func NewNotificationService(repo repository.NotificationRepository, logger utils.Logger) *NotificationService {
	return &NotificationService{
		repo:   repo,
		logger: logger,
	}
}

// FetchUnreadNotifications 未読通知を上限limit件取得
func (s *NotificationService) FetchUnreadNotifications(ctx context.Context, limit int) ([]model.Notification, error) {
	notifications, err := s.repo.QueryByStatus(ctx, model.StatusUnread, limit)
	if err != nil {
		s.logger.Errorf("repo.QueryByStatus failed: %v", err)
		return nil, err
	}
	return notifications, nil
}

// UpdateNotificationStatus ステータス更新を行う
func (s *NotificationService) UpdateNotificationStatus(ctx context.Context, notification *model.Notification) error {
	if notification.NotificationID == "" {
		return errors.New("notification ID is empty")
	}
	err := s.repo.UpdateStatus(ctx, notification.NotificationID, notification.Status, notification.UpdatedAt)
	if err != nil {
		s.logger.Errorf("repo.UpdateStatus failed: %v", err)
		return err
	}
	return nil
}
