package repository

import (
	"context"

	"backend/model"
)

type ReportRepository interface {
	PutReport(ctx context.Context, report *model.Report) error
	// 必要に応じて他メソッド追加可（例: GetReportByUserID, DeleteOldReports など）
}

type BatchRepository interface {
	Report() ReportRepository
	// 他のリポジトリも追加可能（例: Billing(), Notification() など）
}
