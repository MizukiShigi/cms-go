package usecase

import (
	"context"

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
	UserID string
	Name   string
	Email  string
}

type RegisterUserUsecase struct {
	userRepository repository.UserRepository
}

func NewRegisterUserUsecase(userRepository *repository.UserRepository) *RegisterUserUsecase {
	return &RegisterUserUsecase{
		userRepository: *userRepository,
	}
}

func (u *RegisterUserUsecase) Execute(ctx context.Context, input *RegisterUserInput) (*RegisterUserOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InvalidRequestCode, "Invalid email")
	}

	existingUser, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InternalServerErrorCode, "Failed to find user")
	}
	if existingUser != nil {
		return nil, myerror.NewMyError(myerror.ConflictCode, "User already exists")
	}

	user, err := entity.NewUser(input.Name, email, input.Password)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InternalServerErrorCode, "Failed to create user")
	}

	err = u.userRepository.Save(ctx, user)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InternalServerErrorCode, "Failed to save user")
	}

	return &RegisterUserOutput{
		UserID: user.ID.String(),
		Name:   user.Name,
		Email:  user.Email.String(),
	}, nil
}
