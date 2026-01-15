package dto

import "github.com/google/uuid"

type FindByIdOutputDTO struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Document       string    `json:"document"`
	AccessGroupIds []int16   `json:"accessGroupIds"`
}
