package service

import (
	"context"
	"errors"
	"time"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/repository"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, n *model.Notification) error
	GetNotification(ctx context.Context, notificationID string) (*model.Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
	DeleteNotification(ctx context.Context, notificationID string) error
	ListNotificationsByUser(ctx context.Context, userID string) ([]*model.Notification, error)
}

type notificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{
		repo: repo,
	}
}

func (s *notificationService) CreateNotification(ctx context.Context, n *model.Notification) error {
	if n.NotificationID == "" {
		return errors.New("notification_id required")
	}
	if n.UserID == "" {
		return errors.New("user_id required")
	}
	n.Status = model.StatusUnread
	now := time.Now().UTC()
	n.CreatedAt = now
	n.UpdatedAt = now
	return s.repo.Put(ctx, n)
}

func (s *notificationService) GetNotification(ctx context.Context, notificationID string) (*model.Notification, error) {
	return s.repo.GetByID(ctx, notificationID)
}

func (s *notificationService) MarkAsRead(ctx context.Context, notificationID string) error {
	return s.repo.UpdateStatus(ctx, notificationID, model.StatusRead)
}

func (s *notificationService) DeleteNotification(ctx context.Context, notificationID string) error {
	return s.repo.Delete(ctx, notificationID)
}

func (s *notificationService) ListNotificationsByUser(ctx context.Context, userID string) ([]*model.Notification, error) {
	return s.repo.QueryByUserID(ctx, userID)
}
