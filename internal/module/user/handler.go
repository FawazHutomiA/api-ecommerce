package user

import (
	"example/pkg/app"
	"example/pkg/constant"
	"example/pkg/helper"
	"example/pkg/response"
	"example/pkg/validator"
	"net/http"

	"github.com/go-chi/chi"
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

func (handler *UserHandler) Detail(w http.ResponseWriter, r *http.Request) {
	// Init
	var resp response.Response
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")

	id, errs := uuid.Parse(idStr)
	if errs != nil {
		handler.App.Logger.Error(errs)
		resp = response.Error(response.StatusBadRequest, constant.StatusBadRequest, errs)
		resp.JSON(w)
		return
	}

	service, err := handler.UserService.Detail(ctx, id)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var req UserCreateRequest
	var resp response.Response
	ctx := r.Context()

	resp, errV := validator.ValidateRequest(r, &req)
	if errV != nil {
		resp.JSON(w)
		return
	}

	service, err := handler.UserService.Create(ctx, req)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Init
	var req UserUpdateRequest
	var resp response.Response
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")

	id, errs := uuid.Parse(idStr)
	if errs != nil {
		handler.App.Logger.Error(errs)
		resp = response.Error(response.StatusBadRequest, constant.StatusBadRequest, errs)
		resp.JSON(w)
		return
	}

	resp, errV := validator.ValidateRequest(r, &req)
	if errV != nil {
		resp.JSON(w)
		return
	}

	service, err := handler.UserService.Update(ctx, id, req)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Init
	var resp response.Response
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")

	id, errs := uuid.Parse(idStr)
	if errs != nil {
		handler.App.Logger.Error(errs)
		resp = response.Error(response.StatusBadRequest, constant.StatusBadRequest, errs)
		resp.JSON(w)
		return
	}

	service, err := handler.UserService.Delete(ctx, id)
	if err.Errors != nil {
		handler.App.Logger.Error(errs)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *UserHandler) UserMe(w http.ResponseWriter, r *http.Request) {
	// Init
	var resp response.Response
	ctx := r.Context()

	userId := ctx.Value("userID").(uuid.UUID)

	service, err := handler.UserService.Detail(ctx, userId)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}
