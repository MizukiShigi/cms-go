package usecase

import (
	"context"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	repositoryMock "github.com/MizukiShigi/cms-go/mocks/repository"
	serviceMock "github.com/MizukiShigi/cms-go/mocks/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoginUserUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repositoryMock.NewMockUserRepository(ctrl)
	mockAuthService := serviceMock.NewMockAuthService(ctrl)

	t.Run("ログインが成功する", func(t *testing.T) {
		usecase := NewLoginUserUsecase(mockUserRepo, mockAuthService)
		input := &LoginUserInput{
			Email:    "test@example.com",
			Password: "password123",
		}

		email, _ := valueobject.NewEmail("test@example.com")
		user, _ := entity.NewUser("テストユーザー", email, "password123")

		// ユーザーを正常に取得
		mockUserRepo.EXPECT().FindByEmail(context.Background(), email).Return(user, nil)
		
		// トークン生成が成功
		mockAuthService.EXPECT().GenerateToken(context.Background(), user.ID, user.Email).
			Return("test-token", nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "test-token", output.Token)
		assert.Equal(t, user.ID.String(), output.User.ID)
		assert.Equal(t, "テストユーザー", output.User.Name)
		assert.Equal(t, "test@example.com", output.User.Email)
	})

	t.Run("無効なメールアドレスでエラーが発生する", func(t *testing.T) {
		usecase := NewLoginUserUsecase(mockUserRepo, mockAuthService)
		input := &LoginUserInput{
			Email:    "invalid-email",
			Password: "password123",
		}

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("ユーザーが存在しない場合にエラーが発生する", func(t *testing.T) {
		usecase := NewLoginUserUsecase(mockUserRepo, mockAuthService)
		input := &LoginUserInput{
			Email:    "notfound@example.com",
			Password: "password123",
		}

		email, _ := valueobject.NewEmail("notfound@example.com")

		// ユーザーが見つからない
		mockUserRepo.EXPECT().FindByEmail(context.Background(), email).
			Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "User not found"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("ユーザーがnilの場合にエラーが発生する", func(t *testing.T) {
		usecase := NewLoginUserUsecase(mockUserRepo, mockAuthService)
		input := &LoginUserInput{
			Email:    "test@example.com",
			Password: "password123",
		}

		email, _ := valueobject.NewEmail("test@example.com")

		// ユーザーがnilで返される
		mockUserRepo.EXPECT().FindByEmail(context.Background(), email).Return(nil, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
		
		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.NotFoundCode, myErr.Code)
	})

	t.Run("パスワードが一致しない場合にエラーが発生する", func(t *testing.T) {
		usecase := NewLoginUserUsecase(mockUserRepo, mockAuthService)
		input := &LoginUserInput{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		email, _ := valueobject.NewEmail("test@example.com")
		user, _ := entity.NewUser("テストユーザー", email, "password123")

		// ユーザーを正常に取得
		mockUserRepo.EXPECT().FindByEmail(context.Background(), email).Return(user, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
		
		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.UnauthorizedCode, myErr.Code)
	})

	t.Run("トークン生成が失敗する場合にエラーが発生する", func(t *testing.T) {
		usecase := NewLoginUserUsecase(mockUserRepo, mockAuthService)
		input := &LoginUserInput{
			Email:    "test@example.com",
			Password: "password123",
		}

		email, _ := valueobject.NewEmail("test@example.com")
		user, _ := entity.NewUser("テストユーザー", email, "password123")

		// ユーザーを正常に取得
		mockUserRepo.EXPECT().FindByEmail(context.Background(), email).Return(user, nil)
		
		// トークン生成が失敗
		mockAuthService.EXPECT().GenerateToken(context.Background(), user.ID, user.Email).
			Return("", valueobject.NewMyError(valueobject.InternalServerErrorCode, "Token generation failed"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("リポジトリエラーでエラーが発生する", func(t *testing.T) {
		usecase := NewLoginUserUsecase(mockUserRepo, mockAuthService)
		input := &LoginUserInput{
			Email:    "test@example.com",
			Password: "password123",
		}

		email, _ := valueobject.NewEmail("test@example.com")

		// リポジトリでエラーが発生
		mockUserRepo.EXPECT().FindByEmail(context.Background(), email).
			Return(nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Database error"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}