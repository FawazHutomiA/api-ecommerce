package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"example/internal/entity"
	"example/internal/repository/postgresql/token"
	"example/internal/repository/postgresql/user"
	"example/pkg/app"
	"example/pkg/bcrypt"
	"example/pkg/exception"
	"example/pkg/jwt"
	"example/pkg/middleware"
	"example/pkg/response"
	"example/pkg/sqlx"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, params AuthRegisterRequest) (resp AuthRegisterResponse, errData exception.Error)
	Login(ctx context.Context, params AuthLoginRequest) (resp AuthLoginResponse, errData exception.Error)
}

type authService struct {
	app        app.AppConfig
	repository user.UserRepository
	tokenRepo  token.TokenRepository
}

func NewAuthService(app app.AppConfig, repository user.UserRepository, tokenRepo token.TokenRepository) AuthService {
	return &authService{
		app:        app,
		repository: repository,
		tokenRepo:  tokenRepo,
	}
}

func (uc *authService) Register(ctx context.Context, params AuthRegisterRequest) (resp AuthRegisterResponse, errData exception.Error) {
	// Check User
	_, err := uc.repository.UserFindByEmail(ctx, params.Email)
	switch err {
	case nil:
		return resp, exception.Error{
			Status:  response.StatusConflicted,
			Message: "This Email Has Registered, Please Login",
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

	// init data
	userID := uuid.New()
	tokenID := uuid.New()
	roleName := "admin"
	roleID := "df302e3e-2256-488a-86ab-cb3ebbbab046"
	roleIDParse, err := uuid.Parse(roleID)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	// hash password
	hashedPassword, err := bcrypt.HashPassword(10, params.Password)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	paramsToken := jwt.DataToken{
		UserID: userID,
		Role:   roleName,
	}

	jwtToken, err := jwt.GenerateToken(paramsToken)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	user := entity.User{
		ID:          userID,
		RoleID:      roleIDParse,
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

	tokenValidate, err := jwt.ValidateToken(jwtToken.Token)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Wrong token",
			Errors:  exception.ErrBadRequest,
		}
	}

	var claims middleware.Claims
	claimsBytes, err := json.Marshal(tokenValidate.Claims)
	if err != nil {
		return resp, exception.Error{
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
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
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

	resp = AuthRegisterResponse{
		ExpiredAt: jwtToken.Exp,
		Token:     jwtToken.Token,
	}

	return resp, errData
}

func (uc *authService) Login(ctx context.Context, params AuthLoginRequest) (resp AuthLoginResponse, errData exception.Error) {
	userRepo, err := uc.repository.UserFindByEmail(ctx, params.Email)
	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return resp, exception.Error{
			Status:  response.StatusUnauthorized,
			Message: "Invalid Email / Password",
			Errors:  exception.ErrUnauthorized,
		}
	default:
		return resp, exception.Error{
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
		return resp, exception.Error{
			Status:  response.StatusUnauthorized,
			Message: "Invalid Email / Password",
			Errors:  exception.ErrUnauthorized,
		}
	default:
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	paramsToken := jwt.DataToken{
		UserID: userRepo.ID,
		Role:   userRoleRepo.Role,
	}

	jwtToken, err := jwt.GenerateToken(paramsToken)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	tokenValidate, err := jwt.ValidateToken(jwtToken.Token)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Wrong token",
			Errors:  exception.ErrBadRequest,
		}
	}

	var claims middleware.Claims
	claimsBytes, err := json.Marshal(tokenValidate.Claims)
	if err != nil {
		return resp, exception.Error{
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
			return resp, exception.Error{
				Status:  response.StatusBadRequest,
				Message: "Something Wrong",
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
			return resp, exception.Error{
				Status:  response.StatusBadRequest,
				Message: "Something Wrong",
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
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	resp = AuthLoginResponse{
		ExpiredAt: jwtToken.Exp,
		Token:     jwtToken.Token,
	}

	return resp, errData
}
