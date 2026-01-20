package shared

import (
	"log"
	"os"

	"shared/auth"
	"shared/database"
	"shared/ent"
	"shared/helpers"
	"shared/queue"
	repository "shared/repository/users"

	"github.com/hibiken/asynq"
)

type AppState struct {
	AuthService    *auth.AuthService
	UserRepo       repository.UsersRepository
	AsynqClient    *asynq.Client
	AsynqConfig    database.AsynqConfig
	AsynqServer    *asynq.Server
	AsynqInspector *asynq.Inspector
}

func NewAppState(dbPostgresURL, dbRedisURL string) *AppState {

	clientPostgres, err := ent.Open("postgres", dbPostgresURL)
	if err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}

	userRepo := repository.NewUsersRepository(clientPostgres)

	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if accessSecret == "" {
		log.Fatal("JWT_ACCESS_SECRET não definido")
	}

	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		log.Fatal("JWT_REFRESH_SECRET não definido")
	}

	accessExpiryHours := helpers.GetEnv("JWT_ACCESS_EXPIRY_HOURS", 1)
	refreshExpiryDays := helpers.GetEnv("JWT_REFRESH_EXPIRY_DAYS", 7)

	authSvc := &auth.AuthService{
		AccessSecret:      accessSecret,
		RefreshSecret:     refreshSecret,
		AccessExpiryHours: accessExpiryHours,
		RefreshExpiryDays: refreshExpiryDays,
		UserRepo:          userRepo,
	}

	asynqCfg := database.AsynqConfig{
		Addr:     helpers.GetEnv("REDIS_ADDR", dbRedisURL),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       helpers.GetEnv("REDIS_DB", 0),
	}

	asynqClient := database.NewAsynqClient(asynqCfg)

	asynqServer := queue.NewAsynqServer(asynqCfg)

	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr: dbRedisURL,
	})

	return &AppState{
		AuthService:    authSvc,
		UserRepo:       userRepo,
		AsynqClient:    asynqClient,
		AsynqConfig:    asynqCfg,
		AsynqServer:    asynqServer,
		AsynqInspector: inspector,
	}
}
