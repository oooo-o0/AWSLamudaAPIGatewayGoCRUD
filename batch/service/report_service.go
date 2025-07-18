package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/batch/repository"
)

type ReportService struct {
	reportRepo repository.ReportRepository
	userRepo   repository.UserRepository
	careRepo   repository.CarePlanRepository
	claimRepo  repository.ClaimRepository
}

func NewReportService() *ReportService {
	return &ReportService{
		reportRepo: repository.NewReportRepository(),
		userRepo:   repository.NewUserRepository(),
		careRepo:   repository.NewCarePlanRepository(),
		claimRepo:  repository.NewClaimRepository(),
	}
}

func (s *ReportService) ListTargetUsers(ctx context.Context) ([]*model.User, error) {
	return s.userRepo.FindAll(ctx)
}

func (s *ReportService) GenerateWeeklyReport(ctx context.Context, user *model.User, start, end time.Time) error {
	carePlans, err := s.careRepo.FindByUserIDAndDateRange(ctx, user.UserID, start, end)
	if err != nil {
		return fmt.Errorf("CarePlan取得失敗: %w", err)
	}

	claims, err := s.claimRepo.FindByUserIDAndDateRange(ctx, user.UserID, start, end)
	if err != nil {
		return fmt.Errorf("Claim取得失敗: %w", err)
	}

	reportContent := map[string]interface{}{
		"carePlans": carePlans,
		"claims":    claims,
	}

	contentJSON, err := json.Marshal(reportContent)
	if err != nil {
		return fmt.Errorf("JSON変換失敗: %w", err)
	}

	report := &model.Report{
		ReportID:      fmt.Sprintf("weekly-%s-%s", user.UserID, start.Format("20060102")),
		UserID:        user.UserID,
		ReportType:    model.ReportTypeWeekly,
		Title:         fmt.Sprintf("%sの週間レポート（%s〜%s）", user.Name, start.Format("2006-01-02"), end.Format("2006-01-02")),
		Content:       string(contentJSON),
		GeneratedDate: time.Now().UTC(),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	return s.reportRepo.Put(ctx, report)
}
