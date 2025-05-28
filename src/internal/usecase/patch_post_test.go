package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	mock_repository "github.com/MizukiShigi/cms-go/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPatchPostUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock_repository.NewMockPostRepository(ctrl)

	tests := []struct {
		name      string
		input     *PatchPostInput
		setupMock func()
		wantErr   bool
		errCode   string
		checkResult func(t *testing.T, output *PatchPostOutput)
	}{
		{
			name: "正常系: タイトルのみ更新",
			input: &PatchPostInput{
				ID:      valueobject.NewPostID(),
				Title:   func() *valueobject.PostTitle { t := valueobject.PostTitle("新しいタイトル"); return &t }(),
				Content: nil,
				Status:  nil,
				Tags:    nil,
			},
			setupMock: func() {
				originalTitle, _ := valueobject.NewPostTitle("元のタイトル")
				originalContent, _ := valueobject.NewPostContent("元の内容")
				userID := valueobject.NewUserID()
				now := time.Now()
				
				originalPost := entity.ParsePost(
					valueobject.NewPostID(),
					originalTitle,
					originalContent,
					userID,
					valueobject.StatusDraft,
					now,
					now,
					nil,
					&now,
					[]valueobject.TagName{"既存タグ"},
				)

				updatedPost := entity.ParsePost(
					originalPost.ID,
					valueobject.PostTitle("新しいタイトル"),
					originalContent,
					userID,
					valueobject.StatusDraft,
					now,
					now.Add(time.Minute),
					nil,
					&now,
					[]valueobject.TagName{"既存タグ"},
				)

				// 投稿取得（更新前）
				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(originalPost, nil).Times(1)

				// 投稿更新
				mockPostRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)

				// 投稿取得（更新後）
				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(updatedPost, nil).Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *PatchPostOutput) {
				assert.Equal(t, "新しいタイトル", string(output.Title))
				assert.Equal(t, "元の内容", string(output.Content))
				assert.Equal(t, valueobject.StatusDraft, output.Status)
				assert.Len(t, output.Tags, 1)
				assert.Contains(t, output.Tags, valueobject.TagName("既存タグ"))
			},
		},
		{
			name: "正常系: 内容のみ更新",
			input: &PatchPostInput{
				ID:      valueobject.NewPostID(),
				Title:   nil,
				Content: func() *valueobject.PostContent { c := valueobject.PostContent("新しい内容"); return &c }(),
				Status:  nil,
				Tags:    nil,
			},
			setupMock: func() {
				originalTitle, _ := valueobject.NewPostTitle("元のタイトル")
				originalContent, _ := valueobject.NewPostContent("元の内容")
				userID := valueobject.NewUserID()
				now := time.Now()
				
				originalPost := entity.ParsePost(
					valueobject.NewPostID(),
					originalTitle,
					originalContent,
					userID,
					valueobject.StatusDraft,
					now,
					now,
					nil,
					&now,
					[]valueobject.TagName{},
				)

				updatedPost := entity.ParsePost(
					originalPost.ID,
					originalTitle,
					valueobject.PostContent("新しい内容"),
					userID,
					valueobject.StatusDraft,
					now,
					now.Add(time.Minute),
					nil,
					&now,
					[]valueobject.TagName{},
				)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(originalPost, nil).Times(1)

				mockPostRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(updatedPost, nil).Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *PatchPostOutput) {
				assert.Equal(t, "元のタイトル", string(output.Title))
				assert.Equal(t, "新しい内容", string(output.Content))
				assert.Equal(t, valueobject.StatusDraft, output.Status)
			},
		},
		{
			name: "正常系: ステータス更新（下書き→公開）",
			input: &PatchPostInput{
				ID:      valueobject.NewPostID(),
				Title:   nil,
				Content: nil,
				Status:  func() *valueobject.PostStatus { s := valueobject.StatusPublished; return &s }(),
				Tags:    nil,
			},
			setupMock: func() {
				originalTitle, _ := valueobject.NewPostTitle("元のタイトル")
				originalContent, _ := valueobject.NewPostContent("元の内容")
				userID := valueobject.NewUserID()
				now := time.Now()
				
				originalPost := entity.ParsePost(
					valueobject.NewPostID(),
					originalTitle,
					originalContent,
					userID,
					valueobject.StatusDraft,
					now,
					now,
					nil,
					&now,
					[]valueobject.TagName{},
				)

				updatedPost := entity.ParsePost(
					originalPost.ID,
					originalTitle,
					originalContent,
					userID,
					valueobject.StatusPublished,
					now,
					now.Add(time.Minute),
					&now, // 公開されたので初回公開日時が設定される
					&now,
					[]valueobject.TagName{},
				)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(originalPost, nil).Times(1)

				mockPostRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(updatedPost, nil).Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *PatchPostOutput) {
				assert.Equal(t, valueobject.StatusPublished, output.Status)
				assert.NotNil(t, output.FirstPublishedAt)
			},
		},
		{
			name: "正常系: タグ更新",
			input: &PatchPostInput{
				ID:      valueobject.NewPostID(),
				Title:   nil,
				Content: nil,
				Status:  nil,
				Tags:    []valueobject.TagName{"新タグ1", "新タグ2"},
			},
			setupMock: func() {
				originalTitle, _ := valueobject.NewPostTitle("元のタイトル")
				originalContent, _ := valueobject.NewPostContent("元の内容")
				userID := valueobject.NewUserID()
				now := time.Now()
				
				originalPost := entity.ParsePost(
					valueobject.NewPostID(),
					originalTitle,
					originalContent,
					userID,
					valueobject.StatusDraft,
					now,
					now,
					nil,
					&now,
					[]valueobject.TagName{"古いタグ"},
				)

				updatedPost := entity.ParsePost(
					originalPost.ID,
					originalTitle,
					originalContent,
					userID,
					valueobject.StatusDraft,
					now,
					now.Add(time.Minute),
					nil,
					&now,
					[]valueobject.TagName{"新タグ1", "新タグ2"},
				)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(originalPost, nil).Times(1)

				mockPostRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(updatedPost, nil).Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *PatchPostOutput) {
				assert.Len(t, output.Tags, 2)
				assert.Contains(t, output.Tags, valueobject.TagName("新タグ1"))
				assert.Contains(t, output.Tags, valueobject.TagName("新タグ2"))
			},
		},
		{
			name: "正常系: 複数フィールド同時更新",
			input: &PatchPostInput{
				ID:      valueobject.NewPostID(),
				Title:   func() *valueobject.PostTitle { t := valueobject.PostTitle("全て新しいタイトル"); return &t }(),
				Content: func() *valueobject.PostContent { c := valueobject.PostContent("全て新しい内容"); return &c }(),
				Status:  func() *valueobject.PostStatus { s := valueobject.StatusPublished; return &s }(),
				Tags:    []valueobject.TagName{"全て新しいタグ"},
			},
			setupMock: func() {
				originalTitle, _ := valueobject.NewPostTitle("古いタイトル")
				originalContent, _ := valueobject.NewPostContent("古い内容")
				userID := valueobject.NewUserID()
				now := time.Now()
				
				originalPost := entity.ParsePost(
					valueobject.NewPostID(),
					originalTitle,
					originalContent,
					userID,
					valueobject.StatusDraft,
					now,
					now,
					nil,
					&now,
					[]valueobject.TagName{"古いタグ"},
				)

				updatedPost := entity.ParsePost(
					originalPost.ID,
					valueobject.PostTitle("全て新しいタイトル"),
					valueobject.PostContent("全て新しい内容"),
					userID,
					valueobject.StatusPublished,
					now,
					now.Add(time.Minute),
					&now,
					&now,
					[]valueobject.TagName{"全て新しいタグ"},
				)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(originalPost, nil).Times(1)

				mockPostRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(updatedPost, nil).Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *PatchPostOutput) {
				assert.Equal(t, "全て新しいタイトル", string(output.Title))
				assert.Equal(t, "全て新しい内容", string(output.Content))
				assert.Equal(t, valueobject.StatusPublished, output.Status)
				assert.Len(t, output.Tags, 1)
				assert.Contains(t, output.Tags, valueobject.TagName("全て新しいタグ"))
			},
		},
		{
			name: "異常系: 投稿が存在しない",
			input: &PatchPostInput{
				ID:    valueobject.NewPostID(),
				Title: func() *valueobject.PostTitle { t := valueobject.PostTitle("存在しない"); return &t }(),
			},
			setupMock: func() {
				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "post not found"))
			},
			wantErr: true,
			errCode: valueobject.NotFoundCode,
		},
		{
			name: "異常系: 不正なステータス遷移（公開→下書き）",
			input: &PatchPostInput{
				ID:     valueobject.NewPostID(),
				Status: func() *valueobject.PostStatus { s := valueobject.StatusDraft; return &s }(),
			},
			setupMock: func() {
				originalTitle, _ := valueobject.NewPostTitle("公開中の投稿")
				originalContent, _ := valueobject.NewPostContent("内容")
				userID := valueobject.NewUserID()
				now := time.Now()
				
				publishedPost := entity.ParsePost(
					valueobject.NewPostID(),
					originalTitle,
					originalContent,
					userID,
					valueobject.StatusPublished, // 既に公開済み
					now,
					now,
					&now,
					&now,
					[]valueobject.TagName{},
				)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(publishedPost, nil)
			},
			wantErr: true,
			errCode: valueobject.InvalidCode,
		},
		{
			name: "異常系: 投稿更新失敗",
			input: &PatchPostInput{
				ID:    valueobject.NewPostID(),
				Title: func() *valueobject.PostTitle { t := valueobject.PostTitle("更新失敗"); return &t }(),
			},
			setupMock: func() {
				originalTitle, _ := valueobject.NewPostTitle("元のタイトル")
				originalContent, _ := valueobject.NewPostContent("元の内容")
				userID := valueobject.NewUserID()
				now := time.Now()
				
				originalPost := entity.ParsePost(
					valueobject.NewPostID(),
					originalTitle,
					originalContent,
					userID,
					valueobject.StatusDraft,
					now,
					now,
					nil,
					&now,
					[]valueobject.TagName{},
				)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(originalPost, nil)

				mockPostRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to update post"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
		{
			name: "異常系: 更新後投稿取得失敗",
			input: &PatchPostInput{
				ID:    valueobject.NewPostID(),
				Title: func() *valueobject.PostTitle { t := valueobject.PostTitle("取得失敗"); return &t }(),
			},
			setupMock: func() {
				originalTitle, _ := valueobject.NewPostTitle("元のタイトル")
				originalContent, _ := valueobject.NewPostContent("元の内容")
				userID := valueobject.NewUserID()
				now := time.Now()
				
				originalPost := entity.ParsePost(
					valueobject.NewPostID(),
					originalTitle,
					originalContent,
					userID,
					valueobject.StatusDraft,
					now,
					now,
					nil,
					&now,
					[]valueobject.TagName{},
				)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(originalPost, nil).Times(1)

				mockPostRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to get post")).Times(1)
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			usecase := NewPatchPostUsecase(mockPostRepo)
			output, err := usecase.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				myErr, ok := err.(*valueobject.MyError)
				assert.True(t, ok, "エラーはMyError型である必要があります")
				assert.Equal(t, tt.errCode, myErr.Code)
				assert.Nil(t, output)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				if tt.checkResult != nil {
					tt.checkResult(t, output)
				}
			}
		})
	}
}