package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/repository"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	repositoryMock "github.com/MizukiShigi/cms-go/mocks/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListPostsUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := repositoryMock.NewMockPostRepository(ctrl)

	t.Run("投稿一覧の取得が成功する", func(t *testing.T) {
		usecase := NewListPostsUsecase(mockPostRepo)

		userID := valueobject.NewUserID()
		
		// テスト用の投稿データを作成
		post1Title, _ := valueobject.NewPostTitle("投稿1")
		post1Content, _ := valueobject.NewPostContent("投稿1の内容")
		post1, _ := entity.NewPost(post1Title, post1Content, userID, valueobject.StatusPublished)
		
		post2Title, _ := valueobject.NewPostTitle("投稿2")
		post2Content, _ := valueobject.NewPostContent("投稿2の内容")
		post2, _ := entity.NewPost(post2Title, post2Content, userID, valueobject.StatusDraft)

		posts := []*entity.Post{post1, post2}
		totalCount := 2

		// モックの設定
		mockPostRepo.EXPECT().
			List(gomock.Any(), gomock.Any()).
			Return(posts, totalCount, nil)

		// リクエスト作成
		req := &ListPostsRequest{
			Limit:  "20",
			Offset: "0",
			Sort:   "created_at_desc",
		}

		// 実行
		result, err := usecase.Execute(context.Background(), req)

		// 検証
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Posts, 2)
		assert.Equal(t, 2, result.Meta.Total)
		assert.Equal(t, 20, result.Meta.Limit)
		assert.Equal(t, 0, result.Meta.Offset)
		assert.False(t, result.Meta.HasNext)

		// 投稿の詳細を検証
		assert.Equal(t, post1.ID.String(), result.Posts[0].ID)
		assert.Equal(t, "投稿1", result.Posts[0].Title)
		assert.Equal(t, "published", result.Posts[0].Status)
		
		assert.Equal(t, post2.ID.String(), result.Posts[1].ID)
		assert.Equal(t, "投稿2", result.Posts[1].Title)
		assert.Equal(t, "draft", result.Posts[1].Status)
	})

	t.Run("ページネーションが正しく動作する", func(t *testing.T) {
		usecase := NewListPostsUsecase(mockPostRepo)

		posts := []*entity.Post{}
		totalCount := 50

		// モックの設定（has_nextがtrueになることを確認）
		expectedOptions := &repository.ListPostsOptions{
			Limit:  10,
			Offset: 20,
			Status: nil,
			Sort:   "created_at_desc",
		}

		mockPostRepo.EXPECT().
			List(gomock.Any(), expectedOptions).
			Return(posts, totalCount, nil)

		// リクエスト作成
		req := &ListPostsRequest{
			Limit:  "10",
			Offset: "20",
		}

		// 実行
		result, err := usecase.Execute(context.Background(), req)

		// 検証
		assert.NoError(t, err)
		assert.Equal(t, 50, result.Meta.Total)
		assert.Equal(t, 10, result.Meta.Limit)
		assert.Equal(t, 20, result.Meta.Offset)
		assert.True(t, result.Meta.HasNext) // 20 + 10 < 50
	})

	t.Run("ステータスフィルタが正しく動作する", func(t *testing.T) {
		usecase := NewListPostsUsecase(mockPostRepo)

		posts := []*entity.Post{}
		totalCount := 10

		// モックの設定（publishedステータスでフィルタ）
		publishedStatus := valueobject.StatusPublished
		expectedOptions := &repository.ListPostsOptions{
			Limit:  20,
			Offset: 0,
			Status: &publishedStatus,
			Sort:   "created_at_desc",
		}

		mockPostRepo.EXPECT().
			List(gomock.Any(), expectedOptions).
			Return(posts, totalCount, nil)

		// リクエスト作成
		req := &ListPostsRequest{
			Status: "published",
		}

		// 実行
		result, err := usecase.Execute(context.Background(), req)

		// 検証
		assert.NoError(t, err)
		assert.Equal(t, 10, result.Meta.Total)
	})

	t.Run("ソート設定が正しく動作する", func(t *testing.T) {
		usecase := NewListPostsUsecase(mockPostRepo)

		posts := []*entity.Post{}
		totalCount := 5

		// モックの設定（updated_at_ascでソート）
		expectedOptions := &repository.ListPostsOptions{
			Limit:  20,
			Offset: 0,
			Status: nil,
			Sort:   "updated_at_asc",
		}

		mockPostRepo.EXPECT().
			List(gomock.Any(), expectedOptions).
			Return(posts, totalCount, nil)

		// リクエスト作成
		req := &ListPostsRequest{
			Sort:   "updated_at_asc",
		}

		// 実行
		result, err := usecase.Execute(context.Background(), req)

		// 検証
		assert.NoError(t, err)
		assert.Equal(t, 5, result.Meta.Total)
	})

	t.Run("無効なパラメータがデフォルト値で処理される", func(t *testing.T) {
		usecase := NewListPostsUsecase(mockPostRepo)

		posts := []*entity.Post{}
		totalCount := 0

		// モックの設定（デフォルト値での呼び出し）
		expectedOptions := &repository.ListPostsOptions{
			Limit:  20, // デフォルト
			Offset: 0,  // デフォルト
			Status: nil,
			Sort:   "created_at_desc", // デフォルト
		}

		mockPostRepo.EXPECT().
			List(gomock.Any(), expectedOptions).
			Return(posts, totalCount, nil)

		// リクエスト作成（無効な値を設定）
		req := &ListPostsRequest{
			Limit:  "invalid",
			Offset: "-10",
			Status: "invalid_status",
			Sort:   "invalid_sort",
		}

		// 実行
		result, err := usecase.Execute(context.Background(), req)

		// 検証
		assert.NoError(t, err)
		assert.Equal(t, 20, result.Meta.Limit)
		assert.Equal(t, 0, result.Meta.Offset)
	})

	t.Run("リポジトリエラーが適切に処理される", func(t *testing.T) {
		usecase := NewListPostsUsecase(mockPostRepo)

		expectedError := valueobject.NewMyError(valueobject.InternalServerErrorCode, "Database error")

		// モックの設定（エラーを返す）
		mockPostRepo.EXPECT().
			List(gomock.Any(), gomock.Any()).
			Return(nil, 0, expectedError)

		// リクエスト作成
		req := &ListPostsRequest{}

		// 実行
		result, err := usecase.Execute(context.Background(), req)

		// 検証
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestListPostsUsecase_convertToSummary(t *testing.T) {
	usecase := NewListPostsUsecase(nil)

	t.Run("投稿のサマリー変換が正しく動作する", func(t *testing.T) {
		// テストデータ作成
		postID := valueobject.NewPostID()
		userID := valueobject.NewUserID()
		title, _ := valueobject.NewPostTitle("テスト投稿")
		content, _ := valueobject.NewPostContent("テスト内容")
		status := valueobject.StatusPublished
		
		tagName1, _ := valueobject.NewTagName("技術")
		tagName2, _ := valueobject.NewTagName("Go")
		tags := []valueobject.TagName{tagName1, tagName2}

		now := time.Now()
		firstPublishedAt := &now
		contentUpdatedAt := &now

		post := entity.ParsePost(
			postID,
			title,
			content,
			userID,
			status,
			now,
			now,
			firstPublishedAt,
			contentUpdatedAt,
			tags,
		)

		// 変換実行
		summary := usecase.convertToSummary(post)

		// 検証
		assert.Equal(t, postID.String(), summary.ID)
		assert.Equal(t, "テスト投稿", summary.Title)
		assert.Equal(t, "published", summary.Status)
		assert.Len(t, summary.Tags, 2)
		assert.Contains(t, summary.Tags, "技術")
		assert.Contains(t, summary.Tags, "Go")
		assert.NotNil(t, summary.FirstPublishedAt) 
		assert.NotNil(t, summary.ContentUpdatedAt)
	})

	t.Run("null日付が正しく処理される", func(t *testing.T) {
		// テストデータ作成（null日付）
		postID := valueobject.NewPostID()
		userID := valueobject.NewUserID()
		title, _ := valueobject.NewPostTitle("下書き投稿")
		content, _ := valueobject.NewPostContent("下書き内容")
		status := valueobject.StatusDraft

		now := time.Now()
		post := entity.ParsePost(
			postID,
			title,
			content,
			userID,
			status,
			now,
			now,
			nil, // firstPublishedAt is null
			nil, // contentUpdatedAt is null
			nil,
		)

		// 変換実行
		summary := usecase.convertToSummary(post)

		// 検証
		assert.Equal(t, "draft", summary.Status)
		assert.Nil(t, summary.FirstPublishedAt)
		assert.Nil(t, summary.ContentUpdatedAt)
		assert.Empty(t, summary.Tags)
	})
}