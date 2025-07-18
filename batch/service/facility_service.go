package service

import (
	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
)

type FacilityService interface {
	CalculateUsageSummary(claims []model.Claim) map[string]model.FacilityUsageSummary
}

type facilityService struct{}

func NewFacilityService() FacilityService {
	return &facilityService{}
}

func (s *facilityService) CalculateUsageSummary(claims []model.Claim) map[string]model.FacilityUsageSummary {
	summary := make(map[string]model.FacilityUsageSummary)

	for _, claim := range claims {
		s := summary[claim.FacilityID]
		s.FacilityID = claim.FacilityID
		s.TotalClaims += 1
		s.TotalAmount += claim.Amount
		summary[claim.FacilityID] = s
	}

	return summary
}
