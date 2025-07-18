package service

import (
	"context"
	"time"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/repository"
)

type AuditLogService interface {
	Create(ctx context.Context, log *model.AuditLog) error
	QueryByEventType(ctx context.Context, eventType string, from, to time.Time) ([]model.AuditLog, error)
	QueryByUserID(ctx context.Context, userID string, from, to time.Time) ([]model.AuditLog, error)
}

type auditLogService struct {
	repo repository.AuditLogRepository
}

func NewAuditLogService(repo repository.AuditLogRepository) AuditLogService {
	return &auditLogService{repo: repo}
}

func (s *auditLogService) Create(ctx context.Context, log *model.AuditLog) error {
	// 省略可能だがID生成などを行う場合はここで
	return s.repo.Create(ctx, log)
}

func (s *auditLogService) QueryByEventType(ctx context.Context, eventType string, from, to time.Time) ([]model.AuditLog, error) {
	return s.repo.QueryByEventType(ctx, eventType, from, to)
}

func (s *auditLogService) QueryByUserID(ctx context.Context, userID string, from, to time.Time) ([]model.AuditLog, error) {
	return s.repo.QueryByUserID(ctx, userID, from, to)
}
