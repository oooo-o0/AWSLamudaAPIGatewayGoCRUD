package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oooo-o0/kaigo-insurance-system/batch/interface/job_handler"
	"github.com/oooo-o0/kaigo-insurance-system/batch/job/report"
	"github.com/oooo-o0/kaigo-insurance-system/batch/utils/logger"
)

func main() {
	logger.Initialize() // Zap logger初期化

	jobHandler := job_handler.NewJobHandler()
	jobHandler.Register("WeeklyReport", report.ExecuteWeeklyReportJob)

	lambda.Start(jobHandler.Execute)
}
