package service

import (
	"context"
	"fmt"
	"time"

	"kaigo-insurance-system/backend/model"
	"kaigo-insurance-system/batch/repository"
	"kaigo-insurance-system/batch/utils"
)

// BillingServiceInterface はBillingServiceのインターフェース定義
type BillingServiceInterface interface {
	GetUsersForBilling(ctx context.Context) ([]model.User, error)
	ProcessBillingForUser(ctx context.Context, user model.User) error
}

type BillingService struct {
	billingRepo repository.BillingRepositoryInterface
	logger      utils.LoggerInterface
}

func NewBillingService(repo repository.BillingRepositoryInterface, logger utils.LoggerInterface) *BillingService {
	return &BillingService{
		billingRepo: repo,
		logger:      logger,
	}
}

// GetUsersForBilling は日次請求対象ユーザー一覧を取得する
func (s *BillingService) GetUsersForBilling(ctx context.Context) ([]model.User, error) {
	// ここでは例として「Roleが'user'のアクティブユーザーを取得」などの簡単な条件で取得
	users, err := s.billingRepo.FetchUsersForBilling(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users for billing: %w", err)
	}
	return users, nil
}

// ProcessBillingForUser はユーザー単位で請求処理を行う
func (s *BillingService) ProcessBillingForUser(ctx context.Context, user model.User) error {
	// 1. そのユーザーの請求総額を計算する（例：ClaimのAmount合計）
	claims, err := s.billingRepo.FetchApprovedClaimsByUser(ctx, user.UserID)
	if err != nil {
		return fmt.Errorf("failed to fetch claims for user %s: %w", user.UserID, err)
	}

	var totalAmount float64
	for _, c := range claims {
		totalAmount += c.Amount
	}

	// 2. Invoiceレコード作成
	invoice := &model.Invoice{
		InvoiceID:    utils.GenerateUUID(),
		UserID:       user.UserID,
		Amount:       totalAmount,
		InvoiceMonth: time.Now().Format("2006-01"),
		Status:       "pending",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 3. Invoice保存
	if err := s.billingRepo.SaveInvoice(ctx, invoice); err != nil {
		return fmt.Errorf("failed to save invoice for user %s: %w", user.UserID, err)
	}

	s.logger.Infof("invoice created: invoice_id=%s, user_id=%s, amount=%.2f", invoice.InvoiceID, user.UserID, totalAmount)

	return nil
}
