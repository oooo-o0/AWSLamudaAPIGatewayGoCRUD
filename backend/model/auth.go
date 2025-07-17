package model

import "time"

type AuthToken struct {
	TokenID   string    `dynamodbav:"TokenID" json:"tokenId"`
	UserID    string    `dynamodbav:"UserID" json:"userId"`
	Expiry    time.Time `dynamodbav:"Expiry" json:"expiry"`
	IssuedAt  time.Time `dynamodbav:"IssuedAt" json:"issuedAt"`
	IsRevoked bool      `dynamodbav:"IsRevoked" json:"isRevoked"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type AuthResponse struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}
