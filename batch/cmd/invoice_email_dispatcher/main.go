package main

import (
	"batch/config"
	"batch/interface/job_handler"
	"batch/job/invoice"
	"batch/utils/logger"
)

func main() {
	logger.Init()
	cfg := config.Load()

	handler := job_handler.NewJobHandler()
	handler.Register("InvoiceEmailDispatch", invoice.Run)

	handler.Execute("InvoiceEmailDispatch", cfg)
}
