package role

import (
	"example/internal/entity"
	"example/pkg/app"
	"example/pkg/helper"
	sqlxPkg "example/pkg/sqlx"

	"context"
	"database/sql"

	"github.com/google/uuid"
)

type RoleRepository interface {
	RoleFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error)
	RoleFindByID(ctx context.Context, id uuid.UUID) (resp entity.Role, err error)
	RoleInsert(ctx context.Context, tx *sql.Tx, role entity.Role) (err error)
	RoleUpdate(ctx context.Context, role entity.Role) (err error)
	RoleDelete(ctx context.Context, tx *sql.Tx, id uuid.UUID) (err error)
}

type roleRepository struct {
	app app.AppConfig
}

func NewRoleRepository(app app.AppConfig) RoleRepository {
	return &roleRepository{
		app: app,
	}
}

func (repo *roleRepository) RoleFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error) {
	query := FIND_ALL

	var role []entity.Role
	pagination := sqlxPkg.NewPaginationMetadata(repo.app.Db)
	result, err := pagination.GetPagination(query, params, &role)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return result, nil
}

func (repo *roleRepository) RoleFindByID(ctx context.Context, id uuid.UUID) (resp entity.Role, err error) {
	err = repo.app.Db.GetContext(ctx, &resp, FIND_BY_ID, id)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *roleRepository) RoleInsert(ctx context.Context, tx *sql.Tx, role entity.Role) (err error) {
	_, err = tx.ExecContext(ctx, INSERT, role.ToInsert()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *roleRepository) RoleUpdate(ctx context.Context, role entity.Role) (err error) {
	_, err = repo.app.Db.ExecContext(ctx, UPDATE_BY_ID, role.ToUpdate()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *roleRepository) RoleDelete(ctx context.Context, tx *sql.Tx, id uuid.UUID) (err error) {
	_, err = tx.ExecContext(ctx, DELETE_BY_ID, id)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}
