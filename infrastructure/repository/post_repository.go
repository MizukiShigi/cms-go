package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/MizukiShigi/cms-go/infrastructure/db/sqlboiler/models"
	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *entity.Post) error {
	now := time.Now()
	dbPost := &models.Post{
		ID:               post.ID.String(),
		Title:            post.Title.String(),
		Content:          post.Content.String(),
		CreatedAt:        now,
		UpdatedAt:        now,
		FirstPublishedAt: ToNullable(
			post.FirstPublishedAt,
			func(t time.Time) bool { return t.IsZero() },
			null.TimeFrom,
		),
		ContentUpdatedAt: ToNullable(
			post.ContentUpdatedAt,
			func(t time.Time) bool { return t.IsZero() },
			null.TimeFrom,
		),
	}

	if err := dbPost.Insert(ctx, r.db, boil.Infer()); err != nil {
		return err
	}

	return nil
}
