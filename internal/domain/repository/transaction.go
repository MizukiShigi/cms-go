package repository

import "context"

type TransactionManager interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
