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

func TestRegisterUserUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)

	tests := []struct {
		name      string
		input     *RegisterUserInput
		setupMock func()
		wantErr   bool
		errCode   string
		checkResult func(t *testing.T, output *RegisterUserOutput)
	}{
		{
			name: "正常系: ユーザー登録成功",
			input: &RegisterUserInput{
				Name:     "テストユーザー",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				email, _ := valueobject.NewEmail("test@example.com")
				// ユーザーが存在しないことを確認
				mockUserRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "user not found"))
				
				// ユーザー作成処理
				mockUserRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *RegisterUserOutput) {
				assert.NotEmpty(t, output.ID)
				assert.Equal(t, "テストユーザー", output.Name)
				assert.Equal(t, "test@example.com", output.Email)
			},
		},
		{
			name: "異常系: 無効なメールアドレス",
			input: &RegisterUserInput{
				Name:     "テストユーザー",
				Email:    "invalid-email",
				Password: "password123",
			},
			setupMock: func() {
				// モックは呼ばれない
			},
			wantErr: true,
			errCode: valueobject.InvalidCode,
		},
		{
			name: "異常系: 既存ユーザーが存在する",
			input: &RegisterUserInput{
				Name:     "テストユーザー",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				email, _ := valueobject.NewEmail("test@example.com")
				existingUser := &entity.User{
					ID:   valueobject.NewUserID(),
					Name: "既存ユーザー",
					Email: email,
				}
				mockUserRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(existingUser, nil)
			},
			wantErr: true,
			errCode: valueobject.ConflictCode,
		},
		{
			name: "異常系: リポジトリエラー（予期しないエラー）",
			input: &RegisterUserInput{
				Name:     "テストユーザー",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				email, _ := valueobject.NewEmail("test@example.com")
				mockUserRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "database error"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
		{
			name: "異常系: ユーザー作成失敗",
			input: &RegisterUserInput{
				Name:     "テストユーザー",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				email, _ := valueobject.NewEmail("test@example.com")
				mockUserRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "user not found"))
				
				mockUserRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "failed to create user"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
		{
			name: "異常系: 短すぎるパスワード",
			input: &RegisterUserInput{
				Name:     "テストユーザー",
				Email:    "test@example.com",
				Password: "short",
			},
			setupMock: func() {
				email, _ := valueobject.NewEmail("test@example.com")
				mockUserRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "user not found"))
				// NewUserでエラーになるのでCreate は呼ばれない
			},
			wantErr: true,
			errCode: valueobject.InvalidCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			usecase := NewRegisterUserUsecase(mockUserRepo)
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