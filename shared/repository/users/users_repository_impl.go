package repository

import (
	"context"

	"shared/ent"
	"shared/ent/user"
	"shared/ent/usersonaccessgroups"
	"shared/repository/users/dto"

	"github.com/google/uuid"
)

type usersRepository struct {
	client *ent.Client
}

func NewUsersRepository(client *ent.Client) UsersRepository {
	return &usersRepository{client: client}
}

func (r *usersRepository) FindUserToLogin(
	ctx context.Context,
	document string,
) (*dto.FindUserToLoginOutputDTO, error) {

	u, err := r.client.User.
		Query().
		Where(user.DocumentEQ(document)).
		Select(
			user.FieldID,
			user.FieldName,
			user.FieldPassword,
		).
		WithAccessGroups(func(q *ent.UsersOnAccessGroupsQuery) {
			q.Select(usersonaccessgroups.FieldAccessGroupID)
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	accessGroupIds := make([]int16, 0, len(u.Edges.AccessGroups))
	for _, ag := range u.Edges.AccessGroups {
		accessGroupIds = append(accessGroupIds, int16(ag.AccessGroupID))
	}

	return &dto.FindUserToLoginOutputDTO{
		ID:             u.ID,
		Name:           u.Name,
		Password:       u.Password,
		AccessGroupIds: accessGroupIds,
	}, nil
}

func (r *usersRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*dto.FindByIdOutputDTO, error) {

	u, err := r.client.User.
		Query().
		Where(user.IDEQ(id)).
		Select(
			user.FieldID,
			user.FieldName,
			user.FieldDocument,
		).
		WithAccessGroups(func(q *ent.UsersOnAccessGroupsQuery) {
			q.Select(usersonaccessgroups.FieldAccessGroupID)
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	accessGroupIds := make([]int16, 0, len(u.Edges.AccessGroups))
	for _, ag := range u.Edges.AccessGroups {
		accessGroupIds = append(accessGroupIds, int16(ag.AccessGroupID))
	}

	return &dto.FindByIdOutputDTO{
		ID:             u.ID,
		Name:           u.Name,
		Document:       u.Document,
		AccessGroupIds: accessGroupIds,
	}, nil
}
