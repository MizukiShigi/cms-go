package repository

import (
	"context"
	"database/sql"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
)

type ImageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) Create(ctx context.Context, image *entity.Image) error {
	return nil
}
