package main

import (
	"log"

	"kaigo-insurance-system/batch/interface/job_handler"
	"kaigo-insurance-system/batch/job/careplan"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	handler := job_handler.New()
	handler.Register("careplan_auto_expire", careplan.AutoExpireCarePlansJob)
	log.Println("CarePlanAutoExpire lambda started")
	lambda.Start(handler.Execute)
}
