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

func (tm *TransactionManager) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
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
	slog.InfoContext(ctx, "Transaction is set to context")
	if err := fn(ctxWithTx); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			slog.ErrorContext(ctx, "Failed to rollback transaction", "error", rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		slog.ErrorContext(ctx, "Failed to commit transaction", "error", err)
		return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to commit transaction")
	}

	return nil
}
