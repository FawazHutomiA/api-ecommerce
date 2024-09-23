package auth

import (
	"example/pkg/app"
	"example/pkg/response"
	"example/pkg/validator"
	"net/http"
)

type AuthHandler struct {
	App         app.AppConfig
	AuthService AuthService
}

func NewAuthHandler(app app.AppConfig, authService AuthService) *AuthHandler {
	return &AuthHandler{App: app, AuthService: authService}
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Init
	var req AuthRegisterRequest
	var resp response.Response
	ctx := r.Context()

	resp, errV := validator.ValidateRequest(r, &req)
	if errV != nil {
		resp.JSON(w)
		return
	}

	service, err := handler.AuthService.Register(ctx, req)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Register Success", service)
	resp.JSON(w)
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Init
	var req AuthLoginRequest
	var resp response.Response
	ctx := r.Context()

	resp, errV := validator.ValidateRequest(r, &req)
	if errV != nil {
		resp.JSON(w)
		return
	}

	service, err := handler.AuthService.Login(ctx, req)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Login Success", service)
	resp.JSON(w)
}
