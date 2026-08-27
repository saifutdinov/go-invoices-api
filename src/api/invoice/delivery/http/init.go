package http

import (
	"github.com/labstack/echo/v4"
	"github.com/saifutdinov/go-invoices-api/api/domain"
	"github.com/saifutdinov/go-invoices-api/pkg/toml"
)

type InvoiceHandler struct {
	InvoiceUsecase domain.InvoiceUsecase
	Config         *toml.Config
}

func InitInvoiceHandlers(echoServer *echo.Echo, invoiceUsecase domain.InvoiceUsecase, config *toml.Config) {
	handler := InvoiceHandler{
		InvoiceUsecase: invoiceUsecase,
		Config:         config,
	}

	apiGroup := echoServer.Group("/api")

	// creates invoice record
	apiGroup.PUT("/invoice/create", handler.CreateInvoice)
}
