package transaction

import (
	"context"
	"example/internal/entity"
	"example/pkg/app"
	"example/pkg/helper"
	"example/pkg/sqlx"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	TransactionFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error)
	TransactionFindByID(ctx context.Context, id uuid.UUID) (resp entity.Transaction, err error)
	TransactionFindByProductID(ctx context.Context, productID uuid.UUID) (resp entity.Transaction, err error)
	TransactionFindByUserID(ctx context.Context, userID uuid.UUID) (resp entity.Transaction, err error)
	TransactionInsert(ctx context.Context, transaction entity.Transaction) (err error)
	TransactionUpdateByID(ctx context.Context, transaction entity.Transaction) (err error)
	TransactionDelete(ctx context.Context, id uuid.UUID) (err error)
}

type transactionRepositoryImpl struct {
	app app.AppConfig
}

func NewTransactionRepository(app app.AppConfig) TransactionRepository {
	return &transactionRepositoryImpl{
		app: app,
	}
}

func (repo *transactionRepositoryImpl) TransactionFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error) {
	query := FIND_ALL

	if params.Search != "" {
		escapedSearch := strings.Replace(params.Search, "'", "''", -1)
		addFilter := fmt.Sprintf("AND r.name ILIKE '%%%s%%'", escapedSearch)
		query = fmt.Sprintf(`%s %s`, query, addFilter)
	}

	if params.OrderBy != "" {
		params.SortBy = "r.created_at"
		query = fmt.Sprintf("%s ORDER BY %s %s", query, params.SortBy, params.OrderBy)
	}

	var transaction []entity.Transaction

	pagination := sqlx.NewPaginationMetadata(repo.app.Db)
	result, err := pagination.GetPagination(query, params, &transaction)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return result, nil
}

func (repo *transactionRepositoryImpl) TransactionFindByID(ctx context.Context, id uuid.UUID) (resp entity.Transaction, err error) {
	err = repo.app.Db.Get(&resp, FIND_BY_ID, id)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *transactionRepositoryImpl) TransactionFindByProductID(ctx context.Context, productID uuid.UUID) (resp entity.Transaction, err error) {
	err = repo.app.Db.Get(&resp, FIND_BY_PRODUCT_ID, productID)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *transactionRepositoryImpl) TransactionFindByUserID(ctx context.Context, userID uuid.UUID) (resp entity.Transaction, err error) {
	err = repo.app.Db.Get(&resp, FIND_BY_USER_ID, userID)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *transactionRepositoryImpl) TransactionInsert(ctx context.Context, transaction entity.Transaction) (err error) {
	_, err = repo.app.Db.ExecContext(ctx, INSERT, transaction.ToInsert()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *transactionRepositoryImpl) TransactionUpdateByID(ctx context.Context, transaction entity.Transaction) (err error) {
	_, err = repo.app.Db.ExecContext(ctx, UPDATE_BY_ID, transaction.ToUpdate()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *transactionRepositoryImpl) TransactionDelete(ctx context.Context, id uuid.UUID) (err error) {
	_, err = repo.app.Db.ExecContext(ctx, DELETE_BY_ID, id)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}
