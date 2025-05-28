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

func TestUpdatePostUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionManager := mock_repository.NewMockTransactionManager(ctrl)
	mockPostRepo := mock_repository.NewMockPostRepository(ctrl)
	mockTagRepo := mock_repository.NewMockTagRepository(ctrl)

	tests := []struct {
		name      string
		input     *UpdatePostInput
		setupMock func()
		wantErr   bool
		errCode   string
		checkResult func(t *testing.T, output *UpdatePostOutput)
	}{
		{
			name: "正常系: 投稿更新成功（タイトル・内容・タグ変更）",
			input: &UpdatePostInput{
				ID:      valueobject.NewPostID(),
				Title:   valueobject.PostTitle("更新されたタイトル"),
				Content: valueobject.PostContent("更新された内容です。"),
				Tags:    []valueobject.TagName{"新タグ1", "新タグ2"},
				Status:  valueobject.StatusDraft,
			},
			setupMock: func() {
				// 既存の投稿を取得
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
					[]valueobject.TagName{"元タグ"},
				)

				// 更新後の投稿
				updatedPost := entity.ParsePost(
					originalPost.ID,
					valueobject.PostTitle("更新されたタイトル"),
					valueobject.PostContent("更新された内容です。"),
					userID,
					valueobject.StatusDraft,
					now,
					now.Add(time.Minute),
					nil,
					&now,
					[]valueobject.TagName{"新タグ1", "新タグ2"},
				)

				// 最初の投稿取得（更新前）
				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(originalPost, nil).Times(1)

				// トランザクション実行
				mockTransactionManager.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				// 投稿更新
				mockPostRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)

				// タグ作成・取得
				tag1 := entity.NewTagWithName("新タグ1")
				tag2 := entity.NewTagWithName("新タグ2")
				
				mockTagRepo.EXPECT().
					FindOrCreateByName(gomock.Any(), gomock.Any()).
					Return(tag1, nil).Times(1)
				
				mockTagRepo.EXPECT().
					FindOrCreateByName(gomock.Any(), gomock.Any()).
					Return(tag2, nil).Times(1)

				// タグ設定
				mockPostRepo.EXPECT().
					SetTags(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				// 更新後の投稿取得
				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(updatedPost, nil).Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *UpdatePostOutput) {
				assert.NotEmpty(t, output.ID)
				assert.Equal(t, "更新されたタイトル", string(output.Title))
				assert.Equal(t, "更新された内容です。", string(output.Content))
				assert.Equal(t, valueobject.StatusDraft, output.Status)
				assert.Len(t, output.Tags, 2)
				assert.Contains(t, output.Tags, valueobject.TagName("新タグ1"))
				assert.Contains(t, output.Tags, valueobject.TagName("新タグ2"))
			},
		},
		{
			name: "正常系: 投稿更新成功（タグなし）",
			input: &UpdatePostInput{
				ID:      valueobject.NewPostID(),
				Title:   valueobject.PostTitle("タグなし更新"),
				Content: valueobject.PostContent("タグを削除した投稿です。"),
				Tags:    []valueobject.TagName{},
				Status:  valueobject.StatusDraft,
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
					[]valueobject.TagName{"削除予定タグ"},
				)

				updatedPost := entity.ParsePost(
					originalPost.ID,
					valueobject.PostTitle("タグなし更新"),
					valueobject.PostContent("タグを削除した投稿です。"),
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

				mockTransactionManager.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				mockPostRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)

				// タグが空の場合はSetTagsが空配列で呼ばれる
				mockPostRepo.EXPECT().
					SetTags(gomock.Any(), gomock.Any(), gomock.Eq([]*entity.Tag{})).
					Return(nil)

				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(updatedPost, nil).Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *UpdatePostOutput) {
				assert.Equal(t, "タグなし更新", string(output.Title))
				assert.Equal(t, "タグを削除した投稿です。", string(output.Content))
				assert.Len(t, output.Tags, 0)
			},
		},
		{
			name: "異常系: 投稿が存在しない",
			input: &UpdatePostInput{
				ID:      valueobject.NewPostID(),
				Title:   valueobject.PostTitle("存在しない投稿"),
				Content: valueobject.PostContent("内容"),
				Tags:    []valueobject.TagName{},
				Status:  valueobject.StatusDraft,
			},
			setupMock: func() {
				mockPostRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "post not found"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode, // ReturnMyErrorでラップされるため
		},
		{
			name: "異常系: 投稿更新失敗",
			input: &UpdatePostInput{
				ID:      valueobject.NewPostID(),
				Title:   valueobject.PostTitle("更新失敗"),
				Content: valueobject.PostContent("内容"),
				Tags:    []valueobject.TagName{},
				Status:  valueobject.StatusDraft,
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

				mockTransactionManager.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				mockPostRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to update post"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
		{
			name: "異常系: タグ作成失敗",
			input: &UpdatePostInput{
				ID:      valueobject.NewPostID(),
				Title:   valueobject.PostTitle("タグ作成失敗"),
				Content: valueobject.PostContent("内容"),
				Tags:    []valueobject.TagName{"失敗タグ"},
				Status:  valueobject.StatusDraft,
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

				mockTransactionManager.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				mockPostRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)

				mockTagRepo.EXPECT().
					FindOrCreateByName(gomock.Any(), gomock.Any()).
					Return(nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to create tag"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			usecase := NewUpdatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)
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