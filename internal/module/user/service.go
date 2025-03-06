package user

import (
	"context"
	"database/sql"
	"example/internal/entity"
	"example/internal/repository/postgresql/token"
	"example/internal/repository/postgresql/user"
	"example/pkg/app"
	"example/pkg/bcrypt"
	"example/pkg/exception"
	"example/pkg/helper"
	"example/pkg/response"
	"example/pkg/sqlx"

	"github.com/google/uuid"
)

type UserService interface {
	ListPaginate(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, errData exception.Error)
	Detail(ctx context.Context, id uuid.UUID) (resp UserDetailResponse, errData exception.Error)
	Create(ctx context.Context, params UserCreateRequest) (resp UserCreateResponse, errData exception.Error)
	Update(ctx context.Context, id uuid.UUID, params UserUpdateRequest) (resp UserUpdateResponse, errData exception.Error)
	Delete(ctx context.Context, id uuid.UUID) (resp bool, errData exception.Error)
}

type userService struct {
	app        app.AppConfig
	repository user.UserRepository
	tokenRepo  token.TokenRepository
}

func NewUserService(app app.AppConfig, repository user.UserRepository, tokenRepo token.TokenRepository) UserService {
	return &userService{
		app:        app,
		repository: repository,
		tokenRepo:  tokenRepo,
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
	repoUser, err := uc.repository.UserFindByID(ctx, id)
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
		ID:       id,
		RoleID:   repoUser.RoleID,
		Name:     repoUser.Name,
		Email:    repoUser.Email,
		Phone:    repoUser.Phone,
		Gender:   repoUser.Gender,
		Birth:    repoUser.Birth,
		IsActive: repoUser.IsActive,
	}

	return resp, errData
}

func (uc *userService) Create(ctx context.Context, params UserCreateRequest) (resp UserCreateResponse, errData exception.Error) {
	// query user
	_, err := uc.repository.UserFindByEmail(ctx, params.Email)

	switch err {
	case nil:
		return resp, exception.Error{
			Status:  response.StatusConflicted,
			Message: "This Email Has Registered",
			Errors:  exception.ErrConflicted,
		}
	case sql.ErrNoRows:
		err = nil
	default:
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	// hash password
	hashedPassword, err := bcrypt.HashPassword(10, *params.Password)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	userID := uuid.New()

	// map insert
	user := entity.User{
		ID:       userID,
		RoleID:   params.RoleID,
		Name:     params.Name,
		Email:    params.Email,
		Password: &hashedPassword,
		Phone:    params.Phone,
		Gender:   params.Gender,
		Birth:    params.Birth,
		IsActive: params.IsActive,
	}

	// Transaction
	tx, err := sqlx.BeginTx(uc.app.Db, ctx)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	// save to db
	err = uc.repository.UserInsert(ctx, tx, user)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	err = sqlx.Commit(tx, ctx)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}
	// End Transaction

	// map response
	resp = UserCreateResponse{
		ID:       userID,
		RoleID:   params.RoleID,
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
		Phone:    params.Phone,
		Gender:   params.Gender,
		Birth:    params.Birth,
		IsActive: params.IsActive,
	}

	return resp, errData
}

func (uc *userService) Update(ctx context.Context, id uuid.UUID, params UserUpdateRequest) (resp UserUpdateResponse, errData exception.Error) {
	// query user
	_, err := uc.repository.UserFindByID(ctx, id)

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

	// hash password
	hashedPassword, err := bcrypt.HashPassword(10, *params.Password)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	// map insert
	user := entity.User{
		ID:       id,
		RoleID:   params.RoleID,
		Name:     params.Name,
		Email:    params.Email,
		Password: &hashedPassword,
		Phone:    params.Phone,
		Gender:   params.Gender,
		Birth:    params.Birth,
		IsActive: params.IsActive,
	}

	err = uc.repository.UserUpdate(ctx, user)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Error",
			Errors:  exception.ErrBadRequest,
		}
	}

	// map response
	resp = UserUpdateResponse{
		ID:       id,
		RoleID:   params.RoleID,
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
		Phone:    params.Phone,
		Gender:   params.Gender,
		Birth:    params.Birth,
		IsActive: params.IsActive,
	}

	return resp, errData
}

func (uc *userService) Delete(ctx context.Context, id uuid.UUID) (resp bool, errData exception.Error) {
	// query user
	userRepo, err := uc.repository.UserFindByID(ctx, id)

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

	// Transaction
	tx, err := sqlx.BeginTx(uc.app.Db, ctx)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	// update to db
	err = uc.repository.UserDelete(ctx, tx, id)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Error",
			Errors:  exception.ErrBadRequest,
		}
	}

	_, errV := uc.tokenRepo.TokenFindByUserID(ctx, userRepo.ID)
	if errV == nil {
		err = uc.tokenRepo.TokenDeleteByUserID(ctx, tx, userRepo.ID)
		if err != nil {
			return resp, exception.Error{
				Status:  response.StatusBadRequest,
				Message: "Error",
				Errors:  exception.ErrBadRequest,
			}
		}
	}

	err = sqlx.Commit(tx, ctx)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}
	// End Transaction

	return resp, errData
}
