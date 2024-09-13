package user

import (
	"context"
	"database/sql"
	"example/internal/repository/postgresql/user"
	"example/pkg/app"
	"example/pkg/exception"
	"example/pkg/helper"
	"example/pkg/response"

	"github.com/google/uuid"
)

type UserService interface {
	ListPaginate(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, errData exception.Error)
	Detail(ctx context.Context, id uuid.UUID) (resp UserDetailResponse, errData exception.Error)
}

type userService struct {
	app        app.AppConfig
	repository user.UserRepository
}

func NewUserService(app app.AppConfig, repository user.UserRepository) UserService {
	return &userService{
		app:        app,
		repository: repository,
	}
}

func (uc *userService) ListPaginate(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, errData exception.Error) {
	// repository
	repo, err := uc.repository.UserFindAll(ctx, params)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	return repo, errData
}

func (uc *userService) Detail(ctx context.Context, id uuid.UUID) (resp UserDetailResponse, errData exception.Error) {
	// repository
	repo, err := uc.repository.UserFindByID(ctx, id)
	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return resp, exception.Error{
			Status:  response.StatusNotFound,
			Message: "User Not Found",
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
	resp = UserDetailResponse{
		ID:         id,
		Name:       repo.Name,
		Email:      repo.Email,
		Occupation: repo.Occupation,
		Phone:      repo.Phone,
		Gender:     repo.Gender,
		Role:       repo.Role,
		IsGoogle:   repo.IsGoogle,
		IsActive:   repo.IsActive,
		IsVerify:   repo.IsVerify,
	}

	return resp, errData
}
