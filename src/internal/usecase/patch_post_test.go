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

func TestPatchPostUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := repositoryMock.NewMockPostRepository(ctrl)

	t.Run("タイトルのみの更新が成功する", func(t *testing.T) {
		usecase := NewPatchPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()
		oldTitle, _ := valueobject.NewPostTitle("旧タイトル")
		newTitle, _ := valueobject.NewPostTitle("新タイトル")
		content, _ := valueobject.NewPostContent("内容")
		userID, _ := valueobject.NewUserID()
		status := valueobject.Published

		// 既存の投稿
		oldPost, _ := entity.NewPost(oldTitle, content, userID, status)
		oldPost.ID = postID

		// 更新後の投稿
		updatedPost, _ := entity.NewPost(newTitle, content, userID, status)
		updatedPost.ID = postID

		input := &PatchPostInput{
			ID:    postID,
			Title: &newTitle,
		}

		// 既存投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(oldPost, nil)
		
		// 投稿更新
		mockPostRepo.EXPECT().Update(context.Background(), gomock.Any()).Return(nil)
		
		// 更新後の投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(updatedPost, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, newTitle, output.Title)
		assert.Equal(t, content, output.Content)
		assert.Equal(t, status, output.Status)
	})

	t.Run("内容のみの更新が成功する", func(t *testing.T) {
		usecase := NewPatchPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()
		title, _ := valueobject.NewPostTitle("タイトル")
		oldContent, _ := valueobject.NewPostContent("旧内容")
		newContent, _ := valueobject.NewPostContent("新内容")
		userID, _ := valueobject.NewUserID()
		status := valueobject.Published

		// 既存の投稿
		oldPost, _ := entity.NewPost(title, oldContent, userID, status)
		oldPost.ID = postID

		// 更新後の投稿
		updatedPost, _ := entity.NewPost(title, newContent, userID, status)
		updatedPost.ID = postID

		input := &PatchPostInput{
			ID:      postID,
			Content: &newContent,
		}

		// 既存投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(oldPost, nil)
		
		// 投稿更新
		mockPostRepo.EXPECT().Update(context.Background(), gomock.Any()).Return(nil)
		
		// 更新後の投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(updatedPost, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, title, output.Title)
		assert.Equal(t, newContent, output.Content)
		assert.Equal(t, status, output.Status)
	})

	t.Run("ステータスのみの更新が成功する", func(t *testing.T) {
		usecase := NewPatchPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()
		title, _ := valueobject.NewPostTitle("タイトル")
		content, _ := valueobject.NewPostContent("内容")
		userID, _ := valueobject.NewUserID()
		oldStatus := valueobject.Draft
		newStatus := valueobject.Published

		// 既存の投稿（下書き）
		oldPost, _ := entity.NewPost(title, content, userID, oldStatus)
		oldPost.ID = postID

		// 更新後の投稿（公開済み）
		updatedPost, _ := entity.NewPost(title, content, userID, newStatus)
		updatedPost.ID = postID
		now := time.Now()
		updatedPost.FirstPublishedAt = &now

		input := &PatchPostInput{
			ID:     postID,
			Status: &newStatus,
		}

		// 既存投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(oldPost, nil)
		
		// 投稿更新
		mockPostRepo.EXPECT().Update(context.Background(), gomock.Any()).Return(nil)
		
		// 更新後の投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(updatedPost, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, title, output.Title)
		assert.Equal(t, content, output.Content)
		assert.Equal(t, newStatus, output.Status)
		assert.NotNil(t, output.FirstPublishedAt)
	})

	t.Run("タグのみの更新が成功する", func(t *testing.T) {
		usecase := NewPatchPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()
		title, _ := valueobject.NewPostTitle("タイトル")
		content, _ := valueobject.NewPostContent("内容")
		userID, _ := valueobject.NewUserID()
		status := valueobject.Published
		tag1, _ := valueobject.NewTagName("タグ1")
		tag2, _ := valueobject.NewTagName("タグ2")

		// 既存の投稿
		oldPost, _ := entity.NewPost(title, content, userID, status)
		oldPost.ID = postID

		// 更新後の投稿
		updatedPost, _ := entity.NewPost(title, content, userID, status)
		updatedPost.ID = postID
		updatedPost.Tags = []valueobject.TagName{tag1, tag2}

		input := &PatchPostInput{
			ID:   postID,
			Tags: []valueobject.TagName{tag1, tag2},
		}

		// 既存投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(oldPost, nil)
		
		// 投稿更新
		mockPostRepo.EXPECT().Update(context.Background(), gomock.Any()).Return(nil)
		
		// 更新後の投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(updatedPost, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, title, output.Title)
		assert.Equal(t, content, output.Content)
		assert.Equal(t, status, output.Status)
		assert.Len(t, output.Tags, 2)
		assert.Contains(t, output.Tags, tag1)
		assert.Contains(t, output.Tags, tag2)
	})

	t.Run("複数フィールドの同時更新が成功する", func(t *testing.T) {
		usecase := NewPatchPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()
		oldTitle, _ := valueobject.NewPostTitle("旧タイトル")
		newTitle, _ := valueobject.NewPostTitle("新タイトル")
		oldContent, _ := valueobject.NewPostContent("旧内容")
		newContent, _ := valueobject.NewPostContent("新内容")
		userID, _ := valueobject.NewUserID()
		oldStatus := valueobject.Draft
		newStatus := valueobject.Published
		tag, _ := valueobject.NewTagName("タグ1")

		// 既存の投稿
		oldPost, _ := entity.NewPost(oldTitle, oldContent, userID, oldStatus)
		oldPost.ID = postID

		// 更新後の投稿
		updatedPost, _ := entity.NewPost(newTitle, newContent, userID, newStatus)
		updatedPost.ID = postID
		updatedPost.Tags = []valueobject.TagName{tag}

		input := &PatchPostInput{
			ID:      postID,
			Title:   &newTitle,
			Content: &newContent,
			Status:  &newStatus,
			Tags:    []valueobject.TagName{tag},
		}

		// 既存投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(oldPost, nil)
		
		// 投稿更新
		mockPostRepo.EXPECT().Update(context.Background(), gomock.Any()).Return(nil)
		
		// 更新後の投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(updatedPost, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, newTitle, output.Title)
		assert.Equal(t, newContent, output.Content)
		assert.Equal(t, newStatus, output.Status)
		assert.Len(t, output.Tags, 1)
		assert.Contains(t, output.Tags, tag)
	})

	t.Run("投稿が存在しない場合にエラーが発生する", func(t *testing.T) {
		usecase := NewPatchPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()
		title, _ := valueobject.NewPostTitle("タイトル")

		input := &PatchPostInput{
			ID:    postID,
			Title: &title,
		}

		// 投稿が見つからない
		mockPostRepo.EXPECT().Get(context.Background(), postID).
			Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "Post not found"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
		
		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.NotFoundCode, myErr.Code)
	})

	t.Run("不正なステータス遷移でエラーが発生する", func(t *testing.T) {
		usecase := NewPatchPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()
		title, _ := valueobject.NewPostTitle("タイトル")
		content, _ := valueobject.NewPostContent("内容")
		userID, _ := valueobject.NewUserID()
		currentStatus := valueobject.Published
		invalidStatus := valueobject.Draft // 公開済みから下書きには戻せない

		// 既存の投稿（公開済み）
		post, _ := entity.NewPost(title, content, userID, currentStatus)
		post.ID = postID

		input := &PatchPostInput{
			ID:     postID,
			Status: &invalidStatus,
		}

		// 既存投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(post, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("投稿更新に失敗する", func(t *testing.T) {
		usecase := NewPatchPostUsecase(mockPostRepo)
		
		postID, _ := valueobject.NewPostID()
		title, _ := valueobject.NewPostTitle("タイトル")
		content, _ := valueobject.NewPostContent("内容")
		userID, _ := valueobject.NewUserID()
		status := valueobject.Published

		post, _ := entity.NewPost(title, content, userID, status)
		post.ID = postID

		newTitle, _ := valueobject.NewPostTitle("新タイトル")

		input := &PatchPostInput{
			ID:    postID,
			Title: &newTitle,
		}

		// 既存投稿取得
		mockPostRepo.EXPECT().Get(context.Background(), postID).Return(post, nil)
		
		// 投稿更新が失敗
		mockPostRepo.EXPECT().Update(context.Background(), gomock.Any()).
			Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Update failed"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}