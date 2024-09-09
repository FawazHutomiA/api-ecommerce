package auth

import (
	"context"
	"encoding/json"

	"example/internal/module/user"
	"example/internal/service"
	"example/pkg/helper"

	userInput "example/internal/input"
	jwtValidate "example/pkg/jwt"

	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	authService service.AuthService
	userService service.UserService
}

func NewAuthHandler(authService service.AuthService, userService service.UserService) *AuthHandler {
	return &AuthHandler{authService, userService}
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

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var input userInput.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err.(validator.ValidationErrors))

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Account failed to register", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	checkEmailInput := userInput.CheckEmailInput{Email: input.Email}

	isEmailAvailable, _ := h.authService.CheckEmail(checkEmailInput)
	if input.Email == isEmailAvailable.Email {
		response := helper.APIResponse("Email is already taken.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := h.authService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Account failed to register", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := jwtValidate.GenerateToken(newUser.ID)

	if err != nil {
		response := helper.APIResponse("Account failed to register", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input userInput.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err.(validator.ValidationErrors))

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := h.authService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := jwtValidate.GenerateToken(loggedInUser.ID)

	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, token)

	response := helper.APIResponse("Successfully logged in", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) CheckEmailAvailability(c *gin.Context) {
	var input userInput.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err.(validator.ValidationErrors))

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.authService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) HandleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	c.Redirect(http.StatusTemporaryRedirect, url)
	log.Println(url, "urllll")
}

func (h *AuthHandler) HandleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != randomState {
		fmt.Println("State is not valid")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("could not get token \n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("could not create request \n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("could not read response body \n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	var userInfo userInput.GoogleUser
	if err := json.Unmarshal(content, &userInfo); err != nil {
		fmt.Println("could not unmarshal JSON \n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	userLogin, err := h.userService.GetOrSaveUser(userInfo)

	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	tokenLogin, err := jwtValidate.GenerateToken(userLogin.ID)

	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(userLogin, tokenLogin)

	response := helper.APIResponse("Successfully logged in", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}
