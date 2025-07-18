package service

import (
	"context"
	"errors"
	"time"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/backend/repository"

	"github.com/google/uuid"
)

type ReportService interface {
	Create(ctx context.Context, report *model.Report) error
	Get(ctx context.Context, id string) (*model.Report, error)
	QueryByTypeAndDate(ctx context.Context, reportType model.ReportType, from, to time.Time) ([]model.Report, error)
	Update(ctx context.Context, report *model.Report) error
	Delete(ctx context.Context, id string) error
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) Create(ctx context.Context, report *model.Report) error {
	if report.UserID == "" || report.Title == "" || report.ReportType == "" {
		return errors.New("missing required fields")
	}
	report.ReportID = uuid.NewString()
	report.GeneratedDate = time.Now()
	return s.repo.Create(ctx, report)
}

func (s *reportService) Get(ctx context.Context, id string) (*model.Report, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *reportService) QueryByTypeAndDate(ctx context.Context, reportType model.ReportType, from, to time.Time) ([]model.Report, error) {
	return s.repo.QueryByTypeAndDate(ctx, reportType, from, to)
}

func (s *reportService) Update(ctx context.Context, report *model.Report) error {
	if report.ReportID == "" {
		return errors.New("missing reportId")
	}
	return s.repo.Update(ctx, report)
}

func (s *reportService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
