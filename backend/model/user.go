package model

import "time"

type User struct {
	UserID    string    `json:"user_id" dynamodbav:"UserID"`
	Name      string    `json:"name" dynamodbav:"Name"`
	Email     string    `json:"email" dynamodbav:"Email"`
	Role      string    `json:"role" dynamodbav:"Role"`
	CreatedAt time.Time `json:"created_at" dynamodbav:"CreatedAt"`
	UpdatedAt time.Time `json:"updated_at" dynamodbav:"UpdatedAt"`
}
