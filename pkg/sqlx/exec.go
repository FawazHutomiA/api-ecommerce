package sqlx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Exec
func Exec(sqlx *sqlx.DB, query string, args ...interface{}) (sql.Result, error) {
	result, err := sqlx.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ExecWithContext(sqlx *sql.Tx, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := sqlx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
