package model

import "time"

type CarePlan struct {
	CarePlanID  string    `json:"care_plan_id" dynamodbav:"CarePlanID"`
	UserID      string    `json:"user_id" dynamodbav:"UserID"`
	Title       string    `json:"title" dynamodbav:"Title"`
	Description string    `json:"description" dynamodbav:"Description"`
	StartDate   string    `json:"start_date" dynamodbav:"StartDate"`
	EndDate     string    `json:"end_date" dynamodbav:"EndDate"`
	Services    []string  `json:"services" dynamodbav:"Services"`
	CreatedAt   time.Time `json:"created_at" dynamodbav:"CreatedAt"`
	UpdatedAt   time.Time `json:"updated_at" dynamodbav:"UpdatedAt"`
}
