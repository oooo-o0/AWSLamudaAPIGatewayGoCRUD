package health

import (
	"context"

	"github.com/oooo-o0/kaigo-insurance-system/batch/interface/service"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils"
)

type HealthCheckJob struct {
	service service.HealthService
	logger  utils.Logger
}

func NewHealthCheckJob(service service.HealthService, logger utils.Logger) *HealthCheckJob {
	return &HealthCheckJob{
		service: service,
		logger:  logger,
	}
}

func (h *HealthCheckJob) Handle(ctx context.Context) error {
	h.logger.Info("HealthCheckJob started")

	// 任意の内部診断（データストア、設定、タイムチェック等）
	if err := h.service.CheckHealth(ctx); err != nil {
		h.logger.Error("Health check failed", "error", err)
		return err
	}

	h.logger.Info("HealthCheckJob completed successfully")
	return nil
}
