package http

import (
	"net/http"
	"uuid"

	"github.com/labstack/echo/v4"
	"github.com/saifutdinov/go-invoices-api/api/domain"
)

func (ph *PaymentHandler) CreatePayment(c echo.Context) error {
	ctx := c.Request().Context()
	newPayment := new(CreatePaymentForm)
	if err := c.Bind(newPayment); err != nil {
		return c.JSON(echo.ErrBadRequest.Code, err)
	}

	if err := ph.PaymentUsecase.Save(ctx, domain.Payment{
		ExternalID:  newPayment.ID,
		PaymentDate: newPayment.PaymentDate,
		Reference:   newPayment.Reference,
		Source:      newPayment.Source,
		Amount:      domain.ParseAmount(newPayment.Amount.String()),
		Status:      domain.PaymentStatusPending,
	}); err != nil {
		return c.JSON(echo.ErrInternalServerError.Code, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"message": "ok"})
}

func (ph *PaymentHandler) UpdateStatus(c echo.Context) error {

	ctx := c.Request().Context()
	paymentID, _ := uuid.Parse(c.Param("id"))
	newStatus := new(UpdatePaymentStatusForm)
	if err := c.Bind(newStatus); err != nil {
		return c.JSON(echo.ErrBadRequest.Code, err)
	}

	if err := ph.PaymentUsecase.UpdateStatus(ctx, paymentID, newStatus.Status); err != nil {
		return c.JSON(echo.ErrInternalServerError.Code, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"message": "ok"})

}
