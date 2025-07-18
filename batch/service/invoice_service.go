package service

import (
	"batch/utils/logger"
	"context"
	"fmt"
	"time"

	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
)

type InvoiceService struct{}

func NewInvoiceService() *InvoiceService {
	return &InvoiceService{}
}

func (s *InvoiceService) SendInvoiceEmail(ctx context.Context, user *model.User, invoice *model.Invoice) error {
	logger := logger.Get()

	subject := fmt.Sprintf("【請求書】%s様の%s分", user.Name, invoice.IssueDate.Format("2006年1月"))
	body := fmt.Sprintf("請求額: ¥%.2f\n詳細: %s\n発行日: %s",
		invoice.Amount, invoice.Description, invoice.IssueDate.Format("2006-01-02"))

	logger.Infof("Sending invoice email to %s\nSubject: %s\nBody: %s", user.Email, subject, body)

	// 本番はSESやSendGridなどのメールサービスを呼び出す
	// ダミー処理
	time.Sleep(100 * time.Millisecond)

	return nil
}
