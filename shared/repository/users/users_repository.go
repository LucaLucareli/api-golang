package repository

import (
	"context"
	"shared/repository/users/dto"

	"github.com/google/uuid"
)

type UsersRepository interface {
	FindUserToLogin(ctx context.Context, document string) (*dto.FindUserToLoginOutputDTO, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.FindByIdOutputDTO, error)
	FindManyUsersToReport(ctx context.Context) (<-chan dto.FindManyUsersToReportStreamItem, error)
}
