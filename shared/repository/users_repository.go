package repository

import (
	"context"
)

type Users struct {
	ID    int
	Name  string
	Email string
}

type UsersRepository interface {
	FindByID(ctx context.Context, id int) (*Users, error)
}

