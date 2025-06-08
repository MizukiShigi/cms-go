package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/MizukiShigi/cms-go/infrastructure/db/sqlboiler/models"
	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ImageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) Create(ctx context.Context, image *entity.Image) error {
	now := time.Now()
	dbImage := &models.Image{
		ID:               image.ID.String(),
		OriginalFilename: image.OriginalFilename.String(),
		StoredFilename:   image.StoredFilename,
		GCSURL:           image.GCSURL,
		PostID:           image.PostID.String(),
		UserID:           image.UserID.String(),
		SortOrder:        null.IntFrom(image.SortOrder),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := dbImage.Insert(ctx, GetExecDB(ctx, r.db), boil.Infer()); err != nil {
		return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to create image")
	}

	return nil
}
