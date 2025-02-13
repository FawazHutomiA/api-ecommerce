package warehouse

import (
	"context"
	"database/sql"
	"example/internal/repository/postgresql/token"
	"example/internal/repository/postgresql/warehouse"
	"example/pkg/app"
	"example/pkg/exception"
	"example/pkg/helper"
	"example/pkg/response"

	"github.com/google/uuid"
)

type WarehouseService interface {
	ListPaginate(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, errData exception.Error)
	Detail(ctx context.Context, id uuid.UUID) (resp WarehouseDetailResponse, errData exception.Error)
}

type warehouseService struct {
	app        app.AppConfig
	repository warehouse.WarehouseRepository
	tokenRepo  token.TokenRepository
}

func NewWarehouseService(app app.AppConfig, repository warehouse.WarehouseRepository, tokenRepo token.TokenRepository) WarehouseService {
	return &warehouseService{
		app:        app,
		repository: repository,
		tokenRepo:  tokenRepo,
	}
}

func (uc *warehouseService) ListPaginate(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, errData exception.Error) {
	// repository
	repo, err := uc.repository.WarehouseFindAll(ctx, params)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	return repo, errData
}

func (uc *warehouseService) Detail(ctx context.Context, id uuid.UUID) (resp WarehouseDetailResponse, errData exception.Error) {
	// repository
	warehouseRepo, err := uc.repository.WarehouseFindByID(ctx, id)
	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return resp, exception.Error{
			Status:  response.StatusNotFound,
			Message: "Warehouse Not Found",
			Errors:  exception.ErrNotFound,
		}
	default:
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	// map response
	resp = WarehouseDetailResponse(warehouseRepo)

	return resp, errData
}
