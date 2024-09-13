package user

import (
	"example/pkg/app"
	"example/pkg/helper"
	"example/pkg/response"
	"net/http"

	"github.com/google/uuid"
)

type UserHandler struct {
	App         app.AppConfig
	UserService UserService
}

func NewUserHandler(app app.AppConfig, userService UserService) *UserHandler {
	return &UserHandler{App: app, UserService: userService}
}

// ListPaginate handles paginated user listing with Gin context
func (handler *UserHandler) ListPaginate(w http.ResponseWriter, r *http.Request) {
	// Init
	var resp response.Response
	ctx := r.Context()

	param := helper.PaginationParams{}
	param = param.GetPaginateParam(r)

	service, err := handler.UserService.ListPaginate(ctx, param)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *UserHandler) DetailByJwt(w http.ResponseWriter, r *http.Request) {
	// Init
	var resp response.Response
	ctx := r.Context()

	userID := ctx.Value("userID").(uuid.UUID)

	service, err := handler.UserService.Detail(ctx, userID)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}
