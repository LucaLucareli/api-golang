package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"shared"

	"shared/container"
	"shared/interceptors"
	"shared/logger"
	"shared/validation"
	exceptionfactory "shared/validation/exception-factory"

	"employee-api/modules"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	AppName     = "EmployeeApi"
	DefaultDB   = "postgresql://postgres:pass@localhost:5432/employee_db?sslmode=disable"
	DefaultPort = "3003"
	DefaultEnv  = "DEV"
)

func init() {
	env := getEnv("LOG", DefaultEnv)
	logger.Init(AppName, logger.ColorPurple, env)
}

func main() {
	dbURL := getEnv("DATABASE_URL", DefaultDB)
	port := getEnv("EMPLOYEE_API_PORT", DefaultPort)

	log.Info().Msgf("Starting %s application...", AppName)

	_, _, err := container.Build(dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Falha ao criar container")
	}

	appState := shared.NewAppState(dbURL)

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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
