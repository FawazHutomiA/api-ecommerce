package sqlx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Transaction
func BeginTx(sqlx *sqlx.DB, ctx context.Context) (*sql.Tx, error) {
	tx, err := sqlx.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func Commit(sqlx *sql.Tx, ctx context.Context) error {
	err := sqlx.Commit()
	if err != nil {
		sqlx.Rollback()
		return err
	}
	return nil
}

func CommitOrRollback(tx *sql.Tx) error {
	err := recover()
	if err != nil {
		errorRollback := tx.Rollback()
		if errorRollback != nil {
			return errorRollback
		}
	} else {
		errorCommit := tx.Commit()
		if errorCommit != nil {
			return errorCommit
		}
	}
	return nil
}
