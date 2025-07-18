package main

import (
	"batch/config"
	"batch/interface/job_handler"
	"batch/job/facility"
	"batch/utils/logger"
)

func main() {
	logger.Init() // Zap logger init
	cfg := config.Load()

	handler := job_handler.NewJobHandler()
	handler.Register("FacilityUsageAggregator", facility.AggregateFacilityUsage)

	if err := handler.Execute("FacilityUsageAggregator", cfg); err != nil {
		logger.Fatal("Job execution failed", err)
	}
}
