package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/MizukiShigi/cms-go/infrastructure/db/sqlboiler/models"
	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (tr *TagRepository) FindOrCreateByName(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	isInsert := false
	execDB := GetExecDB(ctx, tr.db)
	dbTag, err := models.Tags(
		qm.Where("name = ?", tag.Name.String()),
	).One(ctx, execDB)
	if err != nil {
		if err != sql.ErrNoRows {
			errMsg := "Failed to find tag"
			ctx := domaincontext.WithValue(ctx, "error", err.Error())
			slog.ErrorContext(ctx, errMsg)
			return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, errMsg)
		} else {
			isInsert = true
		}
	}

	// 既に存在する場合は登録不要のため、そのまま返す
	if !isInsert {
		tagID, err := valueobject.ParseTagID(dbTag.ID)
		if err != nil {
			errMsg := "Failed to parse tag ID"
			ctx := domaincontext.WithValue(ctx, "error", err.Error())
			slog.ErrorContext(ctx, errMsg)
			return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, errMsg)
		}

		tagName, err := valueobject.NewTagName(dbTag.Name)
		if err != nil {
			errMsg := "Failed to parse tag name"
			ctx := domaincontext.WithValue(ctx, "error", err.Error())
			slog.ErrorContext(ctx, errMsg)
			return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, errMsg)
		}

		return &entity.Tag{
			ID:        tagID,
			Name:      tagName,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}, nil
	}

	dbTag = &models.Tag{
		ID:   valueobject.NewTagID().String(),
		Name: tag.Name.String(),
	}
	if err := dbTag.Insert(ctx, execDB, boil.Infer()); err != nil {
		errMsg := "Failed to create tag"
		ctx := domaincontext.WithValue(ctx, "error", err.Error())
		slog.ErrorContext(ctx, errMsg)
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, errMsg)
	}

	tagID, err := valueobject.ParseTagID(dbTag.ID)
	if err != nil {
		errMsg := "Failed to parse tag ID"
		ctx := domaincontext.WithValue(ctx, "error", err.Error())
		slog.ErrorContext(ctx, errMsg)
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, errMsg)
	}

	tagName, err := valueobject.NewTagName(dbTag.Name)
	if err != nil {
		errMsg := "Failed to parse tag name"
		ctx := domaincontext.WithValue(ctx, "error", err.Error())
		slog.ErrorContext(ctx, errMsg)
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, errMsg)
	}

	return &entity.Tag{
		ID:        tagID,
		Name:      tagName,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}, nil
}
