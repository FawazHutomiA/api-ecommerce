package token

import (
	"context"
	"database/sql"
	"example/internal/entity"
	"example/pkg/app"
	"example/pkg/helper"
	sqlxPkg "example/pkg/sqlx"

	"github.com/google/uuid"
)

type TokenRepository interface {
	TokenFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error)
	TokenFindByID(ctx context.Context, id uuid.UUID) (resp entity.Token, err error)
	TokenFindByUserID(ctx context.Context, userID uuid.UUID) (resp entity.Token, err error)
	TokenInsert(ctx context.Context, tx *sql.Tx, token entity.Token) (err error)
	TokenUpdate(ctx context.Context, token entity.Token) (err error)
	TokenDelete(ctx context.Context, tx *sql.Tx, id uuid.UUID) (err error)
	TokenDeleteByUserID(ctx context.Context, tx *sql.Tx, userID uuid.UUID) (err error)
}

type tokenRepository struct {
	app app.AppConfig
}

func NewTokenRepository(app app.AppConfig) TokenRepository {
	return &tokenRepository{
		app: app,
	}
}

func (repo *tokenRepository) TokenFindAll(ctx context.Context, params helper.PaginationParams) (resp helper.Pagination, err error) {
	query := FIND_ALL

	var token []entity.Token
	pagination := sqlxPkg.NewPaginationMetadata(repo.app.Db)
	result, err := pagination.GetPagination(query, params, &token)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return result, nil
}

func (repo *tokenRepository) TokenFindByID(ctx context.Context, id uuid.UUID) (resp entity.Token, err error) {
	err = repo.app.Db.GetContext(ctx, &resp, FIND_BY_ID, id)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *tokenRepository) TokenFindByUserID(ctx context.Context, userID uuid.UUID) (resp entity.Token, err error) {
	err = repo.app.Db.GetContext(ctx, &resp, FIND_BY_USER_ID, userID)
	if err != nil {
		repo.app.Logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func (repo *tokenRepository) TokenInsert(ctx context.Context, tx *sql.Tx, token entity.Token) (err error) {
	_, err = tx.ExecContext(ctx, INSERT, token.ToInsert()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *tokenRepository) TokenUpdate(ctx context.Context, token entity.Token) (err error) {
	_, err = repo.app.Db.ExecContext(ctx, UPDATE_BY_ID, token.ToUpdate()...)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *tokenRepository) TokenDelete(ctx context.Context, tx *sql.Tx, id uuid.UUID) (err error) {
	_, err = tx.ExecContext(ctx, DELETE_BY_ID, id)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}

func (repo *tokenRepository) TokenDeleteByUserID(ctx context.Context, tx *sql.Tx, userID uuid.UUID) (err error) {
	_, err = tx.ExecContext(ctx, DELETE_BY_USER_ID, userID)
	if err != nil {
		repo.app.Logger.Error(err)
		return err
	}
	return nil
}
