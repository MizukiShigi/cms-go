package repository

import (
	"context"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
)

type ImageRepository interface {
	Create(ctx context.Context, image *entity.Image) error
}
