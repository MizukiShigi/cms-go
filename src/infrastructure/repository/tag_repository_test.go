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

func TestTagRepository_FindByPostID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{db: db}
	ctx := context.Background()

	postID, _ := valueobject.NewPostID("123e4567-e89b-12d3-a456-426614174000")
	tag1ID := "456e7890-e89b-12d3-a456-426614174000"
	tag1Name := "go"
	tag2ID := "789e0123-e89b-12d3-a456-426614174000"
	tag2Name := "testing"

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(tag1ID, tag1Name, "2023-01-01 00:00:00", "2023-01-01 00:00:00").
		AddRow(tag2ID, tag2Name, "2023-01-01 00:00:00", "2023-01-01 00:00:00")

	mock.ExpectQuery("SELECT t.id, t.name, t.created_at, t.updated_at FROM tags t INNER JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = ?").
		WithArgs(postID.String()).
		WillReturnRows(rows)

	tags, err := repo.FindByPostID(ctx, postID)

	assert.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, tag1ID, tags[0].ID().String())
	assert.Equal(t, tag1Name, tags[0].Name().String())
	assert.Equal(t, tag2ID, tags[1].ID().String())
	assert.Equal(t, tag2Name, tags[1].Name().String())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_FindByPostID_NoTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{db: db}
	ctx := context.Background()

	postID, _ := valueobject.NewPostID("123e4567-e89b-12d3-a456-426614174000")

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})

	mock.ExpectQuery("SELECT t.id, t.name, t.created_at, t.updated_at FROM tags t INNER JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = ?").
		WithArgs(postID.String()).
		WillReturnRows(rows)

	tags, err := repo.FindByPostID(ctx, postID)

	assert.NoError(t, err)
	assert.Len(t, tags, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_FindByPostID_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{db: db}
	ctx := context.Background()

	postID, _ := valueobject.NewPostID("123e4567-e89b-12d3-a456-426614174000")

	mock.ExpectQuery("SELECT t.id, t.name, t.created_at, t.updated_at FROM tags t INNER JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = ?").
		WithArgs(postID.String()).
		WillReturnError(sql.ErrConnDone)

	tags, err := repo.FindByPostID(ctx, postID)

	assert.Error(t, err)
	assert.Nil(t, tags)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Failed to find tags by post ID")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_FindByPostID_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{db: db}
	ctx := context.Background()

	postID, _ := valueobject.NewPostID("123e4567-e89b-12d3-a456-426614174000")

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow("invalid-uuid", "go", "2023-01-01 00:00:00", "2023-01-01 00:00:00")

	mock.ExpectQuery("SELECT t.id, t.name, t.created_at, t.updated_at FROM tags t INNER JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = ?").
		WithArgs(postID.String()).
		WillReturnRows(rows)

	tags, err := repo.FindByPostID(ctx, postID)

	assert.Error(t, err)
	assert.Nil(t, tags)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Failed to scan tag")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_FindOrCreateByName_ExistingTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{db: db}
	ctx := context.Background()

	tagName, _ := valueobject.NewTagName("go")
	inputTag, _ := entity.NewTagWithName(tagName)

	existingTagID := "456e7890-e89b-12d3-a456-426614174000"

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(existingTagID, tagName.String(), "2023-01-01 00:00:00", "2023-01-01 00:00:00")

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM tags WHERE name = ?").
		WithArgs(tagName.String()).
		WillReturnRows(rows)

	tag, err := repo.FindOrCreateByName(ctx, inputTag)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, existingTagID, tag.ID().String())
	assert.Equal(t, tagName.String(), tag.Name().String())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_FindOrCreateByName_CreateNewTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{db: db}
	ctx := context.Background()

	tagName, _ := valueobject.NewTagName("newtag")
	inputTag, _ := entity.NewTagWithName(tagName)

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM tags WHERE name = ?").
		WithArgs(tagName.String()).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO tags").
		WithArgs(
			sqlmock.AnyArg(),
			tagName.String(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tag, err := repo.FindOrCreateByName(ctx, inputTag)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, tagName.String(), tag.Name().String())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_FindOrCreateByName_FindDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{db: db}
	ctx := context.Background()

	tagName, _ := valueobject.NewTagName("go")
	inputTag, _ := entity.NewTagWithName(tagName)

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM tags WHERE name = ?").
		WithArgs(tagName.String()).
		WillReturnError(sql.ErrConnDone)

	tag, err := repo.FindOrCreateByName(ctx, inputTag)

	assert.Error(t, err)
	assert.Nil(t, tag)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Failed to find tag by name")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_FindOrCreateByName_CreateDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{db: db}
	ctx := context.Background()

	tagName, _ := valueobject.NewTagName("newtag")
	inputTag, _ := entity.NewTagWithName(tagName)

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM tags WHERE name = ?").
		WithArgs(tagName.String()).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO tags").
		WithArgs(sqlmock.AnyArg(), tagName.String(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	tag, err := repo.FindOrCreateByName(ctx, inputTag)

	assert.Error(t, err)
	assert.Nil(t, tag)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Failed to create tag")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_FindOrCreateByName_FindScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{db: db}
	ctx := context.Background()

	tagName, _ := valueobject.NewTagName("go")
	inputTag, _ := entity.NewTagWithName(tagName)

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow("invalid-uuid", tagName.String(), "2023-01-01 00:00:00", "2023-01-01 00:00:00")

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM tags WHERE name = ?").
		WithArgs(tagName.String()).
		WillReturnRows(rows)

	tag, err := repo.FindOrCreateByName(ctx, inputTag)

	assert.Error(t, err)
	assert.Nil(t, tag)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Failed to scan tag")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_FindOrCreateByName_EmptyTagName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{db: db}
	ctx := context.Background()

	tagName, _ := valueobject.NewTagName("")
	inputTag, _ := entity.NewTagWithName(tagName)

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM tags WHERE name = ?").
		WithArgs("").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO tags").
		WithArgs(sqlmock.AnyArg(), "", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tag, err := repo.FindOrCreateByName(ctx, inputTag)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "", tag.Name().String())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_FindOrCreateByName_DuplicateInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{db: db}
	ctx := context.Background()

	tagName, _ := valueobject.NewTagName("duplicate")
	inputTag, _ := entity.NewTagWithName(tagName)

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM tags WHERE name = ?").
		WithArgs(tagName.String()).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO tags").
		WithArgs(sqlmock.AnyArg(), tagName.String(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 0))

	tag, err := repo.FindOrCreateByName(ctx, inputTag)

	assert.Error(t, err)
	assert.Nil(t, tag)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Failed to create tag")
	assert.NoError(t, mock.ExpectationsWereMet())
}