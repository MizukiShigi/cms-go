package repository

import (
	"context"
	"database/sql"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func FindOrCreateByName(ctx context.Context, tagName valueobject.TagName) (*entity.Tag, error) {
	
}
