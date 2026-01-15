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
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "JSON inválido"})
	}

	if err := ctx.Validate(&req); err != nil {
		if httpErr := exceptionfactory.CustomExceptionFactory(err); httpErr != nil {
			return ctx.JSON(httpErr.Code, httpErr.Message)
		}
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "erro desconhecido"})
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
