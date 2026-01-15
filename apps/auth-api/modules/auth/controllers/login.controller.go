package controllers

import (
	"auth-api/modules/auth/dto/io"
	"auth-api/modules/auth/dto/request"
	"auth-api/modules/auth/services"
	"shared/interfaces"

	"github.com/labstack/echo/v4"
)

type LoginController struct {
	loginService *services.LoginService
}

func NewLoginController(s *services.LoginService) *LoginController {
	return &LoginController{
		loginService: s,
	}
}

func (ctrl *LoginController) LoginController(ctx echo.Context) error {
	var req request.LoginRequestDTO

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "payload inválido")
	}

	if err := ctx.Validate(&req); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	output, err := ctrl.loginService.LoginService(ctx.Request().Context(), io.LoginInputDTO{
		Document: req.Document,
		Password: req.Password,
	})
	if err != nil {
		return echo.NewHTTPError(401, err.Error())
	}

	interfaces.Set(ctx, interfaces.ResponseInterface[io.LoginOutputDTO]{
		Message: "Login realizado com sucesso",
		Result:  *output,
	})

	return nil
}
