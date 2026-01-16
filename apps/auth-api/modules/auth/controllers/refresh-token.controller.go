package controllers

import (
	"auth-api/modules/auth/dto/io"
	"auth-api/modules/auth/dto/request"
	"auth-api/modules/auth/services"
	"net/http"
	"shared/interfaces"
	exceptionfactory "shared/validation/exception-factory"

	"github.com/labstack/echo/v4"
)

type RefreshTokenController struct {
	refreshTokenService *services.RefreshTokenService
}

func NewRefreshTokenController(s *services.RefreshTokenService) *RefreshTokenController {
	return &RefreshTokenController{
		refreshTokenService: s,
	}
}

func (ctrl *RefreshTokenController) RefreshTokenController(ctx echo.Context) error {
	var req request.RefreshTokenRequestDTO

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "JSON inválido"})
	}

	if err := ctx.Validate(&req); err != nil {
		if httpErr := exceptionfactory.CustomExceptionFactory(err); httpErr != nil {
			return ctx.JSON(httpErr.Code, httpErr.Message)
		}
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "erro desconhecido"})
	}

	output, err := ctrl.refreshTokenService.RefreshTokenService(ctx.Request().Context(), io.RefreshTokenInputDTO{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return echo.NewHTTPError(401, err.Error())
	}

	interfaces.Set(ctx, interfaces.ResponseInterface[io.RefreshTokenOutputDTO]{
		Result: *output,
	})

	return nil
}
