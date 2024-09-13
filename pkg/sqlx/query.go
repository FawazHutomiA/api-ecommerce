package sqlx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Query
func Query(sqlx *sqlx.DB, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := sqlx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func QueryWithContext(sqlx *sqlx.DB, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := sqlx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// QueryRow
func QueryRow(sqlx *sqlx.DB, query string, args ...interface{}) (*sql.Row, error) {
	return sqlx.QueryRow(query, args...), nil
}

func QueryRowWithContext(sqlx *sqlx.DB, ctx context.Context, query string, args ...interface{}) (*sql.Row, error) {
	return sqlx.QueryRowContext(ctx, query, args...), nil
}

// Select
func Select(sqlx *sqlx.DB, dest interface{}, query string, args ...interface{}) error {
	err := sqlx.Select(dest, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func SelectWithContext(sqlx *sqlx.DB, ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := sqlx.SelectContext(ctx, dest, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Get
func Get(sqlx *sqlx.DB, dest interface{}, query string, args ...interface{}) error {
	err := sqlx.Get(dest, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func GetWithContext(sqlx *sqlx.DB, ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := sqlx.GetContext(ctx, dest, query, args...)
	if err != nil {
		return err
	}
	return nil
}
