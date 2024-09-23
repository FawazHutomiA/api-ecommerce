package auth

import (
	"context"
	"database/sql"
	"example/internal/entity"
	"example/internal/repository/postgresql/user"
	"example/pkg/app"
	"example/pkg/exception"
	"example/pkg/jwt"
	"example/pkg/response"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, params AuthRegisterRequest) (resp AuthRegisterResponse, errData exception.Error)
	Login(ctx context.Context, params AuthLoginRequest) (resp AuthLoginResponse, errData exception.Error)
}

type authService struct {
	app        app.AppConfig
	repository user.UserRepository
}

func NewAuthService(app app.AppConfig, repository user.UserRepository) AuthService {
	return &authService{
		app:        app,
		repository: repository,
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

	userID := uuid.New()
	var paswordInput *string

	if !params.IsGoogle {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.MinCost)
		if err != nil {
			return resp, exception.Error{
				Status:  response.StatusBadRequest,
				Message: "Something Wrong",
				Errors:  exception.ErrBadRequest,
			}
		}
		passwordStr := string(passwordHash) // Convert []byte to string
		paswordInput = &passwordStr         // Assign pointer to passwordStr
	} else {
		paswordInput = nil
	}

	paramsToken := jwt.DataToken{
		UserID: userID,
		Role:   "user",
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
		ID:         userID,
		Name:       params.Name,
		Email:      params.Email,
		Occupation: params.Occupation,
		Password:   paswordInput,
		Phone:      params.Phone,
		Role:       "user",
		Gender:     params.Gender,
		Token:      jwtToken.Token,
	}
	err = uc.repository.UserInsert(ctx, user)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

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

	paramsToken := jwt.DataToken{
		UserID: userRepo.ID,
		Role:   userRepo.Role,
	}

	jwtToken, err := jwt.GenerateToken(paramsToken)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	dataToken := entity.User{
		ID:    userRepo.ID,
		Token: jwtToken.Token,
	}

	err = uc.repository.UserUpdateTokenByID(ctx, dataToken)
	if err != nil {
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
