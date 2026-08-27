package http

import (
	"github.com/labstack/echo/v4"
	"github.com/saifutdinov/go-invoices-api/api/domain"
	"github.com/saifutdinov/go-invoices-api/pkg/toml"
)

type PaymentHandler struct {
	PaymentUsecase domain.PaymentUsecase
	Config         *toml.Config
}

func InitPaymentHandlers(echoServer *echo.Echo, paymentUsecase domain.PaymentUsecase, config *toml.Config) {
	handler := PaymentHandler{
		PaymentUsecase: paymentUsecase,
		Config:         config,
	}

	apiGroup := echoServer.Group("/api")

	// create payment
	apiGroup.PUT("/payment/create", handler.CreatePayment)

	apiGroup.PATCH("/payment/:id/status", handler.UpdateStatus)
}
