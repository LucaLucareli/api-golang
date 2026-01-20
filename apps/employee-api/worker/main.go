package main

import (
	"fmt"
	"log"

	"employee-api/modules/report/consumer"
	reports "employee-api/modules/report/services"
	"shared"
	"shared/helpers"
	"shared/logger"
	"shared/queue"
)

const (
	DefaultDB        = "postgresql://postgres:pass@localhost:5432/employee_db?sslmode=disable"
	DefaultRedisPort = "1234"
	DefaultEnv       = "DEV"
	Name             = "EmployeeApi - QUEUE"
)

func init() {
	env := helpers.GetEnv("LOG", DefaultEnv)
	logger.Init(Name, logger.ColorPurple, env)
}

func main() {
	dbURL := helpers.GetEnv("DATABASE_URL", DefaultDB)

	redisPort := helpers.GetEnv("REDIS_CACHE_PORT", DefaultRedisPort)
	redisURL := fmt.Sprintf("localhost:%s", redisPort)

	state := shared.NewAppState(dbURL, redisURL)

	registry := consumer.NewReportProcessorRegistry(
		reports.NewProcessExcelReport(*state),
	)

	reportProcessor := consumer.NewReportProcessor(registry)

	mux := queue.NewMux(
		queue.QueueLifecycleOptions{
			QueueName: "reports",
			Inspector: state.AsynqInspector,
		},
		queue.TaskHandler{
			Task:    queue.TaskGenerateReport,
			Handler: reportProcessor.Handler(),
		},
	)

	log.Println("🚀 Report Worker started")

	if err := state.AsynqServer.Run(mux); err != nil {
		log.Fatal(err)
	}
}
