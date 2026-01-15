package dto

import "github.com/google/uuid"

type FindUserToLoginOutputDTO struct {
	ID             uuid.UUID `json:"id"`
	Password       string    `json:"password"`
	Name           string    `json:"name"`
	AccessGroupIds []int16   `json:"accessGroupIds"`
}
