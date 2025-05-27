package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

// Mock implementation of UserRepository for testing
type mockUserRepository struct {
	users map[string]*entity.User
	shouldFailOnFind bool
	shouldFailOnCreate bool
	findError error
	createError error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	if m.shouldFailOnFind {
		if m.findError != nil {
			return nil, m.findError
		}
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Database error")
	}
	
	user, exists := m.users[email.String()]
	if !exists {
		return nil, valueobject.NewMyError(valueobject.NotFoundCode, "User not found")
	}
	return user, nil
}

func (m *mockUserRepository) Create(ctx context.Context, user *entity.User) error {
	if m.shouldFailOnCreate {
		if m.createError != nil {
			return m.createError
		}
		return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Create failed")
	}
	
	m.users[user.Email.String()] = user
	return nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error) {
	// Not used in register use case, but needed to implement interface
	return nil, valueobject.NewMyError(valueobject.NotFoundCode, "Not implemented")
}

func TestNewRegisterUserUsecase(t *testing.T) {
	repo := newMockUserRepository()
	usecase := NewRegisterUserUsecase(repo)
	
	if usecase.userRepository != repo {
		t.Error("NewRegisterUserUsecase() should set userRepository")
	}
}

func TestRegisterUserUsecase_Execute_Success(t *testing.T) {
	repo := newMockUserRepository()
	usecase := NewRegisterUserUsecase(repo)
	ctx := context.Background()
	
	input := &RegisterUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	
	output, err := usecase.Execute(ctx, input)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
		return
	}
	
	if output == nil {
		t.Error("Execute() should return output")
		return
	}
	
	// Check output fields
	if output.Name != input.Name {
		t.Errorf("Execute() Name = %v, want %v", output.Name, input.Name)
	}
	if output.Email != input.Email {
		t.Errorf("Execute() Email = %v, want %v", output.Email, input.Email)
	}
	if output.ID == "" {
		t.Error("Execute() should generate ID")
	}
	
	// Check user was created in repository
	email, _ := valueobject.NewEmail(input.Email)
	user, err := repo.FindByEmail(ctx, email)
	if err != nil {
		t.Errorf("User should be created in repository: %v", err)
	}
	if user.Name != input.Name {
		t.Errorf("Created user Name = %v, want %v", user.Name, input.Name)
	}
}

func TestRegisterUserUsecase_Execute_InvalidEmail(t *testing.T) {
	repo := newMockUserRepository()
	usecase := NewRegisterUserUsecase(repo)
	ctx := context.Background()
	
	input := &RegisterUserInput{
		Name:     "John Doe",
		Email:    "invalid-email",
		Password: "password123",
	}
	
	output, err := usecase.Execute(ctx, input)
	if err == nil {
		t.Error("Execute() should return error for invalid email")
	}
	if output != nil {
		t.Error("Execute() should not return output on error")
	}
	
	// Check that error is MyError with InvalidCode
	var myErr *valueobject.MyError
	if !errors.As(err, &myErr) {
		t.Error("Error should be MyError")
	} else if myErr.Code != valueobject.InvalidCode {
		t.Errorf("Error code = %v, want %v", myErr.Code, valueobject.InvalidCode)
	}
}

func TestRegisterUserUsecase_Execute_UserAlreadyExists(t *testing.T) {
	repo := newMockUserRepository()
	usecase := NewRegisterUserUsecase(repo)
	ctx := context.Background()
	
	// Create existing user
	email, _ := valueobject.NewEmail("existing@example.com")
	existingUser, _ := entity.NewUser("Existing User", email, "password123")
	repo.users[email.String()] = existingUser
	
	input := &RegisterUserInput{
		Name:     "John Doe",
		Email:    "existing@example.com",
		Password: "password123",
	}
	
	output, err := usecase.Execute(ctx, input)
	if err == nil {
		t.Error("Execute() should return error when user exists")
	}
	if output != nil {
		t.Error("Execute() should not return output on error")
	}
	
	// Check that error is MyError with ConflictCode
	var myErr *valueobject.MyError
	if !errors.As(err, &myErr) {
		t.Error("Error should be MyError")
	} else if myErr.Code != valueobject.ConflictCode {
		t.Errorf("Error code = %v, want %v", myErr.Code, valueobject.ConflictCode)
	}
}

func TestRegisterUserUsecase_Execute_FindUserError(t *testing.T) {
	repo := newMockUserRepository()
	repo.shouldFailOnFind = true
	repo.findError = valueobject.NewMyError(valueobject.InternalServerErrorCode, "Database connection failed")
	
	usecase := NewRegisterUserUsecase(repo)
	ctx := context.Background()
	
	input := &RegisterUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	
	output, err := usecase.Execute(ctx, input)
	if err == nil {
		t.Error("Execute() should return error when repository fails")
	}
	if output != nil {
		t.Error("Execute() should not return output on error")
	}
}

func TestRegisterUserUsecase_Execute_CreateUserError(t *testing.T) {
	repo := newMockUserRepository()
	repo.shouldFailOnCreate = true
	repo.createError = valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to create user")
	
	usecase := NewRegisterUserUsecase(repo)
	ctx := context.Background()
	
	input := &RegisterUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	
	output, err := usecase.Execute(ctx, input)
	if err == nil {
		t.Error("Execute() should return error when create fails")
	}
	if output != nil {
		t.Error("Execute() should not return output on error")
	}
}

func TestRegisterUserUsecase_Execute_InvalidUserData(t *testing.T) {
	repo := newMockUserRepository()
	usecase := NewRegisterUserUsecase(repo)
	ctx := context.Background()
	
	tests := []struct {
		name  string
		input *RegisterUserInput
	}{
		{
			name: "Empty name",
			input: &RegisterUserInput{
				Name:     "",
				Email:    "john@example.com",
				Password: "password123",
			},
		},
		{
			name: "Short password",
			input: &RegisterUserInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "short",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := usecase.Execute(ctx, tt.input)
			if err == nil {
				t.Error("Execute() should return error for invalid user data")
			}
			if output != nil {
				t.Error("Execute() should not return output on error")
			}
		})
	}
}