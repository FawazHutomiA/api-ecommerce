package auth

import (
	"example/internal/entity"
	"example/internal/repository/postgresql/role"
	"example/internal/repository/postgresql/token"
	"example/internal/repository/postgresql/user"
	"example/pkg/app"
	"example/pkg/bcrypt"
	"example/pkg/exception"
	"example/pkg/jwt"
	"example/pkg/middleware"
	"example/pkg/response"
	"example/pkg/sqlx"

	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, params AuthRegisterRequest) (resp *AuthRegisterResponse, errData exception.Error)
	Login(ctx context.Context, params AuthLoginRequest) (resp *AuthLoginResponse, errData exception.Error)
}

type authService struct {
	app        app.AppConfig
	repository user.UserRepository
	tokenRepo  token.TokenRepository
	roleRepo   role.RoleRepository
}

func NewAuthService(app app.AppConfig, repository user.UserRepository, tokenRepo token.TokenRepository, roleRepo role.RoleRepository) AuthService {
	return &authService{
		app:        app,
		repository: repository,
		tokenRepo:  tokenRepo,
		roleRepo:   roleRepo,
	}
}

func (uc *authService) Register(ctx context.Context, params AuthRegisterRequest) (resp *AuthRegisterResponse, errData exception.Error) {
	// Check User
	_, err := uc.repository.UserFindByEmail(ctx, params.Email)
	switch err {
	case nil:
		return nil, exception.Error{
			Status:  response.StatusConflicted,
			Message: "This Email Has Registered, Please Login",
			Errors:  exception.ErrConflicted,
		}
	case sql.ErrNoRows:
		err = nil
	default:
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	// init data
	userID := uuid.New()
	tokenID := uuid.New()

	roleRepo, err := uc.roleRepo.RoleFindByID(ctx, params.RoleID)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Failed to get data role by id",
			Errors:  exception.ErrBadRequest,
		}
	}

	// hash password
	hashedPassword, err := bcrypt.HashPassword(10, params.Password)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Failed to hash password",
			Errors:  exception.ErrInternalServer,
		}
	}

	paramsToken := jwt.DataToken{
		UserID: userID,
		Role:   roleRepo.Name,
	}

	jwtToken, err := jwt.GenerateToken(paramsToken)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Failed to generate token",
			Errors:  exception.ErrBadRequest,
		}
	}

	user := entity.User{
		ID:          userID,
		RoleID:      params.RoleID,
		WarehouseID: params.WarehouseID,
		Name:        params.Name,
		Email:       params.Email,
		Password:    &hashedPassword,
		Phone:       params.Phone,
		Gender:      params.Gender,
		Birth:       params.Birth,
		IsActive:    true,
	}

	// Transaction
	tx, err := sqlx.BeginTx(uc.app.Db, ctx)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	// save to db
	err = uc.repository.UserInsert(ctx, tx, user)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Failed to insert data user",
			Errors:  exception.ErrInternalServer,
		}
	}

	tokenValidate, err := jwt.ValidateToken(jwtToken.Token)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Wrong token",
			Errors:  exception.ErrBadRequest,
		}
	}

	var claims middleware.Claims
	claimsBytes, err := json.Marshal(tokenValidate.Claims)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Wrong token",
			Errors:  exception.ErrBadRequest,
		}
	}
	json.Unmarshal(claimsBytes, &claims)

	// Konversi expiredAt dari int64 ke time.Time
	expiredAt := time.Unix(claims.Exp, 0)

	// Insert Token
	dataToken := entity.Token{
		ID:        tokenID,
		UserID:    userID,
		Token:     jwtToken.Token,
		ExpiredAt: expiredAt,
	}

	err = uc.tokenRepo.TokenInsert(ctx, tx, dataToken)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Failed to insert data token",
			Errors:  exception.ErrBadRequest,
		}
	}

	err = sqlx.Commit(tx, ctx)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}
	// End Transaction

	resp = &AuthRegisterResponse{
		ExpiredAt: jwtToken.Exp,
		Token:     jwtToken.Token,
	}

	return resp, errData
}

func (uc *authService) Login(ctx context.Context, params AuthLoginRequest) (resp *AuthLoginResponse, errData exception.Error) {
	userRepo, err := uc.repository.UserFindByEmail(ctx, params.Email)
	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return nil, exception.Error{
			Status:  response.StatusUnauthorized,
			Message: "Invalid Email / Password",
			Errors:  exception.ErrUnauthorized,
		}
	default:
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	userRoleRepo, err := uc.repository.UserRoleFindByEmail(ctx, params.Email)
	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return nil, exception.Error{
			Status:  response.StatusUnauthorized,
			Message: "Invalid Email / Password",
			Errors:  exception.ErrUnauthorized,
		}
	default:
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	// Check Password
	valid := bcrypt.ComparePasswordHash(params.Password, *userRoleRepo.Password)
	if !valid {
		return resp, exception.Error{
			Status:  response.StatusUnauthorized,
			Message: "Invalid Password",
			Errors:  exception.ErrUnauthorized,
		}
	}

	if !userRepo.IsActive {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Your account has not active yet",
			Errors:  exception.ErrBadRequest,
		}
	}

	paramsToken := jwt.DataToken{
		UserID: userRepo.ID,
		Role:   userRoleRepo.Role,
	}

	jwtToken, err := jwt.GenerateToken(paramsToken)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Failed to generate token",
			Errors:  exception.ErrBadRequest,
		}
	}

	tokenValidate, err := jwt.ValidateToken(jwtToken.Token)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Wrong token",
			Errors:  exception.ErrBadRequest,
		}
	}

	var claims middleware.Claims
	claimsBytes, err := json.Marshal(tokenValidate.Claims)
	if err != nil {
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Wrong token",
			Errors:  exception.ErrBadRequest,
		}
	}
	json.Unmarshal(claimsBytes, &claims)

	// Konversi expiredAt dari int64 ke time.Time
	expiredAt := time.Unix(claims.Exp, 0)

	tokenRepo, err := uc.tokenRepo.TokenFindByUserID(ctx, userRepo.ID)
	switch err {
	case nil:
		err = nil
		dataToken := entity.Token{
			ID:        tokenRepo.ID,
			UserID:    userRepo.ID,
			Token:     jwtToken.Token,
			ExpiredAt: expiredAt,
		}

		err = uc.tokenRepo.TokenUpdate(ctx, dataToken)
		if err != nil {
			return nil, exception.Error{
				Status:  response.StatusBadRequest,
				Message: "Failed to update token",
				Errors:  exception.ErrBadRequest,
			}
		}
	case sql.ErrNoRows:
		// Transaction
		tx, err := sqlx.BeginTx(uc.app.Db, ctx)
		if err != nil {
			return resp, exception.Error{
				Status:  response.StatusInternalServerError,
				Message: "Error",
				Errors:  exception.ErrInternalServer,
			}
		}

		tokenID := uuid.New()

		// Insert Token
		dataToken := entity.Token{
			ID:        tokenID,
			UserID:    userRepo.ID,
			Token:     jwtToken.Token,
			ExpiredAt: expiredAt,
		}

		err = uc.tokenRepo.TokenInsert(ctx, tx, dataToken)
		if err != nil {
			return nil, exception.Error{
				Status:  response.StatusBadRequest,
				Message: "Failed to inser data token",
				Errors:  exception.ErrBadRequest,
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
	default:
		return nil, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	resp = &AuthLoginResponse{
		ExpiredAt: jwtToken.Exp,
		Token:     jwtToken.Token,
	}

	return resp, errData
}
