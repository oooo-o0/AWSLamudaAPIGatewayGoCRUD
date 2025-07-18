package facility

import (
	"context"
	"time"

	"batch/interface/repository"
	"batch/service"
	"batch/utils/logger"
)

func AggregateFacilityUsage(ctx context.Context, repo repository.FacilityRepository, svc service.FacilityService) error {
	logger.Info("Starting Facility Usage Aggregation")

	today := time.Now().Format("2006-01-02")
	claims, err := repo.FetchClaimsForDate(ctx, today)
	if err != nil {
		logger.Error("Failed to fetch claims", err)
		return err
	}

	summary := svc.CalculateUsageSummary(claims)

	if err := repo.StoreFacilityUsageSummary(ctx, today, summary); err != nil {
		logger.Error("Failed to store usage summary", err)
		return err
	}

	logger.Info("Facility Usage Aggregation completed successfully")
	return nil
}
