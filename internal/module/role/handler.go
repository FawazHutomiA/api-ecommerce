package role

import (
	"example/pkg/app"
	"example/pkg/constant"
	"example/pkg/helper"
	"example/pkg/response"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type RoleHandler struct {
	App         app.AppConfig
	RoleService RoleService
}

func NewRoleHandler(app app.AppConfig, roleService RoleService) *RoleHandler {
	return &RoleHandler{App: app, RoleService: roleService}
}

// ListPaginate handles paginated role listing with Gin context
func (handler *RoleHandler) ListPaginate(w http.ResponseWriter, r *http.Request) {
	// Init
	var resp response.Response
	ctx := r.Context()

	param := helper.PaginationParams{}
	param = param.GetPaginateParam(r)

	service, err := handler.RoleService.ListPaginate(ctx, param)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *RoleHandler) Detail(w http.ResponseWriter, r *http.Request) {
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

	service, err := handler.RoleService.Detail(ctx, id)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}
