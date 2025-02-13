package warehouse

import (
	"example/pkg/app"
	"example/pkg/constant"
	"example/pkg/helper"
	"example/pkg/response"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type WarehouseHandler struct {
	App              app.AppConfig
	WarehouseService WarehouseService
}

func NewWarehouseHandler(app app.AppConfig, warehouseService WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{App: app, WarehouseService: warehouseService}
}

// ListPaginate handles paginated warehouse listing with Gin context
func (handler *WarehouseHandler) ListPaginate(w http.ResponseWriter, r *http.Request) {
	// Init
	var resp response.Response
	ctx := r.Context()

	param := helper.PaginationParams{}
	param = param.GetPaginateParam(r)

	service, err := handler.WarehouseService.ListPaginate(ctx, param)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *WarehouseHandler) Detail(w http.ResponseWriter, r *http.Request) {
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

	service, err := handler.WarehouseService.Detail(ctx, id)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}
