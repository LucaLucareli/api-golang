package report

import (
	"employee-api/modules/report/controllers"
	reports "employee-api/modules/report/services"
	"shared"

	"shared/enums"
	"shared/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/quantumsheep/plouf"
)

type ReportModule struct {
	plouf.Module
}

func (m *ReportModule) RegisterRoutes(e *echo.Group, state *shared.AppState) {
	reportService := reports.NewRequestReportService(state.AsynqClient)
	reportController := controllers.NewRequestReportController(reportService)

	e.POST(
		"/report",
		middlewares.RequireAccess(
			state.AuthService,
			enums.AccessGroupEmployee,
		)(reportController.RequestReportController),
	)
}
