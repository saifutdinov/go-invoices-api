package api

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"
	invoiceHttp "github.com/saifutdinov/go-invoices-api/api/invoice/delivery/http"
	invoicePsql "github.com/saifutdinov/go-invoices-api/api/invoice/repository/psql"
	invoceUsecase "github.com/saifutdinov/go-invoices-api/api/invoice/usecase"
	logPsql "github.com/saifutdinov/go-invoices-api/api/log/repository/psql"

	recoHttp "github.com/saifutdinov/go-invoices-api/api/reconciliation/delivery/http"
	recoPsql "github.com/saifutdinov/go-invoices-api/api/reconciliation/repository/psql"
	recoUsercase "github.com/saifutdinov/go-invoices-api/api/reconciliation/usecase"

	paymentHttp "github.com/saifutdinov/go-invoices-api/api/payment/delivery/http"
	paymentPsql "github.com/saifutdinov/go-invoices-api/api/payment/repository/psql"
	paymentUsecase "github.com/saifutdinov/go-invoices-api/api/payment/usecase"
	"github.com/saifutdinov/go-invoices-api/pkg/chronos"
	"github.com/saifutdinov/go-invoices-api/pkg/db"
	"github.com/saifutdinov/go-invoices-api/pkg/toml"
)

func initHandlers(echoServer *echo.Echo, dbo *sql.DB, config *toml.Config) {

	echoServer.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	tx := db.NewTransaction(dbo)

	// log repository
	logRepo := logPsql.NewLogRepository(dbo)

	// invoice
	invoiceRepo := invoicePsql.NewInvoiceRepository(dbo)
	invoiceUsecase := invoceUsecase.NewInvoiceUsecase(tx, logRepo, invoiceRepo, config)
	invoiceHttp.InitInvoiceHandlers(echoServer, invoiceUsecase, config)

	// payment
	paymentRepo := paymentPsql.NewPaymentRepository(dbo)
	paymentUsecase := paymentUsecase.NewPaymentUsecase(tx, logRepo, paymentRepo, config)
	paymentHttp.InitPaymentHandlers(echoServer, paymentUsecase, config)

	// reconciliation
	recoRepo := recoPsql.NewReconciliationRepository(dbo)

	// chronos
	ctx := context.Background()

	chronosScheduler := chronos.New(ctx)
	reconciliactionUsecase := recoUsercase.NewReconciliationUsecase(tx, logRepo, recoRepo, invoiceRepo, paymentRepo, chronosScheduler)

	recoHttp.InitReconciliationHandler(echoServer, reconciliactionUsecase, config)

	reconciliactionUsecase.StartBackgroundProcessing()
	reconciliactionUsecase.ProcessPayments(ctx)

}
