package users

import (
	"employee-api/modules/users/controllers"
	"employee-api/modules/users/services"
	"shared"

	"shared/enums"
	"shared/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/quantumsheep/plouf"
)

type UserModule struct {
	plouf.Module
}

func (m *UserModule) RegisterRoutes(e *echo.Group, state *shared.AppState) {
	svc := services.NewUserService()
	ctrl := controllers.NewUserController(svc)

	e.GET("/oi", middlewares.RequireAccess(state.AuthService, enums.AccessGroupEmployee)(ctrl.SayHello))
}
