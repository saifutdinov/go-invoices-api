package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/saifutdinov/go-invoices-api/api/domain"
)

func (ih *InvoiceHandler) CreateInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	newInvoice := new(CreateInvoiceForm)
	if err := c.Bind(newInvoice); err != nil {
		return c.JSON(echo.ErrBadRequest.Code, err)
	}

	if err := ih.InvoiceUsecase.Save(ctx, domain.Invoice{
		ExternalID: newInvoice.ID,
		DueDate:    newInvoice.DueDate,
		Amount:     domain.ParseAmount(newInvoice.Amount.String()),
		Status:     domain.InvoiceStatus(newInvoice.Status),
	}); err != nil {
		return c.JSON(echo.ErrInternalServerError.Code, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"message": "ok"})
}
