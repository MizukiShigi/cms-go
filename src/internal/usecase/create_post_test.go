package usecase

import (
	"context"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	mock_repository "github.com/MizukiShigi/cms-go/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionManager := mock_repository.NewMockTransactionManager(ctrl)
	mockPostRepo := mock_repository.NewMockPostRepository(ctrl)
	mockTagRepo := mock_repository.NewMockTagRepository(ctrl)

	tests := []struct {
		name      string
		input     *CreatePostInput
		setupMock func()
		wantErr   bool
		errCode   string
		checkResult func(t *testing.T, output *CreatePostOutput)
	}{
		{
			name: "正常系: 投稿作成成功（タグあり）",
			input: &CreatePostInput{
				Title:   valueobject.PostTitle("テスト投稿タイトル"),
				Content: valueobject.PostContent("テスト投稿の内容です。"),
				Tags:    []valueobject.TagName{"タグ1", "タグ2"},
				UserID:  valueobject.NewUserID(),
				Status:  valueobject.StatusDraft,
			},
			setupMock: func() {
				// トランザクション実行をモック
				mockTransactionManager.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				// 投稿作成をモック
				mockPostRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)

				// タグ作成をモック（2つのタグ）
				tag1 := entity.NewTagWithName("タグ1")
				tag2 := entity.NewTagWithName("タグ2")
				
				mockTagRepo.EXPECT().
					FindOrCreateByName(gomock.Any(), gomock.Any()).
					Return(tag1, nil).Times(1)
				
				mockTagRepo.EXPECT().
					FindOrCreateByName(gomock.Any(), gomock.Any()).
					Return(tag2, nil).Times(1)

				// タグ設定をモック
				mockPostRepo.EXPECT().
					SetTags(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *CreatePostOutput) {
				assert.NotEmpty(t, output.ID)
				assert.Equal(t, "テスト投稿タイトル", string(output.Title))
				assert.Equal(t, "テスト投稿の内容です。", string(output.Content))
				assert.Len(t, output.Tags, 2)
				assert.Contains(t, output.Tags, valueobject.TagName("タグ1"))
				assert.Contains(t, output.Tags, valueobject.TagName("タグ2"))
			},
		},
		{
			name: "正常系: 投稿作成成功（タグなし）",
			input: &CreatePostInput{
				Title:   valueobject.PostTitle("タグなしの投稿"),
				Content: valueobject.PostContent("タグがない投稿です。"),
				Tags:    []valueobject.TagName{},
				UserID:  valueobject.NewUserID(),
				Status:  valueobject.StatusDraft,
			},
			setupMock: func() {
				mockTransactionManager.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				mockPostRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)

				// タグがないのでSetTagsが空配列で呼ばれる
				mockPostRepo.EXPECT().
					SetTags(gomock.Any(), gomock.Any(), gomock.Eq([]*entity.Tag{})).
					Return(nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *CreatePostOutput) {
				assert.NotEmpty(t, output.ID)
				assert.Equal(t, "タグなしの投稿", string(output.Title))
				assert.Equal(t, "タグがない投稿です。", string(output.Content))
				assert.Len(t, output.Tags, 0)
			},
		},
		{
			name: "異常系: 投稿作成失敗",
			input: &CreatePostInput{
				Title:   valueobject.PostTitle("テスト投稿"),
				Content: valueobject.PostContent("内容"),
				Tags:    []valueobject.TagName{},
				UserID:  valueobject.NewUserID(),
				Status:  valueobject.StatusDraft,
			},
			setupMock: func() {
				mockTransactionManager.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				mockPostRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to create post"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
		{
			name: "異常系: タグ作成失敗",
			input: &CreatePostInput{
				Title:   valueobject.PostTitle("テスト投稿"),
				Content: valueobject.PostContent("内容"),
				Tags:    []valueobject.TagName{"失敗タグ"},
				UserID:  valueobject.NewUserID(),
				Status:  valueobject.StatusDraft,
			},
			setupMock: func() {
				mockTransactionManager.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				mockPostRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)

				mockTagRepo.EXPECT().
					FindOrCreateByName(gomock.Any(), gomock.Any()).
					Return(nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to create tag"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
		{
			name: "異常系: タグ設定失敗",
			input: &CreatePostInput{
				Title:   valueobject.PostTitle("テスト投稿"),
				Content: valueobject.PostContent("内容"),
				Tags:    []valueobject.TagName{"テストタグ"},
				UserID:  valueobject.NewUserID(),
				Status:  valueobject.StatusDraft,
			},
			setupMock: func() {
				mockTransactionManager.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				mockPostRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)

				tag := entity.NewTagWithName("テストタグ")
				mockTagRepo.EXPECT().
					FindOrCreateByName(gomock.Any(), gomock.Any()).
					Return(tag, nil)

				mockPostRepo.EXPECT().
					SetTags(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to set tags"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
		{
			name: "異常系: トランザクション失敗",
			input: &CreatePostInput{
				Title:   valueobject.PostTitle("テスト投稿"),
				Content: valueobject.PostContent("内容"),
				Tags:    []valueobject.TagName{},
				UserID:  valueobject.NewUserID(),
				Status:  valueobject.StatusDraft,
			},
			setupMock: func() {
				mockTransactionManager.EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Transaction failed"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			usecase := NewCreatePostUsecase(mockTransactionManager, mockPostRepo, mockTagRepo)
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