package main

import (
	"fmt"
	"io"
	"net/http"
	"shared"

	"shared/helpers"
	"shared/interceptors"
	"shared/logger"
	"shared/validation"
	exceptionfactory "shared/validation/exception-factory"

	"auth-api/modules"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	AppName          = "AuthApi"
	DefaultDB        = "postgresql://postgres:pass@localhost:5432/auth_db?sslmode=disable"
	DefaultPort      = "3001"
	DefaultRedisPort = "1234"
	DefaultEnv       = "DEV"
)

func init() {
	env := helpers.GetEnv("LOG", DefaultEnv)
	logger.Init(AppName, logger.ColorRed, env)
}

func main() {
	dbURL := helpers.GetEnv("DATABASE_URL", DefaultDB)

	redisPort := helpers.GetEnv("REDIS_CACHE_PORT", DefaultRedisPort)
	redisURL := fmt.Sprintf("localhost:%s", redisPort)

	port := helpers.GetEnv("AUTH_API_PORT", DefaultPort)

	log.Info().Msgf("Starting %s application...", AppName)

	appState := shared.NewAppState(dbURL, redisURL)

	e := echo.New()

	e.HideBanner = true
	e.Logger.SetOutput(io.Discard)

	e.Use(interceptors.TransformInterceptor)

	e.Use(interceptors.RequestLogger)

	e.Validator = validation.NewValidator()

	e.Use(validation.ValidationMiddleware)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if validationErr := exceptionfactory.CustomExceptionFactory(err); validationErr != nil {
			_ = c.JSON(validationErr.Code, validationErr.Message)
			return
		}

		e.DefaultHTTPErrorHandler(err, c)
	}

	app := modules.NewAppModule()
	app.RegisterAllRoutes(e, appState)

	address := fmt.Sprintf(":%s", port)
	log.Info().Msgf("Application is running on: http://localhost:%s/api", port)

	logger.PrintRoutes(e)

	if err := e.Start(address); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Erro fatal no servidor")
	}
}
