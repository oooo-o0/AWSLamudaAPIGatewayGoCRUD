package service

import (
	"context"
	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/repository"
)

type MigrationStateService interface {
	Create(ctx context.Context, state *model.MigrationState) error
	GetLatest(ctx context.Context) (*model.MigrationState, error)
}

type migrationStateService struct {
	repo repository.MigrationStateRepository
}

func NewMigrationStateService(repo repository.MigrationStateRepository) MigrationStateService {
	return &migrationStateService{repo: repo}
}

func (s *migrationStateService) Create(ctx context.Context, state *model.MigrationState) error {
	return s.repo.Create(ctx, state)
}

func (s *migrationStateService) GetLatest(ctx context.Context) (*model.MigrationState, error) {
	return s.repo.GetLatest(ctx)
}
