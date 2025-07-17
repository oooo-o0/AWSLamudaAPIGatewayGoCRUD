package service

import (
	"context"
	"errors"
	"strings"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/repository"

	"github.com/google/uuid"
)

type FacilityService interface {
	CreateFacility(ctx context.Context, f *model.Facility) (*model.Facility, error)
	GetFacility(ctx context.Context, id string) (*model.Facility, error)
	UpdateFacility(ctx context.Context, f *model.Facility) (*model.Facility, error)
	DeleteFacility(ctx context.Context, id string) error
	FindFacilitiesByName(ctx context.Context, name string) ([]*model.Facility, error)
}

type facilityService struct {
	repo repository.FacilityRepository
}

func NewFacilityService(repo repository.FacilityRepository) FacilityService {
	return &facilityService{repo: repo}
}

func (s *facilityService) CreateFacility(ctx context.Context, f *model.Facility) (*model.Facility, error) {
	// バリデーション
	if strings.TrimSpace(f.Name) == "" {
		return nil, errors.New("facility name is required")
	}

	// UUID発行
	f.FacilityID = uuid.New().String()

	err := s.repo.Create(ctx, f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *facilityService) GetFacility(ctx context.Context, id string) (*model.Facility, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *facilityService) UpdateFacility(ctx context.Context, f *model.Facility) (*model.Facility, error) {
	if strings.TrimSpace(f.FacilityID) == "" {
		return nil, errors.New("facility ID is required")
	}
	if strings.TrimSpace(f.Name) == "" {
		return nil, errors.New("facility name is required")
	}

	err := s.repo.Update(ctx, f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *facilityService) DeleteFacility(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("facility ID is required")
	}
	return s.repo.Delete(ctx, id)
}

func (s *facilityService) FindFacilitiesByName(ctx context.Context, name string) ([]*model.Facility, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name is required")
	}
	return s.repo.ListByName(ctx, name)
}
