package repository

import (
	"context"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
)

type TagRepository interface {
	FindOrCreateByName(ctx context.Context, tag *entity.Tag) (*entity.Tag, error)
}
