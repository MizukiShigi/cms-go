package usecase

import (
	"context"

	"github.com/MizukiShigi/cms-go/internal/domain/myerror"
	"github.com/MizukiShigi/cms-go/internal/domain/repository"
	"github.com/MizukiShigi/cms-go/internal/domain/service"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	Token string
	User  UserOutput
}

type UserOutput struct {
	ID    string
	Name  string
	Email string
}

type LoginUserUsecase struct {
	userRepository repository.UserRepository
	authService    service.AuthService
}

func NewLoginUserUsecase(userRepository repository.UserRepository, authService service.AuthService) *LoginUserUsecase {
	return &LoginUserUsecase{
		userRepository: userRepository,
		authService:    authService,
	}
}

func (u *LoginUserUsecase) Execute(ctx context.Context, input *LoginUserInput) (*LoginUserOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, myerror.NewMyError(myerror.NotFoundCode, "user not found")
	}

	token, err := u.authService.GenerateToken(ctx, user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &LoginUserOutput{
		Token: token,
		User: UserOutput{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email.String(),
		},
	}, nil
}
