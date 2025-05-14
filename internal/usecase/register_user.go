package usecase

import (
	"context"
	"errors"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/myerror"
	"github.com/MizukiShigi/cms-go/internal/domain/repository"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type RegisterUserInput struct {
	Name     string
	Email    string
	Password string
}

type RegisterUserOutput struct {
	ID    string
	Name  string
	Email string
}

type RegisterUserUsecase struct {
	userRepository repository.UserRepository
}

func NewRegisterUserUsecase(userRepository repository.UserRepository) *RegisterUserUsecase {
	return &RegisterUserUsecase{
		userRepository: userRepository,
	}
}

func (u *RegisterUserUsecase) Execute(ctx context.Context, input *RegisterUserInput) (*RegisterUserOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InvalidCode, "Invalid email")
	}

	existingUser, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		var myErr *myerror.MyError
		if !errors.As(err, &myErr) {
			return nil, myerror.NewMyError(myerror.InternalServerErrorCode, "Failed to find user")
		}
		if myErr.Code != myerror.NotFoundCode {
			// 独自エラーのユーザーが存在しない場合のエラーのみ登録処理を続行
			return nil, err
		}
	}
	if existingUser != nil {
		return nil, myerror.NewMyError(myerror.ConflictCode, "User already exists")
	}

	user, err := entity.NewUser(input.Name, email, input.Password)
	if err != nil {
		return nil, err
	}

	err = u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &RegisterUserOutput{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email.String(),
	}, nil
}
