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

func TestUpdatePostUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionManager := repositoryMock.NewMockTransactionManager(ctrl)
	mockPostRepo := repositoryMock.NewMockPostRepository(ctrl)
	mockTagRepo := repositoryMock.NewMockTagRepository(ctrl)

	t.Run("全項目の投稿更新が成功する", func(t *testing.T) {
		usecase := NewUpdatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)

		postID := valueobject.NewPostID()
		oldTitle, _ := valueobject.NewPostTitle("旧タイトル")
		oldContent, _ := valueobject.NewPostContent("旧内容")
		newTitle, _ := valueobject.NewPostTitle("新タイトル")
		newContent, _ := valueobject.NewPostContent("新内容")
		userID := valueobject.NewUserID()
		status := valueobject.StatusPublished
		tagName, _ := valueobject.NewTagName("タグ1")

		// 既存の投稿
		oldPost, _ := entity.NewPost(oldTitle, oldContent, userID, status)
		oldPost.ID = postID
		oldTime := time.Now().Add(-time.Hour)
		oldPost.FirstPublishedAt = &oldTime

		// 更新後の投稿
		updatedPost, _ := entity.NewPost(newTitle, newContent, userID, status)
		updatedPost.ID = postID
		updatedPost.FirstPublishedAt = &oldTime
		newTime := time.Now()
		updatedPost.ContentUpdatedAt = &newTime

		input := &UpdatePostInput{
			ID:      postID,
			Title:   newTitle,
			Content: newContent,
			Tags:    []valueobject.TagName{tagName},
			Status:  status,
		}

		tag := entity.NewTagWithName(tagName)

		// 既存投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(oldPost, nil)

		// トランザクション内の処理をモック
		mockTransactionManager.EXPECT().Transaction(context.Background(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				// 投稿更新
				mockPostRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

				// タグ作成/取得
				mockTagRepo.EXPECT().FindOrCreateByName(ctx, gomock.Any()).Return(tag, nil)

				// タグ設定
				mockPostRepo.EXPECT().SetTags(ctx, gomock.Any(), gomock.Any()).Return(nil)

				return fn(ctx)
			})

		// 更新後の投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(updatedPost, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, postID, output.ID)
		assert.Equal(t, newTitle, output.Title)
		assert.Equal(t, newContent, output.Content)
		assert.Equal(t, status, output.Status)
		assert.Equal(t, &oldTime, output.FirstPublishedAt)
		assert.Equal(t, &newTime, output.ContentUpdatedAt)
	})

	t.Run("タグなしの投稿更新が成功する", func(t *testing.T) {
		usecase := NewUpdatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)

		postID := valueobject.NewPostID()
		oldTitle, _ := valueobject.NewPostTitle("旧タイトル")
		oldContent, _ := valueobject.NewPostContent("旧内容")
		newTitle, _ := valueobject.NewPostTitle("新タイトル")
		newContent, _ := valueobject.NewPostContent("新内容")
		userID := valueobject.NewUserID()
		status := valueobject.StatusPublished

		// 既存の投稿
		oldPost, _ := entity.NewPost(oldTitle, oldContent, userID, status)
		oldPost.ID = postID

		// 更新後の投稿
		updatedPost, _ := entity.NewPost(newTitle, newContent, userID, status)
		updatedPost.ID = postID

		input := &UpdatePostInput{
			ID:      postID,
			Title:   newTitle,
			Content: newContent,
			Tags:    []valueobject.TagName{},
			Status:  status,
		}

		// 既存投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(oldPost, nil)

		// トランザクション内の処理をモック
		mockTransactionManager.EXPECT().Transaction(context.Background(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				// 投稿更新
				mockPostRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

				// タグなしなのでSetTagsのみ
				mockPostRepo.EXPECT().SetTags(ctx, gomock.Any(), gomock.Any()).Return(nil)

				return fn(ctx)
			})

		// 更新後の投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(updatedPost, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, postID, output.ID)
		assert.Equal(t, newTitle, output.Title)
		assert.Equal(t, newContent, output.Content)
	})

	t.Run("投稿が存在しない場合にエラーが発生する", func(t *testing.T) {
		usecase := NewUpdatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)

		postID := valueobject.NewPostID()
		title, _ := valueobject.NewPostTitle("タイトル")
		content, _ := valueobject.NewPostContent("内容")

		input := &UpdatePostInput{
			ID:      postID,
			Title:   title,
			Content: content,
			Tags:    []valueobject.TagName{},
			Status:  valueobject.StatusPublished,
		}

		// 投稿が見つからない
		mockPostRepo.EXPECT().Get(context.Background(), postID).
			Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "Post not found"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("投稿更新に失敗する", func(t *testing.T) {
		usecase := NewUpdatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)

		postID := valueobject.NewPostID()
		title, _ := valueobject.NewPostTitle("タイトル")
		content, _ := valueobject.NewPostContent("内容")
		userID := valueobject.NewUserID()

		post, _ := entity.NewPost(title, content, userID, valueobject.StatusPublished)
		post.ID = postID

		input := &UpdatePostInput{
			ID:      postID,
			Title:   title,
			Content: content,
			Tags:    []valueobject.TagName{},
			Status:  valueobject.StatusPublished,
		}

		// 既存投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(post, nil)

		// トランザクション内で投稿更新が失敗
		mockTransactionManager.EXPECT().Transaction(context.Background(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockPostRepo.EXPECT().Update(ctx, gomock.Any()).
					Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Update failed"))

				return fn(ctx)
			})

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("タグ作成に失敗する", func(t *testing.T) {
		usecase := NewUpdatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)

		postID := valueobject.NewPostID()
		title, _ := valueobject.NewPostTitle("タイトル")
		content, _ := valueobject.NewPostContent("内容")
		userID := valueobject.NewUserID()
		tagName, _ := valueobject.NewTagName("タグ1")

		post, _ := entity.NewPost(title, content, userID, valueobject.StatusPublished)
		post.ID = postID

		input := &UpdatePostInput{
			ID:      postID,
			Title:   title,
			Content: content,
			Tags:    []valueobject.TagName{tagName},
			Status:  valueobject.StatusPublished,
		}

		// 既存投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(post, nil)

		// トランザクション内でタグ作成が失敗
		mockTransactionManager.EXPECT().Transaction(context.Background(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockPostRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
				mockTagRepo.EXPECT().FindOrCreateByName(ctx, gomock.Any()).
					Return(nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Tag creation failed"))

				return fn(ctx)
			})

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
