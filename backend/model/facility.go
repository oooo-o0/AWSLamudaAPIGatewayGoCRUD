package model

import (
	"time"
)

type Facility struct {
	FacilityID string    `json:"facility_id" dynamodbav:"FacilityID"`
	Name       string    `json:"name" dynamodbav:"Name"`
	Address    string    `json:"address,omitempty" dynamodbav:"Address,omitempty"`
	Phone      string    `json:"phone,omitempty" dynamodbav:"Phone,omitempty"`
	CreatedAt  time.Time `json:"created_at" dynamodbav:"CreatedAt"`
	UpdatedAt  time.Time `json:"updated_at" dynamodbav:"UpdatedAt"`
}
