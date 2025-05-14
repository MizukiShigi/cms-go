package repository

import (
	"context"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
)

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) error
	Get(ctx context.Context, id string) (*entity.Post, error)
}
