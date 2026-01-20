package controllers

import (
	"employee-api/modules/report/dto/io"
	"employee-api/modules/report/dto/request"
	services "employee-api/modules/report/services"
	"net/http"
	"shared/interfaces"
	exceptionfactory "shared/validation/exception-factory"

	"github.com/labstack/echo/v4"
)

type RequestReportController struct {
	refreshTokenService *services.RequestReportService
}

func NewRequestReportController(s *services.RequestReportService) *RequestReportController {
	return &RequestReportController{
		refreshTokenService: s,
	}
}

func (ctrl *RequestReportController) RequestReportController(ctx echo.Context) error {
	var req request.RequestReportDTO

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "JSON inválido"})
	}

	if err := ctx.Validate(&req); err != nil {
		if httpErr := exceptionfactory.CustomExceptionFactory(err); httpErr != nil {
			return ctx.JSON(httpErr.Code, httpErr.Message)
		}
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "erro desconhecido"})
	}

	err := ctrl.refreshTokenService.Execute(
		ctx.Request().Context(),
		io.RequestReportInputDTO{
			ReportTypeID: req.ReportTypeID,
			UserID:       req.UserID,
		},
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	interfaces.Set(ctx, interfaces.ResponseInterface[any]{
		Message: "Relatório solicitado",
	})

	return nil
}
