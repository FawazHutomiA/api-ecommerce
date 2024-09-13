package user

import (
	"context"
	"example/internal/entity"
	"example/pkg/app"
	"example/pkg/helper"
	sqlxPkg "example/pkg/sqlx"

	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UserRepository interface {
	UserFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error)
	UserFindByID(ctx context.Context, id uuid.UUID) (resp entity.User, err error)
	UserFindByEmail(ctx context.Context, email string) (resp entity.User, err error)
	UserInsert(ctx context.Context, user entity.User) (err error)
	UserUpdateTokenByID(ctx context.Context, user entity.User) (err error)
}

type userRepository struct {
	app app.AppConfig
}

func NewUserRepository(app app.AppConfig) UserRepository {
	return &userRepository{
		app: app,
	}
}

func (repo *userRepository) UserFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error) {
	query := FIND_ALL

	if params.Search != "" {
		escapedSearch := strings.Replace(params.Search, "'", "''", -1)
		addFilter := fmt.Sprintf("AND (u.full_name ILIKE '%%%s%%' OR u.email ILIKE '%%%s%%')", escapedSearch, escapedSearch)
		query = fmt.Sprintf("%s %s", query, addFilter)
	}

	if params.OrderBy != "" {
		params.SortBy = "u.created_at"
		query = fmt.Sprintf("%s ORDER BY %s %s", query, params.SortBy, params.OrderBy)
	}

	var user []entity.User
	pagination := sqlxPkg.NewPaginationMetadata(repo.app.Db)
	result, err := pagination.GetPagination(query, params, &user)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return result, nil
}

func (repo *userRepository) UserFindByID(ctx context.Context, id uuid.UUID) (resp entity.User, err error) {
	err = repo.app.Db.GetContext(ctx, &resp, FIND_BY_ID, id)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *userRepository) UserFindByEmail(ctx context.Context, email string) (resp entity.User, err error) {
	err = repo.app.Db.GetContext(ctx, &resp, FIND_BY_EMAIL, email)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *userRepository) UserInsert(ctx context.Context, user entity.User) (err error) {
	_, err = repo.app.Db.ExecContext(ctx, INSERT, user.ToInsert()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *userRepository) UserUpdateTokenByID(ctx context.Context, user entity.User) (err error) {
	_, err = repo.app.Db.ExecContext(ctx, UPDATE_TOKEN_USER, user.ToUpdate()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}
