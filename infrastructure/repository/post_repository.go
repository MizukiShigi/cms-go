package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/MizukiShigi/cms-go/infrastructure/db/sqlboiler/models"
	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
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
		ID:        post.ID.String(),
		Title:     post.Title.String(),
		Content:   post.Content.String(),
		UserID:    post.UserID.String(),
		CreatedAt: now,
		UpdatedAt: now,
		FirstPublishedAt: ToNullable(
			*post.FirstPublishedAt,
			func(t time.Time) bool { return t.IsZero() },
			null.TimeFrom,
		),
		ContentUpdatedAt: ToNullable(
			*post.ContentUpdatedAt,
			func(t time.Time) bool { return t.IsZero() },
			null.TimeFrom,
		),
	}

	if err := dbPost.Insert(ctx, r.db, boil.Infer()); err != nil {
		return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to create post")
	}

	return nil
}

func (r *PostRepository) Get(ctx context.Context, id string) (*entity.Post, error) {
	dbPost, err := models.FindPost(ctx, r.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, valueobject.NewMyError(valueobject.NotFoundCode, "Post not found")
		}
		errMsg := "Failed to find post"
		ctx := domaincontext.WithValue(ctx, "error", err.Error())
		slog.ErrorContext(ctx, errMsg)
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, errMsg)
	}

	voPostID, err := valueobject.ParsePostID(dbPost.ID)
	if err != nil {
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Invalid post ID")
	}

	voUserID, err := valueobject.ParseUserID(dbPost.UserID)
	if err != nil {
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Invalid user ID")
	}

	voTitle, err := valueobject.NewPostTitle(dbPost.Title)
	if err != nil {
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Invalid post title")
	}

	voContent, err := valueobject.NewPostContent(dbPost.Content)
	if err != nil {
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Invalid post content")
	}

	voStatus, err := valueobject.NewPostStatus(dbPost.Status)
	if err != nil {
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Invalid post status")
	}

	firstPublishedAt := &dbPost.FirstPublishedAt.Time
	if !dbPost.FirstPublishedAt.Valid {
		firstPublishedAt = nil
	}

	contentUpdatedAt := &dbPost.ContentUpdatedAt.Time
	if !dbPost.ContentUpdatedAt.Valid {
		contentUpdatedAt = nil
	}

	var voTags []valueobject.Tag
	if dbPost.R != nil && dbPost.R.Tags != nil {
		voTags = make([]valueobject.Tag, 0, len(dbPost.R.Tags))
		for _, tag := range dbPost.R.Tags {
			voTags = append(voTags, valueobject.Tag(tag.Name))
		}
	}

	post := &entity.Post{
		ID:               voPostID,
		UserID:           voUserID,
		Title:            voTitle,
		Content:          voContent,
		Status:           voStatus,
		CreatedAt:        dbPost.CreatedAt,
		UpdatedAt:        dbPost.UpdatedAt,
		FirstPublishedAt: firstPublishedAt,
		ContentUpdatedAt: contentUpdatedAt,
		Tags:             voTags,
	}

	return post, nil
}
