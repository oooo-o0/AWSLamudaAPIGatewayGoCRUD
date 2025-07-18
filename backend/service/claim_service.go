package service

import (
	"context"
	"errors"

	"github.com/yourorg/kaigo/backend/model"
	"github.com/yourorg/kaigo/backend/repository"
)

type ClaimService interface {
	CreateClaim(ctx context.Context, claim *model.Claim) error
	GetClaim(ctx context.Context, claimID string) (*model.Claim, error)
	UpdateClaim(ctx context.Context, claim *model.Claim) error
	DeleteClaim(ctx context.Context, claimID string) error
	QueryClaimsByFacility(ctx context.Context, facilityID string) ([]*model.Claim, error)
	QueryClaimsByDate(ctx context.Context, claimDate string) ([]*model.Claim, error)
}

type claimService struct {
	repo repository.ClaimRepository
}

func NewClaimService(repo repository.ClaimRepository) ClaimService {
	return &claimService{
		repo: repo,
	}
}

func (s *claimService) CreateClaim(ctx context.Context, claim *model.Claim) error {
	if err := claim.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, claim)
}

func (s *claimService) GetClaim(ctx context.Context, claimID string) (*model.Claim, error) {
	if claimID == "" {
		return nil, errors.New("claimID is required")
	}
	return s.repo.GetByID(ctx, claimID)
}

func (s *claimService) UpdateClaim(ctx context.Context, claim *model.Claim) error {
	if err := claim.Validate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, claim)
}

func (s *claimService) DeleteClaim(ctx context.Context, claimID string) error {
	if claimID == "" {
		return errors.New("claimID is required")
	}
	return s.repo.Delete(ctx, claimID)
}

func (s *claimService) QueryClaimsByFacility(ctx context.Context, facilityID string) ([]*model.Claim, error) {
	if facilityID == "" {
		return nil, errors.New("facilityID is required")
	}
	return s.repo.QueryByFacilityID(ctx, facilityID)
}

func (s *claimService) QueryClaimsByDate(ctx context.Context, claimDate string) ([]*model.Claim, error) {
	if claimDate == "" {
		return nil, errors.New("claimDate is required")
	}
	return s.repo.QueryByClaimDate(ctx, claimDate)
}
