package repository

import (
	"context"
	"time"

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

func (r *usersRepository) FindManyUsersToReport(ctx context.Context) (<-chan dto.FindManyUsersToReportStreamItem, error) {
	out := make(chan dto.FindManyUsersToReportStreamItem)

	go func() {
		defer close(out)
		const pageSize = 100
		var lastCreatedAt time.Time

		for {
			users, err := r.client.User.
				Query().
				Where(user.CreatedAtGT(lastCreatedAt)).
				Order(ent.Asc(user.FieldID)).
				Limit(pageSize).
				Select(user.FieldID, user.FieldName, user.FieldDocument, user.FieldCreatedAt).
				All(ctx)
			if err != nil {
				out <- dto.FindManyUsersToReportStreamItem{Err: err}
				return
			}

			if len(users) == 0 {
				break
			}

			for _, u := range users {
				out <- dto.FindManyUsersToReportStreamItem{
					User: dto.FindManyUsersToReportOutputDTO{
						ID:       u.ID,
						Name:     u.Name,
						Document: u.Document,
					},
				}
				lastCreatedAt = u.CreatedAt
			}
		}
	}()

	return out, nil
}
