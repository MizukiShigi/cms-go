package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestPostRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Post")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")
	post, _ := entity.NewPost(userID, title, content, status)

	mock.ExpectExec("INSERT INTO posts").
		WithArgs(
			sqlmock.AnyArg(),
			post.UserID().String(),
			post.Title().String(),
			post.Content().String(),
			post.Status().String(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(ctx, post)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Create_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Post")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")
	post, _ := entity.NewPost(userID, title, content, status)

	mock.ExpectExec("INSERT INTO posts").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(ctx, post)

	assert.Error(t, err)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Failed to create post")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Get_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	postID, _ := valueobject.NewPostID("123e4567-e89b-12d3-a456-426614174000")
	userIDStr := "456e7890-e89b-12d3-a456-426614174000"
	titleStr := "Test Post"
	contentStr := "Test content"
	statusStr := "draft"

	rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "status", "created_at", "updated_at", "first_published_at", "content_updated_at"}).
		AddRow(postID.String(), userIDStr, titleStr, contentStr, statusStr, "2023-01-01 00:00:00", "2023-01-01 00:00:00", nil, nil)

	mock.ExpectQuery("SELECT (.+) FROM posts WHERE id = ?").
		WithArgs(postID.String()).
		WillReturnRows(rows)

	post, err := repo.Get(ctx, postID)

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, postID, post.ID())
	assert.Equal(t, userIDStr, post.UserID().String())
	assert.Equal(t, titleStr, post.Title().String())
	assert.Equal(t, contentStr, post.Content().String())
	assert.Equal(t, statusStr, post.Status().String())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	postID, _ := valueobject.NewPostID("123e4567-e89b-12d3-a456-426614174000")

	mock.ExpectQuery("SELECT (.+) FROM posts WHERE id = ?").
		WithArgs(postID.String()).
		WillReturnError(sql.ErrNoRows)

	post, err := repo.Get(ctx, postID)

	assert.Error(t, err)
	assert.Nil(t, post)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.NotFoundCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Post not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Get_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	postID, _ := valueobject.NewPostID("123e4567-e89b-12d3-a456-426614174000")

	mock.ExpectQuery("SELECT (.+) FROM posts WHERE id = ?").
		WithArgs(postID.String()).
		WillReturnError(sql.ErrConnDone)

	post, err := repo.Get(ctx, postID)

	assert.Error(t, err)
	assert.Nil(t, post)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Failed to get post")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Updated Post")
	content, _ := valueobject.NewPostContent("Updated content")
	status, _ := valueobject.NewPostStatus("published")
	post, _ := entity.NewPost(userID, title, content, status)

	mock.ExpectExec("UPDATE posts SET").
		WithArgs(
			post.Title().String(),
			post.Content().String(),
			post.Status().String(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			post.ID().String(),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(ctx, post)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Update_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Updated Post")
	content, _ := valueobject.NewPostContent("Updated content")
	status, _ := valueobject.NewPostStatus("published")
	post, _ := entity.NewPost(userID, title, content, status)

	mock.ExpectExec("UPDATE posts SET").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Update(ctx, post)

	assert.Error(t, err)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.NotFoundCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Post not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Update_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Updated Post")
	content, _ := valueobject.NewPostContent("Updated content")
	status, _ := valueobject.NewPostStatus("published")
	post, _ := entity.NewPost(userID, title, content, status)

	mock.ExpectExec("UPDATE posts SET").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	err = repo.Update(ctx, post)

	assert.Error(t, err)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Failed to update post")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_SetTags_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Post")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")
	post, _ := entity.NewPost(userID, title, content, status)

	tagID1, _ := valueobject.NewTagID("tag1-id")
	tagName1, _ := valueobject.NewTagName("go")
	tag1, _ := entity.NewTag(tagID1, tagName1)

	tagID2, _ := valueobject.NewTagID("tag2-id")
	tagName2, _ := valueobject.NewTagName("testing")
	tag2, _ := entity.NewTag(tagID2, tagName2)

	tags := []*entity.Tag{tag1, tag2}

	mock.ExpectExec("DELETE FROM post_tags WHERE post_id = ?").
		WithArgs(post.ID().String()).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("INSERT INTO post_tags").
		WithArgs(post.ID().String(), tag1.ID().String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO post_tags").
		WithArgs(post.ID().String(), tag2.ID().String()).
		WillReturnResult(sqlmock.NewResult(2, 1))

	err = repo.SetTags(ctx, post, tags)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_SetTags_EmptyTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Post")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")
	post, _ := entity.NewPost(userID, title, content, status)

	tags := []*entity.Tag{}

	mock.ExpectExec("DELETE FROM post_tags WHERE post_id = ?").
		WithArgs(post.ID().String()).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.SetTags(ctx, post, tags)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_SetTags_DeleteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Post")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")
	post, _ := entity.NewPost(userID, title, content, status)

	tags := []*entity.Tag{}

	mock.ExpectExec("DELETE FROM post_tags WHERE post_id = ?").
		WithArgs(post.ID().String()).
		WillReturnError(sql.ErrConnDone)

	err = repo.SetTags(ctx, post, tags)

	assert.Error(t, err)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Failed to delete post tags")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_SetTags_InsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostRepository{db: db}
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Post")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")
	post, _ := entity.NewPost(userID, title, content, status)

	tagID, _ := valueobject.NewTagID("tag-id")
	tagName, _ := valueobject.NewTagName("go")
	tag, _ := entity.NewTag(tagID, tagName)

	tags := []*entity.Tag{tag}

	mock.ExpectExec("DELETE FROM post_tags WHERE post_id = ?").
		WithArgs(post.ID().String()).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("INSERT INTO post_tags").
		WithArgs(post.ID().String(), tag.ID().String()).
		WillReturnError(sql.ErrConnDone)

	err = repo.SetTags(ctx, post, tags)

	assert.Error(t, err)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Failed to insert post tag")
	assert.NoError(t, mock.ExpectationsWereMet())
}