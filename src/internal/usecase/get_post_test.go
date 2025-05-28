package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	repositoryMock "github.com/MizukiShigi/cms-go/mocks/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetPostUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := repositoryMock.NewMockPostRepository(ctrl)

	t.Run("公開済み投稿の取得が成功する", func(t *testing.T) {
		usecase := NewGetPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()
		title, _ := valueobject.NewPostTitle("テスト投稿")
		content, _ := valueobject.NewPostContent("テスト内容")
		userID, _ := valueobject.NewUserID()
		status := valueobject.Published
		
		post, _ := entity.NewPost(title, content, userID, status)
		post.ID = postID
		
		now := time.Now()
		post.FirstPublishedAt = &now
		post.ContentUpdatedAt = &now

		input := &GetPostInput{
			ID: postID,
		}

		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(post, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, postID, output.ID)
		assert.Equal(t, title, output.Title)
		assert.Equal(t, content, output.Content)
		assert.Equal(t, status, output.Status)
		assert.Equal(t, post.FirstPublishedAt, output.FirstPublishedAt)
		assert.Equal(t, post.ContentUpdatedAt, output.ContentUpdatedAt)
	})

	t.Run("下書き投稿の取得が成功する", func(t *testing.T) {
		usecase := NewGetPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()
		title, _ := valueobject.NewPostTitle("下書き投稿")
		content, _ := valueobject.NewPostContent("下書き内容")
		userID, _ := valueobject.NewUserID()
		status := valueobject.Draft
		
		post, _ := entity.NewPost(title, content, userID, status)
		post.ID = postID

		input := &GetPostInput{
			ID: postID,
		}

		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(post, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, postID, output.ID)
		assert.Equal(t, title, output.Title)
		assert.Equal(t, content, output.Content)
		assert.Equal(t, status, output.Status)
		assert.Nil(t, output.FirstPublishedAt)
		assert.Nil(t, output.ContentUpdatedAt)
	})

	t.Run("投稿が存在しない場合にエラーが発生する", func(t *testing.T) {
		usecase := NewGetPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()

		input := &GetPostInput{
			ID: postID,
		}

		mockPostRepo.EXPECT().Get(context.Background(), postID).
			Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "Post not found"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
		
		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.NotFoundCode, myErr.Code)
	})

	t.Run("リポジトリエラーでエラーが発生する", func(t *testing.T) {
		usecase := NewGetPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()

		input := &GetPostInput{
			ID: postID,
		}

		mockPostRepo.EXPECT().Get(context.Background(), postID).
			Return(nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Database error"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
		
		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	})
}