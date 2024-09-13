package product

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

type ProductRepository interface {
	ProductFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error)
	ProductFindByID(ctx context.Context, id uuid.UUID) (resp entity.Product, err error)
	ProductFindByName(ctx context.Context, name string) (resp entity.Product, err error)
	ProductFindBySlug(ctx context.Context, slug string) (resp entity.Product, err error)
	ProductInsert(ctx context.Context, product entity.Product) (err error)
	ProductUpdateByID(ctx context.Context, product entity.Product) (err error)
	ProductDelete(ctx context.Context, id uuid.UUID) (err error)
}

type productRepositoryImpl struct {
	app app.AppConfig
}

func NewProductRepository(app app.AppConfig) ProductRepository {
	return &productRepositoryImpl{
		app: app,
	}
}

func (repo *productRepositoryImpl) ProductFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error) {
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

	var product []entity.Product

	pagination := sqlx.NewPaginationMetadata(repo.app.Db)
	result, err := pagination.GetPagination(query, params, &product)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return result, nil
}

func (repo *productRepositoryImpl) ProductFindByID(ctx context.Context, id uuid.UUID) (resp entity.Product, err error) {
	err = repo.app.Db.Get(&resp, FIND_BY_ID, id)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *productRepositoryImpl) ProductFindByName(ctx context.Context, name string) (resp entity.Product, err error) {
	err = repo.app.Db.GetContext(ctx, &resp, FIND_BY_NAME, name)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *productRepositoryImpl) ProductFindBySlug(ctx context.Context, slug string) (resp entity.Product, err error) {
	err = repo.app.Db.GetContext(ctx, &resp, FIND_BY_SLUG, slug)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *productRepositoryImpl) ProductInsert(ctx context.Context, product entity.Product) (err error) {
	_, err = repo.app.Db.ExecContext(ctx, INSERT, product.ToInsert()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *productRepositoryImpl) ProductUpdateByID(ctx context.Context, product entity.Product) (err error) {
	_, err = repo.app.Db.ExecContext(ctx, UPDATE_BY_ID, product.ToUpdate()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *productRepositoryImpl) ProductDelete(ctx context.Context, id uuid.UUID) (err error) {
	_, err = repo.app.Db.ExecContext(ctx, DELETE_BY_ID, id)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}
