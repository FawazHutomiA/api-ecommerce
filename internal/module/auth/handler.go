package auth

import (
	"context"
	"encoding/json"
	"example/pkg/app"
	"example/pkg/helper"
	"example/pkg/response"
	"example/pkg/validator"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	App         app.AppConfig
	AuthService AuthService
}

func NewAuthHandler(app app.AppConfig, authService AuthService) *AuthHandler {
	return &AuthHandler{App: app, AuthService: authService}
}

var (
	googleOauthConfig *oauth2.Config
	randomState       = "random"
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/v1/callback",
		ClientID:     helper.GetENV("CLIENT_ID"),
		ClientSecret: helper.GetENV("CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
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

func (handler *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	log.Println(url, "urllll")
}

func (handler *AuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	ctx := r.Context()
	// Init
	state := r.URL.Query().Get("state")
	if state != randomState {
		handler.App.Logger.Error("State is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		handler.App.Logger.Error("Could not get token", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	responseGoogle, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		handler.App.Logger.Error("Could not create request", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer responseGoogle.Body.Close()

	content, err := io.ReadAll(responseGoogle.Body)
	if err != nil {
		handler.App.Logger.Error("Could not read response body", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var userInfo GoogleUser
	if err := json.Unmarshal(content, &userInfo); err != nil {
		handler.App.Logger.Error("Could not unmarshal JSON", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	userLogin, errV := handler.AuthService.GetOrSaveUser(ctx, userInfo)
	if errV.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(errV.Status, errV.Message, errV.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Login Success", userLogin)
	resp.JSON(w)
}
