package role

import (
	"context"
	"database/sql"
	"example/internal/repository/postgresql/role"
	"example/internal/repository/postgresql/token"
	"example/pkg/app"
	"example/pkg/exception"
	"example/pkg/helper"
	"example/pkg/response"

	"github.com/google/uuid"
)

type RoleService interface {
	ListPaginate(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, errData exception.Error)
	Detail(ctx context.Context, id uuid.UUID) (resp RoleDetailResponse, errData exception.Error)
}

type roleService struct {
	app        app.AppConfig
	repository role.RoleRepository
	tokenRepo  token.TokenRepository
}

func NewRoleService(app app.AppConfig, repository role.RoleRepository, tokenRepo token.TokenRepository) RoleService {
	return &roleService{
		app:        app,
		repository: repository,
		tokenRepo:  tokenRepo,
	}
}

func (uc *roleService) ListPaginate(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, errData exception.Error) {
	// repository
	repo, err := uc.repository.RoleFindAll(ctx, params)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	return repo, errData
}

func (uc *roleService) Detail(ctx context.Context, id uuid.UUID) (resp RoleDetailResponse, errData exception.Error) {
	// repository
	roleRepo, err := uc.repository.RoleFindByID(ctx, id)
	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return resp, exception.Error{
			Status:  response.StatusNotFound,
			Message: "Role Not Found",
			Errors:  exception.ErrNotFound,
		}
	default:
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	// map response
	resp = RoleDetailResponse(roleRepo)

	return resp, errData
}
