package repository

import (
	"context"

	"shared/ent"
	"shared/ent/users"
)

type usersRepository struct {
	client *ent.Client
}

func NewUsersRepository(client *ent.Client) UsersRepository {
	return &usersRepository{client: client}
}

func (r *usersRepository) FindByID(ctx context.Context, id int) (*Users, error) {
	u, err := r.client.Users.
		Query().
		Where(users.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return &Users{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}, nil
}
