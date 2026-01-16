package modules

import (
	"auth-api/modules/auth"
	"shared"

	"github.com/labstack/echo/v4"
	"github.com/quantumsheep/plouf"
)

type AppModule struct {
	plouf.Module
	AuthModule *auth.AuthModule
}

func NewAppModule() *AppModule {
	return &AppModule{
		AuthModule: &auth.AuthModule{},
	}
}

func (m *AppModule) RegisterAllRoutes(e *echo.Echo, state *shared.AppState) {
	api := e.Group("/api")

	authGroup := api.Group("/auth")
	m.AuthModule.RegisterRoutes(authGroup, state)
}
