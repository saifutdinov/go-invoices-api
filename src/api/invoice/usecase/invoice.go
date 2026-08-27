package usecase

import (
	"context"

	"github.com/saifutdinov/go-invoices-api/api/domain"
)

func (iu *InvoiceUsecase) Save(ctx context.Context, invoice domain.Invoice) error {
	return iu.InvoiceRepository.Save(ctx, invoice)
}
