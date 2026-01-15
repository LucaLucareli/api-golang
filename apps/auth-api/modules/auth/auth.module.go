package auth

import (
	"auth-api/modules/auth/controllers"
	"auth-api/modules/auth/services"
	"shared"

	"github.com/labstack/echo/v4"
	"github.com/quantumsheep/plouf"
)

type AuthModule struct {
	plouf.Module
}

func (m *AuthModule) RegisterRoutes(e *echo.Group, state *shared.AppState) {
	loginService := services.NewLoginService(state.AuthService)
	loginController := controllers.NewLoginController(loginService)

	e.POST("/login", loginController.LoginController)
}
