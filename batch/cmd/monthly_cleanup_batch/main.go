package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oooo-o0/kaigo-insurance-system/batch/interface/job_handler"
	"github.com/oooo-o0/kaigo-insurance-system/batch/job/cleanup"
)

func main() {
	handler := job_handler.New()
	handler.Register("MonthlyCleanupJob", cleanup.ExecuteMonthlyCleanup)
	lambda.Start(handler.Handle)
}
