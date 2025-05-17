package repository

import (
	"context"
	"database/sql"
	"log/slog"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type TransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (r *PostRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to begin transaction")
	}

	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback()
			if err != nil {
				slog.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
			}
			panic(p)
		}
	}()

	ctxWithTx := context.WithValue(ctx, domaincontext.TransactionDB, tx)

	if err := fn(ctxWithTx); err != nil {
		err := tx.Rollback()
		if err != nil {
			slog.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to commit transaction")
	}

	return nil
}
