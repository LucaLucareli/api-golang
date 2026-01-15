package shared

import (
	"log"
	"os"
	"shared/auth"
	"shared/ent"
	repository "shared/repository/users"
	"strconv"
)

type AppState struct {
	AuthService *auth.AuthService
	UserRepo    repository.UsersRepository
}

func NewAppState(dbURL string) *AppState {

	client, err := ent.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}

	userRepo := repository.NewUsersRepository(client)

	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if accessSecret == "" {
		log.Fatal("JWT_ACCESS_SECRET não definido")
	}

	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		log.Fatal("JWT_REFRESH_SECRET não definido")
	}

	accessExpiryHours := 1
	if v := os.Getenv("JWT_ACCESS_EXPIRY_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			accessExpiryHours = n
		}
	}

	refreshExpiryDays := 7
	if v := os.Getenv("JWT_REFRESH_EXPIRY_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			refreshExpiryDays = n
		}
	}

	authSvc := &auth.AuthService{
		AccessSecret:      accessSecret,
		RefreshSecret:     refreshSecret,
		AccessExpiryHours: accessExpiryHours,
		RefreshExpiryDays: refreshExpiryDays,
		UserRepo:          userRepo,
	}

	return &AppState{
		AuthService: authSvc,
		UserRepo:    userRepo,
	}
}
