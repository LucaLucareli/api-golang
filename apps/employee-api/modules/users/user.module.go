package users

import (
	"employee-api/modules/users/controllers"
	"employee-api/modules/users/services"

	"github.com/labstack/echo/v4"
	"github.com/quantumsheep/plouf"
)

type UserModule struct {
	plouf.Module
}

func (m *UserModule) RegisterRoutes(e *echo.Group) {
	svc := services.NewUserService()
	ctrl := controllers.NewUserController(svc)

	e.GET("/oi", ctrl.SayHello)
}
