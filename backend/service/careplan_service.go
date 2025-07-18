package service

import (
	"context"
	"time"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/repository"

	"github.com/google/uuid"
)

type CarePlanService interface {
	Create(ctx context.Context, plan *model.CarePlan) error
	Get(ctx context.Context, id string) (*model.CarePlan, error)
	Update(ctx context.Context, plan *model.CarePlan) error
	Delete(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userID string) ([]model.CarePlan, error)
}

type carePlanService struct {
	repo repository.CarePlanRepository
}

func NewCarePlanService(repo repository.CarePlanRepository) CarePlanService {
	return &carePlanService{repo: repo}
}

func (s *carePlanService) Create(ctx context.Context, plan *model.CarePlan) error {
	plan.CarePlanID = uuid.New().String()
	now := time.Now()
	plan.CreatedAt = now
	plan.UpdatedAt = now
	return s.repo.Create(ctx, plan)
}

func (s *carePlanService) Get(ctx context.Context, id string) (*model.CarePlan, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *carePlanService) Update(ctx context.Context, plan *model.CarePlan) error {
	plan.UpdatedAt = time.Now()
	return s.repo.Update(ctx, plan)
}

func (s *carePlanService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *carePlanService) ListByUser(ctx context.Context, userID string) ([]model.CarePlan, error) {
	return s.repo.ListByUserID(ctx, userID)
}
