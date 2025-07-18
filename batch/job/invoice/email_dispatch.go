package invoice

import (
	"batch/interface/repository"
	"batch/interface/service"
	"batch/utils/logger"
	"context"
)

func Run(ctx context.Context, svc service.ServiceContainer, repo repository.RepositoryContainer) error {
	logger := logger.Get()

	invoices, err := repo.InvoiceRepository.FetchPendingInvoices(ctx)
	if err != nil {
		logger.Error("Failed to fetch invoices", err)
		return err
	}

	for _, invoice := range invoices {
		user, err := repo.UserRepository.GetUserByID(ctx, invoice.UserID)
		if err != nil {
			logger.Warnf("User not found for InvoiceID=%s, skipping: %v", invoice.InvoiceID, err)
			continue
		}

		err = svc.InvoiceService.SendInvoiceEmail(ctx, user, &invoice)
		if err != nil {
			logger.Error("Failed to send invoice email", err)
			continue
		}

		err = repo.InvoiceRepository.MarkAsSent(ctx, invoice.InvoiceID)
		if err != nil {
			logger.Error("Failed to update invoice status", err)
		}
	}
	return nil
}
