package usecase

import (
	"context"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	repositoryMock "github.com/MizukiShigi/cms-go/mocks/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreatePostUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionManager := repositoryMock.NewMockTransactionManager(ctrl)
	mockPostRepo := repositoryMock.NewMockPostRepository(ctrl)
	mockTagRepo := repositoryMock.NewMockTagRepository(ctrl)

	t.Run("タグありの投稿作成が成功する", func(t *testing.T) {
		usecase := NewCreatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)

		title, _ := valueobject.NewPostTitle("テスト投稿")
		content, _ := valueobject.NewPostContent("テスト内容")
		userID := valueobject.NewUserID()
		status := valueobject.StatusPublished
		tagName1, _ := valueobject.NewTagName("タグ1")
		tagName2, _ := valueobject.NewTagName("タグ2")

		input := &CreatePostInput{
			Title:   title,
			Content: content,
			Tags:    []valueobject.TagName{tagName1, tagName2},
			UserID:  userID,
			Status:  status,
		}

		tag1 := entity.NewTagWithName(tagName1)
		tag2 := entity.NewTagWithName(tagName2)

		// トランザクション内の処理をモック
		mockTransactionManager.EXPECT().Transaction(context.Background(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				// 投稿作成
				mockPostRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

				// タグ作成/取得
				mockTagRepo.EXPECT().FindOrCreateByName(ctx, gomock.Any()).Return(tag1, nil)
				mockTagRepo.EXPECT().FindOrCreateByName(ctx, gomock.Any()).Return(tag2, nil)

				// タグ設定
				mockPostRepo.EXPECT().SetTags(ctx, gomock.Any(), gomock.Any()).Return(nil)

				return fn(ctx)
			})

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, title, output.Title)
		assert.Equal(t, content, output.Content)
		assert.Equal(t, userID, output.UserID)
		assert.Len(t, output.Tags, 2)
	})

	t.Run("タグなしの投稿作成が成功する", func(t *testing.T) {
		usecase := NewCreatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)

		title, _ := valueobject.NewPostTitle("テスト投稿")
		content, _ := valueobject.NewPostContent("テスト内容")
		userID := valueobject.NewUserID()
		status := valueobject.StatusPublished

		input := &CreatePostInput{
			Title:   title,
			Content: content,
			Tags:    []valueobject.TagName{},
			UserID:  userID,
			Status:  status,
		}

		// トランザクション内の処理をモック
		mockTransactionManager.EXPECT().Transaction(context.Background(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				// 投稿作成
				mockPostRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

				// タグなしなのでSetTagsのみ
				mockPostRepo.EXPECT().SetTags(ctx, gomock.Any(), gomock.Any()).Return(nil)

				return fn(ctx)
			})

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, title, output.Title)
		assert.Equal(t, content, output.Content)
		assert.Equal(t, userID, output.UserID)
		assert.Len(t, output.Tags, 0)
	})

	t.Run("投稿作成に失敗する", func(t *testing.T) {
		usecase := NewCreatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)

		title, _ := valueobject.NewPostTitle("テスト投稿")
		content, _ := valueobject.NewPostContent("テスト内容")
		userID := valueobject.NewUserID()
		status := valueobject.StatusPublished

		input := &CreatePostInput{
			Title:   title,
			Content: content,
			Tags:    []valueobject.TagName{},
			UserID:  userID,
			Status:  status,
		}

		// トランザクション内で投稿作成が失敗
		mockTransactionManager.EXPECT().Transaction(context.Background(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockPostRepo.EXPECT().Create(ctx, gomock.Any()).
					Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Database error"))

				return fn(ctx)
			})

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)

		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	})

	t.Run("タグ作成に失敗する", func(t *testing.T) {
		usecase := NewCreatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)

		title, _ := valueobject.NewPostTitle("テスト投稿")
		content, _ := valueobject.NewPostContent("テスト内容")
		userID := valueobject.NewUserID()
		status := valueobject.StatusPublished
		tagName, _ := valueobject.NewTagName("タグ1")

		input := &CreatePostInput{
			Title:   title,
			Content: content,
			Tags:    []valueobject.TagName{tagName},
			UserID:  userID,
			Status:  status,
		}

		// トランザクション内でタグ作成が失敗
		mockTransactionManager.EXPECT().Transaction(context.Background(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockPostRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mockTagRepo.EXPECT().FindOrCreateByName(ctx, gomock.Any()).
					Return(nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Tag creation failed"))

				return fn(ctx)
			})

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)

		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	})

	t.Run("タグ設定に失敗する", func(t *testing.T) {
		usecase := NewCreatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)

		title, _ := valueobject.NewPostTitle("テスト投稿")
		content, _ := valueobject.NewPostContent("テスト内容")
		userID := valueobject.NewUserID()
		status := valueobject.StatusPublished
		tagName, _ := valueobject.NewTagName("タグ1")

		input := &CreatePostInput{
			Title:   title,
			Content: content,
			Tags:    []valueobject.TagName{tagName},
			UserID:  userID,
			Status:  status,
		}

		tag := entity.NewTagWithName(tagName)

		// トランザクション内でタグ設定が失敗
		mockTransactionManager.EXPECT().Transaction(context.Background(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockPostRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mockTagRepo.EXPECT().FindOrCreateByName(ctx, gomock.Any()).Return(tag, nil)
				mockPostRepo.EXPECT().SetTags(ctx, gomock.Any(), gomock.Any()).
					Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Set tags failed"))

				return fn(ctx)
			})

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)

		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
	})

	t.Run("トランザクション自体が失敗する", func(t *testing.T) {
		usecase := NewCreatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)

		title, _ := valueobject.NewPostTitle("テスト投稿")
		content, _ := valueobject.NewPostContent("テスト内容")
		userID := valueobject.NewUserID()
		status := valueobject.StatusPublished

		input := &CreatePostInput{
			Title:   title,
			Content: content,
			Tags:    []valueobject.TagName{},
			UserID:  userID,
			Status:  status,
		}

		// トランザクション自体が失敗
		mockTransactionManager.EXPECT().Transaction(context.Background(), gomock.Any()).
			Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Transaction failed"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
