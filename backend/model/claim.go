package model

import (
	"errors"
	"time"
)

type Claim struct {
	ClaimID     string  `json:"claim_id"`
	FacilityID  string  `json:"facility_id"`
	UserID      string  `json:"user_id"`
	ClaimDate   string  `json:"claim_date"` // ISO8601 "YYYY-MM-DD"
	Amount      float64 `json:"amount"`
	Description string  `json:"description,omitempty"`
	Status      string  `json:"status"` // e.g. "Pending", "Approved", "Rejected"
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Validate performs basic claim data validation
func (c *Claim) Validate() error {
	if c.ClaimID == "" {
		return errors.New("ClaimID is required")
	}
	if c.FacilityID == "" {
		return errors.New("FacilityID is required")
	}
	if c.UserID == "" {
		return errors.New("UserID is required")
	}
	if _, err := time.Parse("2006-01-02", c.ClaimDate); err != nil {
		return errors.New("ClaimDate must be YYYY-MM-DD")
	}
	if c.Amount < 0 {
		return errors.New("Amount must be non-negative")
	}
	if c.Status == "" {
		return errors.New("Status is required")
	}
	return nil
}
