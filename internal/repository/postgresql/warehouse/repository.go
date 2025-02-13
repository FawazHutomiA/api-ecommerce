package warehouse

import (
	"example/internal/entity"
	"example/pkg/app"
	"example/pkg/helper"
	sqlxPkg "example/pkg/sqlx"

	"context"
	"database/sql"

	"github.com/google/uuid"
)

type WarehouseRepository interface {
	WarehouseFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error)
	WarehouseFindByID(ctx context.Context, id uuid.UUID) (resp entity.Warehouse, err error)
	WarehouseInsert(ctx context.Context, tx *sql.Tx, warehouse entity.Warehouse) (err error)
	WarehouseUpdate(ctx context.Context, warehouse entity.Warehouse) (err error)
	WarehouseDelete(ctx context.Context, tx *sql.Tx, id uuid.UUID) (err error)
}

type warehouseRepository struct {
	app app.AppConfig
}

func NewWarehouseRepository(app app.AppConfig) WarehouseRepository {
	return &warehouseRepository{
		app: app,
	}
}

func (repo *warehouseRepository) WarehouseFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error) {
	query := FIND_ALL

	var warehouse []entity.Warehouse
	pagination := sqlxPkg.NewPaginationMetadata(repo.app.Db)
	result, err := pagination.GetPagination(query, params, &warehouse)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return result, nil
}

func (repo *warehouseRepository) WarehouseFindByID(ctx context.Context, id uuid.UUID) (resp entity.Warehouse, err error) {
	err = repo.app.Db.GetContext(ctx, &resp, FIND_BY_ID, id)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *warehouseRepository) WarehouseInsert(ctx context.Context, tx *sql.Tx, warehouse entity.Warehouse) (err error) {
	_, err = tx.ExecContext(ctx, INSERT, warehouse.ToInsert()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *warehouseRepository) WarehouseUpdate(ctx context.Context, warehouse entity.Warehouse) (err error) {
	_, err = repo.app.Db.ExecContext(ctx, UPDATE_BY_ID, warehouse.ToUpdate()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *warehouseRepository) WarehouseDelete(ctx context.Context, tx *sql.Tx, id uuid.UUID) (err error) {
	_, err = tx.ExecContext(ctx, DELETE_BY_ID, id)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}
