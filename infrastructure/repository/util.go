package repository

import (
	"context"
	"database/sql"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// 値が nil または型のゼロ値の場合、NULL として扱う
func ToNullable[T comparable, N any](value *T, isZero func(T) bool, toNull func(T) N) N {
	if value == nil || isZero(*value) {
		var zero N
		return zero
	}
	return toNull(*value)
}

func GetExecDB(ctx context.Context, db *sql.DB) boil.ContextExecutor {
	var execDB boil.ContextExecutor
	execDB = db
	if contexDB, ok := ctx.Value(domaincontext.TransactionDB).(boil.ContextExecutor); ok {
		execDB = contexDB
	}
	return execDB
}
