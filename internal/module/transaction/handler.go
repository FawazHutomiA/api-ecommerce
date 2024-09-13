package transaction

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

type TransactionHandler struct {
	App                app.AppConfig
	TransactionService TransactionService
}

func NewTransactionHandler(app app.AppConfig, transactionService TransactionService) *TransactionHandler {
	return &TransactionHandler{App: app, TransactionService: transactionService}
}

func (handler *TransactionHandler) ListPaginate(w http.ResponseWriter, r *http.Request) {
	// Init
	var resp response.Response
	ctx := r.Context()

	param := helper.PaginationParams{}
	param = param.GetPaginateParam(r)

	service, err := handler.TransactionService.ListPaginate(ctx, param)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *TransactionHandler) Detail(w http.ResponseWriter, r *http.Request) {
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

	service, err := handler.TransactionService.Detail(ctx, id)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *TransactionHandler) DetailByProductID(w http.ResponseWriter, r *http.Request) {
	// Init
	var resp response.Response
	ctx := r.Context()

	productIDStr := chi.URLParam(r, "product_id")

	id, errs := uuid.Parse(productIDStr)
	if errs != nil {
		handler.App.Logger.Error(errs)
		resp = response.Error(response.StatusBadRequest, constant.StatusBadRequest, errs)
		resp.JSON(w)
		return
	}

	service, err := handler.TransactionService.DetailByProductID(ctx, id)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *TransactionHandler) DetailByUserID(w http.ResponseWriter, r *http.Request) {
	// Init
	var resp response.Response
	ctx := r.Context()

	userIDStr := chi.URLParam(r, "user_id")

	id, errs := uuid.Parse(userIDStr)
	if errs != nil {
		handler.App.Logger.Error(errs)
		resp = response.Error(response.StatusBadRequest, constant.StatusBadRequest, errs)
		resp.JSON(w)
		return
	}

	service, err := handler.TransactionService.DetailByUserID(ctx, id)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var req TransactionCreateRequest
	var resp response.Response
	ctx := r.Context()

	userID := ctx.Value("userID").(uuid.UUID)

	resp, errV := validator.ValidateRequest(r, &req)
	if errV != nil {
		resp.JSON(w)
		return
	}

	service, err := handler.TransactionService.Create(ctx, req, userID)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Init
	var req TransactionUpdateRequest
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

	service, err := handler.TransactionService.Update(ctx, id, req)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	service, err := handler.TransactionService.Delete(ctx, id)
	if err.Errors != nil {
		handler.App.Logger.Error(errs)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}

func (handler *TransactionHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	// Init
	var req TransactionNotificationInput
	var resp response.Response
	ctx := r.Context()

	resp, errV := validator.ValidateRequest(r, &req)
	if errV != nil {
		resp.JSON(w)
		return
	}

	payment, err := handler.TransactionService.ProcessPayment(ctx, req)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", payment)
	resp.JSON(w)
}
