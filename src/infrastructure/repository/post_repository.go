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
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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

	if err := dbPost.Insert(ctx, GetExecDB(ctx, r.db), boil.Infer()); err != nil {
		return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to create post")
	}

	return nil
}

func (r *PostRepository) Get(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	dbPost, err := models.Posts(qm.Where("id = ?", id.String()), qm.Load("Tags")).One(ctx, r.db)
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

	var voTags []valueobject.TagName
	if dbPost.R != nil && dbPost.R.Tags != nil {
		voTags = make([]valueobject.TagName, 0, len(dbPost.R.Tags))
		for _, tag := range dbPost.R.Tags {
			voTags = append(voTags, valueobject.TagName(tag.Name))
		}
	}

	post := entity.ParsePost(
		voPostID,
		voTitle,
		voContent,
		voUserID,
		voStatus,
		dbPost.CreatedAt,
		dbPost.UpdatedAt,
		firstPublishedAt,
		contentUpdatedAt,
		voTags,
	)

	return post, nil
}

func (r *PostRepository) Update(ctx context.Context, post *entity.Post) error {
	dbPost := &models.Post{
		ID:      post.ID.String(),
		Title:   post.Title.String(),
		Content: post.Content.String(),
		UserID:  post.UserID.String(),
		Status:  post.Status.String(),
		ContentUpdatedAt: ToNullable(
			post.ContentUpdatedAt,
			func(t time.Time) bool { return t.IsZero() },
			null.TimeFrom,
		),
		FirstPublishedAt: ToNullable(
			post.FirstPublishedAt,
			func(t time.Time) bool { return t.IsZero() },
			null.TimeFrom,
		),
	}

	if _, err := dbPost.Update(ctx, GetExecDB(ctx, r.db), boil.Infer()); err != nil {
		slog.ErrorContext(ctx, "Failed to update post", "error", err)
		return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to update post")
	}

	return nil
}

func (r *PostRepository) SetTags(ctx context.Context, post *entity.Post, tags []*entity.Tag) error {
	if tags == nil {
		slog.InfoContext(ctx, "No tags to set")
		return nil
	}

	dbPost := &models.Post{
		ID: post.ID.String(),
	}

	dbTags := make([]*models.Tag, 0, len(tags))
	for _, tag := range tags {
		dbTags = append(dbTags, &models.Tag{
			ID: tag.ID.String(),
		})
	}

	if err := dbPost.SetTags(ctx, GetExecDB(ctx, r.db), false, dbTags...); err != nil {
		return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to set tags")
	}

	return nil
}
