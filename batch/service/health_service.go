package service

import (
	"context"
	"time"

	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/batch/interface/repository"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils"
)

type HealthService interface {
	CheckHealth(ctx context.Context) error
}

type healthService struct {
	auditRepo repository.AuditLogRepository
	logger    utils.Logger
}

func NewHealthService(auditRepo repository.AuditLogRepository, logger utils.Logger) HealthService {
	return &healthService{
		auditRepo: auditRepo,
		logger:    logger,
	}
}

func (s *healthService) CheckHealth(ctx context.Context) error {
	now := time.Now()

	// 簡単なDynamoへの書き込みでヘルス確認（または複数の診断も可）
	log := model.AuditLog{
		LogID:       utils.GenerateULID(),
		UserID:      "system",
		EventType:   "HEALTH_CHECK",
		Description: "Scheduled health check executed",
		Timestamp:   now,
		Metadata:    "{}",
	}

	if err := s.auditRepo.PutAuditLog(ctx, log); err != nil {
		return err
	}

	s.logger.Info("AuditLog created for health check")
	return nil
}
