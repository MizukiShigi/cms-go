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

func TestGetPostUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock_repository.NewMockPostRepository(ctrl)

	tests := []struct {
		name      string
		input     *GetPostInput
		setupMock func()
		wantErr   bool
		errCode   string
		checkResult func(t *testing.T, output *GetPostOutput)
	}{
		{
			name: "正常系: 投稿取得成功（公開済み）",
			input: &GetPostInput{
				ID: valueobject.NewPostID(),
			},
			setupMock: func() {
				title, _ := valueobject.NewPostTitle("テスト投稿タイトル")
				content, _ := valueobject.NewPostContent("テスト投稿の内容です。")
				userID := valueobject.NewUserID()
				now := time.Now()
				
				post := entity.ParsePost(
					valueobject.NewPostID(),
					title,
					content,
					userID,
					valueobject.StatusPublished,
					now,
					now,
					&now,
					&now,
					[]valueobject.TagName{"タグ1", "タグ2"},
				)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(post, nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *GetPostOutput) {
				assert.NotEmpty(t, output.ID)
				assert.Equal(t, "テスト投稿タイトル", string(output.Title))
				assert.Equal(t, "テスト投稿の内容です。", string(output.Content))
				assert.Equal(t, valueobject.StatusPublished, output.Status)
				assert.Len(t, output.Tags, 2)
				assert.Contains(t, output.Tags, valueobject.TagName("タグ1"))
				assert.Contains(t, output.Tags, valueobject.TagName("タグ2"))
				assert.NotNil(t, output.FirstPublishedAt)
				assert.NotNil(t, output.ContentUpdatedAt)
			},
		},
		{
			name: "正常系: 投稿取得成功（下書き、タグなし）",
			input: &GetPostInput{
				ID: valueobject.NewPostID(),
			},
			setupMock: func() {
				title, _ := valueobject.NewPostTitle("下書き投稿")
				content, _ := valueobject.NewPostContent("まだ公開していない投稿です。")
				userID := valueobject.NewUserID()
				now := time.Now()
				
				post := entity.ParsePost(
					valueobject.NewPostID(),
					title,
					content,
					userID,
					valueobject.StatusDraft,
					now,
					now,
					nil, // 下書きなので初回公開日時はnull
					&now,
					[]valueobject.TagName{}, // タグなし
				)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(post, nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *GetPostOutput) {
				assert.NotEmpty(t, output.ID)
				assert.Equal(t, "下書き投稿", string(output.Title))
				assert.Equal(t, "まだ公開していない投稿です。", string(output.Content))
				assert.Equal(t, valueobject.StatusDraft, output.Status)
				assert.Len(t, output.Tags, 0)
				assert.Nil(t, output.FirstPublishedAt)
				assert.NotNil(t, output.ContentUpdatedAt)
			},
		},
		{
			name: "異常系: 投稿が存在しない",
			input: &GetPostInput{
				ID: valueobject.NewPostID(),
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
			name: "異常系: リポジトリエラー",
			input: &GetPostInput{
				ID: valueobject.NewPostID(),
			},
			setupMock: func() {
				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "database error"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			usecase := NewGetPostUsecase(mockPostRepo)
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