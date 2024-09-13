package product

import (
	"context"
	"database/sql"

	"example/internal/entity"
	"example/internal/repository/postgresql/product"
	"example/pkg/app"
	"example/pkg/exception"
	"example/pkg/helper"
	"example/pkg/response"
	"example/pkg/strings"

	"github.com/google/uuid"
)

type ProductService interface {
	ListPaginate(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, errData exception.Error)
	Detail(ctx context.Context, id uuid.UUID) (resp ProductDetailResponse, errData exception.Error)
	DetailBySlug(ctx context.Context, slug string) (resp ProductDetailResponse, errData exception.Error)
	Create(ctx context.Context, params ProductCreateRequest) (resp ProductCreateResponse, errData exception.Error)
	Update(ctx context.Context, id uuid.UUID, params ProductUpdateRequest) (resp ProductUpdateResponse, errData exception.Error)
	Delete(ctx context.Context, id uuid.UUID) (resp bool, errData exception.Error)
}

type productService struct {
	app        app.AppConfig
	repository product.ProductRepository
}

func NewProductService(app app.AppConfig, repository product.ProductRepository) ProductService {
	return &productService{
		app:        app,
		repository: repository,
	}
}

func (uc *productService) ListPaginate(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, errData exception.Error) {
	// repository
	productRepo, err := uc.repository.ProductFindAll(ctx, params)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	return productRepo, errData
}

func (uc *productService) Detail(ctx context.Context, id uuid.UUID) (resp ProductDetailResponse, errData exception.Error) {
	// repository
	productRepo, err := uc.repository.ProductFindByID(ctx, id)

	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return resp, exception.Error{
			Status:  response.StatusNotFound,
			Message: "Product  Not Found",
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
	resp = ProductDetailResponse(productRepo)

	return resp, errData
}

func (uc *productService) DetailBySlug(ctx context.Context, slug string) (resp ProductDetailResponse, errData exception.Error) {
	// repository
	productRepo, err := uc.repository.ProductFindBySlug(ctx, slug)

	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return resp, exception.Error{
			Status:  response.StatusNotFound,
			Message: "Product  Not Found",
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
	resp = ProductDetailResponse(productRepo)

	return resp, errData
}

// Create usecase
func (uc *productService) Create(ctx context.Context, params ProductCreateRequest) (resp ProductCreateResponse, errData exception.Error) {
	// repository
	_, err := uc.repository.ProductFindByName(ctx, params.Name)
	switch err {
	case nil:
		return resp, exception.Error{
			Status:  response.StatusConflicted,
			Message: "This Product  Already Exsist",
			Errors:  exception.ErrConflicted,
		}
	case sql.ErrNoRows:
		err = nil
	default:
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	productID := uuid.New()

	// map insert
	product := entity.Product{
		ID:               productID,
		UserID:           params.UserID,
		Name:             params.Name,
		ShortDescription: params.ShortDescription,
		Description:      params.Description,
		GoalAmount:       params.GoalAmount,
		CurrentAmount:    params.CurrentAmount,
		Slug:             strings.Slug(params.Name),
		BackerAmount:     params.BackerAmount,
	}

	// save to db
	err = uc.repository.ProductInsert(ctx, product)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	params.Slug = strings.Slug(params.Name)

	// map response
	resp = ProductCreateResponse(params)

	return resp, errData
}

func (uc *productService) Update(ctx context.Context, id uuid.UUID, params ProductUpdateRequest) (resp ProductUpdateResponse, errData exception.Error) {
	_, errData = uc.Detail(ctx, id)
	if errData.Errors != nil {
		return resp, exception.Error(errData)
	}

	// map insert
	product := entity.Product{
		ID:               id,
		UserID:           params.UserID,
		Name:             params.Name,
		ShortDescription: params.ShortDescription,
		Description:      params.Description,
		GoalAmount:       params.GoalAmount,
		CurrentAmount:    params.CurrentAmount,
		Slug:             strings.Slug(params.Name),
		BackerAmount:     params.BackerAmount,
	}

	err := uc.repository.ProductUpdateByID(ctx, product)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Error",
			Errors:  exception.ErrBadRequest,
		}
	}

	params.Slug = strings.Slug(params.Name)

	// map response
	resp = ProductUpdateResponse(params)

	return resp, errData
}

func (uc *productService) Delete(ctx context.Context, id uuid.UUID) (resp bool, errData exception.Error) {
	_, errData = uc.Detail(ctx, id)
	if errData.Errors != nil {
		return resp, exception.Error(errData)
	}

	// update to db
	err := uc.repository.ProductDelete(ctx, id)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	return resp, errData
}
