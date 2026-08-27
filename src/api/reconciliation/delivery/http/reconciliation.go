package http

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (rh *ReconciliationHandler) Report(
	c echo.Context,
) error {

	report, err := rh.ReconciliationUsecase.Report(c.Request().Context())
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to build reconciliation report",
		)
	}

	return c.JSON(http.StatusOK, report)
}

func (rh *ReconciliationHandler) ReportPage(c echo.Context) error {
	return c.File("assets/templates/reconciliation.html")
}
