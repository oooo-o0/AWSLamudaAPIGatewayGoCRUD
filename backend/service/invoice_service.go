package service

import (
	"context"

	"github.com/oooo-o0/kaigo-insurance-system/backend/model"
	"github.com/oooo-o0/kaigo-insurance-system/backend/repository"
)

type InvoiceService struct {
	Repo repository.InvoiceRepository
}

func NewInvoiceService(repo repository.InvoiceRepository) *InvoiceService {
	return &InvoiceService{Repo: repo}
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, invoice *model.Invoice) error {
	return s.Repo.Put(ctx, invoice)
}

func (s *InvoiceService) GetInvoice(ctx context.Context, invoiceID string) (*model.Invoice, error) {
	return s.Repo.GetByID(ctx, invoiceID)
}

func (s *InvoiceService) ListInvoicesByUser(ctx context.Context, userID string) ([]*model.Invoice, error) {
	return s.Repo.ListByUserID(ctx, userID)
}

func (s *InvoiceService) MarkInvoiceSent(ctx context.Context, invoiceID string) error {
	return s.Repo.UpdateStatus(ctx, invoiceID, "sent")
}
