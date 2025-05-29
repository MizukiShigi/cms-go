package usecase

import (
	"context"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/MizukiShigi/cms-go/mocks/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegisterUserUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repository.NewMockUserRepository(ctrl)

	t.Run("新規ユーザー登録が成功する", func(t *testing.T) {
		usecase := NewRegisterUserUsecase(mockUserRepo)
		input := &RegisterUserInput{
			Name:     "テストユーザー",
			Email:    "test@example.com",
			Password: "password123",
		}

		email, _ := valueobject.NewEmail("test@example.com")

		// ユーザーが存在しないことを確認
		mockUserRepo.EXPECT().FindByEmail(context.Background(), email).
			Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "User not found"))

		// ユーザー作成が成功
		mockUserRepo.EXPECT().Create(context.Background(), gomock.Any()).Return(nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "テストユーザー", output.Name)
		assert.Equal(t, "test@example.com", output.Email)
		assert.NotEmpty(t, output.ID)
	})

	t.Run("無効なメールアドレスでエラーが発生する", func(t *testing.T) {
		usecase := NewRegisterUserUsecase(mockUserRepo)
		input := &RegisterUserInput{
			Name:     "テストユーザー",
			Email:    "invalid-email",
			Password: "password123",
		}

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)

		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.InvalidCode, myErr.Code)
	})

	t.Run("既存ユーザーが存在する場合にエラーが発生する", func(t *testing.T) {
		usecase := NewRegisterUserUsecase(mockUserRepo)
		input := &RegisterUserInput{
			Name:     "テストユーザー",
			Email:    "test@example.com",
			Password: "password123",
		}

		email, _ := valueobject.NewEmail("test@example.com")
		existingUser, _ := entity.NewUser("既存ユーザー", email, "password")

		// 既存ユーザーが見つかる
		mockUserRepo.EXPECT().FindByEmail(context.Background(), email).Return(existingUser, nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)

		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.ConflictCode, myErr.Code)
	})

	t.Run("短すぎるパスワードでエラーが発生する", func(t *testing.T) {
		usecase := NewRegisterUserUsecase(mockUserRepo)
		input := &RegisterUserInput{
			Name:     "テストユーザー",
			Email:    "test@example.com",
			Password: "123", // 短すぎるパスワード
		}

		email, _ := valueobject.NewEmail("test@example.com")

		// ユーザーが存在しないことを確認
		mockUserRepo.EXPECT().FindByEmail(context.Background(), email).
			Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "User not found"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("リポジトリエラーでエラーが発生する", func(t *testing.T) {
		usecase := NewRegisterUserUsecase(mockUserRepo)
		input := &RegisterUserInput{
			Name:     "テストユーザー",
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

	t.Run("ユーザー作成時にリポジトリエラーでエラーが発生する", func(t *testing.T) {
		usecase := NewRegisterUserUsecase(mockUserRepo)
		input := &RegisterUserInput{
			Name:     "テストユーザー",
			Email:    "test@example.com",
			Password: "password123",
		}

		email, _ := valueobject.NewEmail("test@example.com")

		// ユーザーが存在しないことを確認
		mockUserRepo.EXPECT().FindByEmail(context.Background(), email).
			Return(nil, valueobject.NewMyError(valueobject.NotFoundCode, "User not found"))

		// ユーザー作成でエラーが発生
		mockUserRepo.EXPECT().Create(context.Background(), gomock.Any()).
			Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Create failed"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
