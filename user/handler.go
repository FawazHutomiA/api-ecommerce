package user

import (
	"context"
	"encoding/json"
	"example/auth"
	"example/helper"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserHandler struct {
	userService Service
	authService auth.Service
}

func NewUserHandler(userService Service, authService auth.Service) *UserHandler {
	return &UserHandler{userService, authService}
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

// var (
// 	googleOAuthConfig = &oauth2.Config{
// 		RedirectURL:  "http://localhost:8080/api/v1/callback", //URL ini harus sama dengan redirect_uri pada google developer console
// 		ClientID:     os.Getenv("CLIENT_ID"),
// 		ClientSecret: os.Getenv("CLIENT_SECRET"),
// 		// scopes API URL bisa didapatkan di https://developers.google.com/identity/protocols/oauth2/scopes
// 		Scopes:   []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
// 		Endpoint: google.Endpoint,
// 	}
// 	randomState = "random"
// )

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err.(validator.ValidationErrors))

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Account failed to register", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	checkEmailInput := CheckEmailInput{Email: input.Email}

	isEmailAvailable, _ := h.userService.CheckEmail(checkEmailInput)
	if input.Email == isEmailAvailable.Email {
		response := helper.APIResponse("Email is already taken.", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Account failed to register", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)

	if err != nil {
		response := helper.APIResponse("Account failed to register", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Login(c *gin.Context) {
	var input LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err.(validator.ValidationErrors))

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)

	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := FormatUser(loggedInUser, token)

	response := helper.APIResponse("Successfully logged in", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) CheckEmailAvailability(c *gin.Context) {
	var input CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err.(validator.ValidationErrors))

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
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

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(User)
	userID := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfully uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUserByJWT(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(User)

	formatter := FormatUser(currentUser, "")

	response := helper.APIResponse("Successfuly get user data", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) HandleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	c.Redirect(http.StatusTemporaryRedirect, url)
	log.Println(url, "urllll")
}

func (h *UserHandler) HandleCallback(c *gin.Context) {
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

	var userInfo GoogleUser
	if err := json.Unmarshal(content, &userInfo); err != nil {
		fmt.Println("could not unmarshal JSON \n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	emailInput := CheckEmailInput{Email: userInfo.Email}
	isEmailAvailable, err := h.userService.IsEmailAvailable(emailInput)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}

		response := helper.APIResponse("Email checking failed.", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if isEmailAvailable {
		registerInput := RegisterUserInput{
			Name:       userInfo.Name,
			Email:      userInfo.Email,
			Occupation: "",
			Password:   "",
			IsGoogle:   true,
		}
		newUser, err := h.userService.RegisterUser(registerInput)
		if err != nil {
			response := helper.APIResponse("Account failed to register", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		tokenRegister, err := h.authService.GenerateToken(newUser.ID)

		if err != nil {
			response := helper.APIResponse("Account failed to register", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		formatter := FormatUser(newUser, tokenRegister)

		response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

		c.JSON(http.StatusOK, response)
	} else {
		getUserByEmail, err := h.userService.CheckEmail(emailInput)
		if err != nil {
			response := helper.APIResponse("Email checking failed.", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		loginInput := LoginInput{
			Email:    getUserByEmail.Email,
			Password: "",
		}

		loggedInUser, err := h.userService.Login(loginInput)

		if err != nil {
			errorMessage := gin.H{"errors": err.Error()}

			response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
			c.JSON(http.StatusUnprocessableEntity, response)
			return
		}

		tokenLogin, err := h.authService.GenerateToken(loggedInUser.ID)

		if err != nil {
			response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		formatter := FormatUser(loggedInUser, tokenLogin)

		response := helper.APIResponse("Successfully logged in", http.StatusOK, "success", formatter)

		c.JSON(http.StatusOK, response)
	}

}
