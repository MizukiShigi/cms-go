package usecase

import (
	"context"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	mock_repository "github.com/MizukiShigi/cms-go/mocks/repository"
	mock_service "github.com/MizukiShigi/cms-go/mocks/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLoginUserUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	tests := []struct {
		name      string
		input     *LoginUserInput
		setupMock func()
		wantErr   bool
		errCode   string
		checkResult func(t *testing.T, output *LoginUserOutput)
	}{
		{
			name: "正常系: ログイン成功",
			input: &LoginUserInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				email, _ := valueobject.NewEmail("test@example.com")
				user, _ := entity.NewUser("テストユーザー", email, "password123")
				
				mockUserRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(user, nil)
				
				mockAuthService.EXPECT().
					GenerateToken(gomock.Any(), user.ID, user.Email).
					Return("test-jwt-token", nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, output *LoginUserOutput) {
				assert.Equal(t, "test-jwt-token", output.Token)
				assert.NotEmpty(t, output.User.ID)
				assert.Equal(t, "テストユーザー", output.User.Name)
				assert.Equal(t, "test@example.com", output.User.Email)
			},
		},
		{
			name: "異常系: 無効なメールアドレス",
			input: &LoginUserInput{
				Email:    "invalid-email",
				Password: "password123",
			},
			setupMock: func() {
				// モックは呼ばれない（Email作成でエラー）
			},
			wantErr: true,
			errCode: valueobject.InvalidCode,
		},
		{
			name: "異常系: ユーザーが存在しない",
			input: &LoginUserInput{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			setupMock: func() {
				email, _ := valueobject.NewEmail("notfound@example.com")
				mockUserRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "user not found"))
			},
			wantErr: true,
			errCode: valueobject.NotFoundCode,
		},
		{
			name: "異常系: FindByEmailがnilユーザーを返す",
			input: &LoginUserInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				email, _ := valueobject.NewEmail("test@example.com")
				mockUserRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(nil, nil)
			},
			wantErr: true,
			errCode: valueobject.NotFoundCode,
		},
		{
			name: "異常系: パスワードが間違っている",
			input: &LoginUserInput{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func() {
				email, _ := valueobject.NewEmail("test@example.com")
				user, _ := entity.NewUser("テストユーザー", email, "correctpassword")
				
				mockUserRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(user, nil)
			},
			wantErr: true,
			errCode: valueobject.UnauthorizedCode,
		},
		{
			name: "異常系: トークン生成失敗",
			input: &LoginUserInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				email, _ := valueobject.NewEmail("test@example.com")
				user, _ := entity.NewUser("テストユーザー", email, "password123")
				
				mockUserRepo.EXPECT().
					FindByEmail(gomock.Any(), email).
					Return(user, nil)
				
				mockAuthService.EXPECT().
					GenerateToken(gomock.Any(), user.ID, user.Email).
					Return("", valueobject.NewMyError(valueobject.InternalServerErrorCode, "failed to generate token"))
			},
			wantErr: true,
			errCode: valueobject.InternalServerErrorCode,
		},
		{
			name: "異常系: リポジトリエラー",
			input: &LoginUserInput{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			usecase := NewLoginUserUsecase(mockUserRepo, mockAuthService)
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