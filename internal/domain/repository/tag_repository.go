package repository

import (
	"context"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type TagRepository interface {
	FindOrCreateByName(ctx context.Context, tagName valueobject.TagName) (*entity.Tag, error)
}
