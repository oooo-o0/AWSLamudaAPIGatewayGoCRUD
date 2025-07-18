package model

import "time"

type Invoice struct {
	InvoiceID    string    `json:"invoice_id" dynamodbav:"InvoiceID"`
	UserID       string    `json:"user_id" dynamodbav:"UserID"`
	Amount       float64   `json:"amount" dynamodbav:"Amount"`
	InvoiceMonth string    `json:"invoice_month" dynamodbav:"InvoiceMonth"` // e.g., "2025-07"
	PDFURL       string    `json:"pdf_url,omitempty" dynamodbav:"PDFURL,omitempty"`
	Status       string    `json:"status" dynamodbav:"Status"` // "pending", "sent", etc.
	CreatedAt    time.Time `json:"created_at" dynamodbav:"CreatedAt"`
	UpdatedAt    time.Time `json:"updated_at" dynamodbav:"UpdatedAt"`
}
