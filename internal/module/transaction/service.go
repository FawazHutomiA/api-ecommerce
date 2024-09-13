package transaction

import (
	"context"
	"database/sql"

	"example/internal/entity"
	"example/internal/module/payment"
	"example/internal/repository/postgresql/product"
	"example/internal/repository/postgresql/transaction"
	"example/internal/repository/postgresql/user"
	"example/pkg/app"
	"example/pkg/exception"
	"example/pkg/helper"
	"example/pkg/response"

	"github.com/google/uuid"
)

type TransactionService interface {
	ListPaginate(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, errData exception.Error)
	Detail(ctx context.Context, id uuid.UUID) (resp TransactionDetailResponse, errData exception.Error)
	DetailByProductID(ctx context.Context, productID uuid.UUID) (resp TransactionDetailResponse, errData exception.Error)
	DetailByUserID(ctx context.Context, userID uuid.UUID) (resp TransactionDetailResponse, errData exception.Error)
	Create(ctx context.Context, params TransactionCreateRequest, userID uuid.UUID) (resp TransactionCreateResponse, errData exception.Error)
	Update(ctx context.Context, id uuid.UUID, params TransactionUpdateRequest) (resp TransactionUpdateResponse, errData exception.Error)
	Delete(ctx context.Context, id uuid.UUID) (resp bool, errData exception.Error)
	ProcessPayment(ctx context.Context, params TransactionNotificationInput) (resp bool, errData exception.Error)
}

type transactionService struct {
	app            app.AppConfig
	repository     transaction.TransactionRepository
	productRepo    product.ProductRepository
	userRepo       user.UserRepository
	paymentService payment.PaymentService
}

func NewTransactionService(app app.AppConfig, repository transaction.TransactionRepository, productRepo product.ProductRepository, userRepo user.UserRepository, paymentService payment.PaymentService) TransactionService {
	return &transactionService{
		app:            app,
		repository:     repository,
		productRepo:    productRepo,
		userRepo:       userRepo,
		paymentService: paymentService,
	}
}

func (uc *transactionService) ListPaginate(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, errData exception.Error) {
	// repository
	transactionRepo, err := uc.repository.TransactionFindAll(ctx, params)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	return transactionRepo, errData
}

func (uc *transactionService) Detail(ctx context.Context, id uuid.UUID) (resp TransactionDetailResponse, errData exception.Error) {
	// repository
	transactionRepo, err := uc.repository.TransactionFindByID(ctx, id)
	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return resp, exception.Error{
			Status:  response.StatusNotFound,
			Message: "Transaction  Not Found",
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
	resp = TransactionDetailResponse(transactionRepo)

	return resp, errData
}

func (uc *transactionService) DetailByProductID(ctx context.Context, productID uuid.UUID) (resp TransactionDetailResponse, errData exception.Error) {
	// repository
	transactionRepo, err := uc.repository.TransactionFindByProductID(ctx, productID)
	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return resp, exception.Error{
			Status:  response.StatusNotFound,
			Message: "Transaction  Not Found",
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
	resp = TransactionDetailResponse(transactionRepo)

	return resp, errData
}

func (uc *transactionService) DetailByUserID(ctx context.Context, userID uuid.UUID) (resp TransactionDetailResponse, errData exception.Error) {
	// repository
	transactionRepo, err := uc.repository.TransactionFindByUserID(ctx, userID)
	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return resp, exception.Error{
			Status:  response.StatusNotFound,
			Message: "Transaction  Not Found",
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
	resp = TransactionDetailResponse(transactionRepo)

	return resp, errData
}

// Create usecase
func (uc *transactionService) Create(ctx context.Context, params TransactionCreateRequest, userID uuid.UUID) (resp TransactionCreateResponse, errData exception.Error) {
	transactionID := uuid.New()

	user, err := uc.userRepo.UserFindByID(ctx, userID)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	paymentReq := entity.PaymentRequest{
		ID:     transactionID,
		Amount: params.Amount,
	}

	paymentURL, err := uc.paymentService.GetPaymentURL(paymentReq, user)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	// map insert
	transaction := entity.Transaction{
		ID:         transactionID,
		ProductID:  params.ProductID,
		UserID:     userID,
		Amount:     params.Amount,
		Status:     "pending",
		Code:       params.Code,
		PaymentURL: paymentURL,
	}

	// save to db
	err = uc.repository.TransactionInsert(ctx, transaction)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	// map response
	resp = TransactionCreateResponse{
		ProductID:  params.ProductID,
		UserID:     userID,
		Amount:     params.Amount,
		Status:     transaction.Status,
		Code:       params.Code,
		PaymentURL: paymentURL,
	}

	return resp, errData
}

func (uc *transactionService) Update(ctx context.Context, id uuid.UUID, params TransactionUpdateRequest) (resp TransactionUpdateResponse, errData exception.Error) {
	_, errData = uc.Detail(ctx, id)
	if errData.Errors != nil {
		return resp, exception.Error(errData)
	}

	// map insert
	transaction := entity.Transaction{
		ID:     id,
		UserID: params.UserID,
	}

	err := uc.repository.TransactionUpdateByID(ctx, transaction)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Error",
			Errors:  exception.ErrBadRequest,
		}
	}

	// map response
	resp = TransactionUpdateResponse(params)

	return resp, errData
}

func (uc *transactionService) Delete(ctx context.Context, id uuid.UUID) (resp bool, errData exception.Error) {
	_, errData = uc.Detail(ctx, id)
	if errData.Errors != nil {
		return resp, exception.Error(errData)
	}

	// update to db
	err := uc.repository.TransactionDelete(ctx, id)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Error",
			Errors:  exception.ErrInternalServer,
		}
	}

	return resp, errData
}

func (uc *transactionService) ProcessPayment(ctx context.Context, params TransactionNotificationInput) (resp bool, errData exception.Error) {
	transaction, err := uc.repository.TransactionFindByID(ctx, params.OrderID)
	switch err {
	case nil:
		err = nil
	case sql.ErrNoRows:
		return resp, exception.Error{
			Status:  response.StatusNotFound,
			Message: "Transaction  Not Found",
			Errors:  exception.ErrNotFound,
		}
	default:
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Something Wrong",
			Errors:  exception.ErrBadRequest,
		}
	}

	if params.PaymentType == "credit_card" && params.TransactionStatus == "capture" && params.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if params.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if params.TransactionStatus == "deny" || params.TransactionStatus == "expire" || params.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	err = uc.repository.TransactionUpdateByID(ctx, transaction)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Error",
			Errors:  exception.ErrBadRequest,
		}
	}

	product, err := uc.productRepo.ProductFindByID(ctx, transaction.ProductID)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Error",
			Errors:  exception.ErrBadRequest,
		}
	}

	if transaction.Status == "paid" {
		product.BackerAmount = product.BackerAmount + 1
		product.CurrentAmount = product.CurrentAmount + transaction.Amount

		err := uc.productRepo.ProductUpdateByID(ctx, product)
		if err != nil {
			return resp, exception.Error{
				Status:  response.StatusBadRequest,
				Message: "Error",
				Errors:  exception.ErrBadRequest,
			}
		}

	}

	return resp, errData
}
