package controllers

import (
	"employee-api/modules/users/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FindUserByIdController struct {
	userService *services.UserService
}

func NewUserController(s *services.UserService) *FindUserByIdController{
	return &FindUserByIdController{
		userService: s,
	}
}

func (ctrl *FindUserByIdController) SayHello(c echo.Context) error {
	message := ctrl.userService.GetGreeting()
	return c.String(http.StatusOK, message)
}
