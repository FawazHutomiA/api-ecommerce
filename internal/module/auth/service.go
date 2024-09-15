package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"example/internal/entity"
	"example/internal/repository/postgresql/user"
	"example/pkg/app"
	"example/pkg/email"
	"example/pkg/exception"
	"example/pkg/jwt"
	jwtValidate "example/pkg/jwt"
	"example/pkg/middleware"
	"example/pkg/response"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, params AuthRegisterRequest) (resp AuthRegisterResponse, errData exception.Error)
	Login(ctx context.Context, params AuthLoginRequest) (resp AuthLoginResponse, errData exception.Error)
	GetOrSaveUser(ctx context.Context, params GoogleUser) (resp AuthLoginResponse, errData exception.Error)
	VerifyUserEmail(ctx context.Context, token string) (resp AuthRegisterResponse, errData exception.Error)
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

	// Send verification email
	verificationLink := fmt.Sprintf("http://localhost:8080/api/v1/verify?token=%s", jwtToken.Token)
	err = email.SendVerificationEmailSMTP(params.Email, verificationLink)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Failed to send verification email",
			Errors:  exception.ErrInternalServer,
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
		IsGoogle:   false,
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

	if !userRepo.IsGoogle {
		// Check Password
		err = bcrypt.CompareHashAndPassword([]byte(*userRepo.Password), []byte(params.Password))
		if err != nil {
			return resp, exception.Error{
				Status:  response.StatusUnauthorized,
				Message: "Invalid Email / Password",
				Errors:  exception.ErrUnauthorized,
			}
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

func (uc *authService) GetOrSaveUser(ctx context.Context, params GoogleUser) (resp AuthLoginResponse, errData exception.Error) {
	userByEmail, _ := uc.repository.UserFindByEmail(ctx, params.Email)

	userID := uuid.New()

	if userByEmail.ID == uuid.Nil {
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
			ID:       userID,
			Name:     params.Name,
			Email:    params.Email,
			Role:     "user",
			IsGoogle: true,
			Token:    jwtToken.Token,
		}

		err = uc.repository.UserInsert(ctx, user)
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
	} else {
		paramsToken := jwt.DataToken{
			UserID: userByEmail.ID,
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
		resp = AuthLoginResponse{
			ExpiredAt: jwtToken.Exp,
			Token:     jwtToken.Token,
		}

		return resp, errData
	}
}

func (uc *authService) VerifyUserEmail(ctx context.Context, token string) (resp AuthRegisterResponse, errData exception.Error) {
	tokenValidate, err := jwtValidate.ValidateToken(token)
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

	userID := claims.Data.UserID

	user, err := uc.repository.UserFindByID(ctx, userID)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	// Update user's IsVerify field to true
	user.IsVerify = true

	err = uc.repository.UserUpdateVerifyStatus(ctx, user)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	resp = AuthRegisterResponse{
		ExpiredAt: claims.Iat,
		Token:     token,
	}

	return resp, errData
}
