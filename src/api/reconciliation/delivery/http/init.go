package http

import (
	"github.com/labstack/echo/v4"
	"github.com/saifutdinov/go-invoices-api/api/domain"
	"github.com/saifutdinov/go-invoices-api/pkg/toml"
)

type ReconciliationHandler struct {
	ReconciliationUsecase domain.ReconciliationUsecase
	Config                *toml.Config
}

func InitReconciliationHandler(
	echoServer *echo.Echo,
	reconciliationUsecase domain.ReconciliationUsecase,
	config *toml.Config,
) {

	handler := ReconciliationHandler{
		ReconciliationUsecase: reconciliationUsecase,
		Config:                config,
	}

	echoServer.GET("/reconciliation", handler.ReportPage)

	apiGroup := echoServer.Group("/api")

	// reconciliation report
	apiGroup.GET("/reconciliation/report", handler.Report)

}
