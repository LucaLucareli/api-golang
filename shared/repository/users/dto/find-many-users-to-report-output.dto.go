package dto

import (
	"github.com/google/uuid"
)

type FindManyUsersToReportOutputDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Document string    `json:"document"`
}

type FindManyUsersToReportStreamItem struct {
	User FindManyUsersToReportOutputDTO
	Err  error
}
