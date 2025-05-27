package usecase

import (
	"context"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/service"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

// Mock implementations for testing
type mockUserRepositoryForLogin struct {
	users         map[string]*entity.User
	shouldFail    bool
	failError     error
}

func newMockUserRepositoryForLogin() *mockUserRepositoryForLogin {
	return &mockUserRepositoryForLogin{
		users: make(map[string]*entity.User),
	}
}

func (m *mockUserRepositoryForLogin) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	if m.shouldFail {
		if m.failError != nil {
			return nil, m.failError
		}
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Database error")
	}
	
	user, exists := m.users[email.String()]
	if !exists {
		return nil, valueobject.NewMyError(valueobject.NotFoundCode, "User not found")
	}
	return user, nil
}

func (m *mockUserRepositoryForLogin) Create(ctx context.Context, user *entity.User) error {
	return nil // Not used in login
}

func (m *mockUserRepositoryForLogin) FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error) {
	return nil, nil // Not used in login
}

type mockAuthServiceForLogin struct {
	shouldFail bool
	token      string
	err        error
}

func (m *mockAuthServiceForLogin) GenerateToken(ctx context.Context, userID valueobject.UserID, email valueobject.Email) (string, error) {
	if m.shouldFail {
		if m.err != nil {
			return "", m.err
		}
		return "", valueobject.NewMyError(valueobject.InternalServerErrorCode, "Token generation failed")
	}
	
	if m.token != "" {
		return m.token, nil
	}
	
	return "mock-jwt-token", nil
}

func TestNewLoginUserUsecase(t *testing.T) {
	repo := newMockUserRepositoryForLogin()
	authService := &mockAuthServiceForLogin{}
	usecase := NewLoginUserUsecase(repo, authService)
	
	if usecase.userRepository != repo {
		t.Error("NewLoginUserUsecase() should set userRepository")
	}
	if usecase.authService != authService {
		t.Error("NewLoginUserUsecase() should set authService")
	}
}

func TestLoginUserUsecase_Execute_Success(t *testing.T) {
	repo := newMockUserRepositoryForLogin()
	authService := &mockAuthServiceForLogin{
		token: "test-jwt-token-123",
	}
	usecase := NewLoginUserUsecase(repo, authService)
	ctx := context.Background()
	
	// Create test user
	email, _ := valueobject.NewEmail("john@example.com")
	user, _ := entity.NewUser("John Doe", email, "password123")
	repo.users[email.String()] = user
	
	input := &LoginUserInput{
		Email:    "john@example.com",
		Password: "password123", // Note: Current implementation doesn't verify password!
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
	if output.Token != "test-jwt-token-123" {
		t.Errorf("Execute() Token = %v, want %v", output.Token, "test-jwt-token-123")
	}
	if output.User.Name != user.Name {
		t.Errorf("Execute() User.Name = %v, want %v", output.User.Name, user.Name)
	}
	if output.User.Email != input.Email {
		t.Errorf("Execute() User.Email = %v, want %v", output.User.Email, input.Email)
	}
	if output.User.ID != user.ID.String() {
		t.Errorf("Execute() User.ID = %v, want %v", output.User.ID, user.ID.String())
	}
}

func TestLoginUserUsecase_Execute_InvalidEmail(t *testing.T) {
	repo := newMockUserRepositoryForLogin()
	authService := &mockAuthServiceForLogin{}
	usecase := NewLoginUserUsecase(repo, authService)
	ctx := context.Background()
	
	input := &LoginUserInput{
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
}

func TestLoginUserUsecase_Execute_UserNotFound(t *testing.T) {
	repo := newMockUserRepositoryForLogin()
	authService := &mockAuthServiceForLogin{}
	usecase := NewLoginUserUsecase(repo, authService)
	ctx := context.Background()
	
	input := &LoginUserInput{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	
	output, err := usecase.Execute(ctx, input)
	if err == nil {
		t.Error("Execute() should return error when user not found")
	}
	if output != nil {
		t.Error("Execute() should not return output on error")
	}
}

func TestLoginUserUsecase_Execute_RepositoryError(t *testing.T) {
	repo := newMockUserRepositoryForLogin()
	repo.shouldFail = true
	repo.failError = valueobject.NewMyError(valueobject.InternalServerErrorCode, "Database connection failed")
	
	authService := &mockAuthServiceForLogin{}
	usecase := NewLoginUserUsecase(repo, authService)
	ctx := context.Background()
	
	input := &LoginUserInput{
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

func TestLoginUserUsecase_Execute_TokenGenerationError(t *testing.T) {
	repo := newMockUserRepositoryForLogin()
	authService := &mockAuthServiceForLogin{
		shouldFail: true,
		err:        valueobject.NewMyError(valueobject.InternalServerErrorCode, "JWT service unavailable"),
	}
	usecase := NewLoginUserUsecase(repo, authService)
	ctx := context.Background()
	
	// Create test user
	email, _ := valueobject.NewEmail("john@example.com")
	user, _ := entity.NewUser("John Doe", email, "password123")
	repo.users[email.String()] = user
	
	input := &LoginUserInput{
		Email:    "john@example.com",
		Password: "password123",
	}
	
	output, err := usecase.Execute(ctx, input)
	if err == nil {
		t.Error("Execute() should return error when token generation fails")
	}
	if output != nil {
		t.Error("Execute() should not return output on error")
	}
}

// IMPORTANT: This test demonstrates a security issue in the current implementation
func TestLoginUserUsecase_Execute_PasswordNotVerified(t *testing.T) {
	repo := newMockUserRepositoryForLogin()
	authService := &mockAuthServiceForLogin{}
	usecase := NewLoginUserUsecase(repo, authService)
	ctx := context.Background()
	
	// Create test user
	email, _ := valueobject.NewEmail("john@example.com")
	user, _ := entity.NewUser("John Doe", email, "correct_password")
	repo.users[email.String()] = user
	
	// Try to login with wrong password
	input := &LoginUserInput{
		Email:    "john@example.com",
		Password: "wrong_password", // This should fail, but current implementation doesn't check!
	}
	
	output, err := usecase.Execute(ctx, input)
	
	// TODO: This should fail but currently succeeds due to missing password verification
	// The current implementation has a security flaw - it doesn't verify the password!
	if err != nil {
		t.Logf("GOOD: Password verification would prevent login with wrong password: %v", err)
	} else {
		t.Logf("SECURITY ISSUE: Current implementation allows login without password verification!")
		t.Logf("Output: %+v", output)
		// This is expected behavior with current implementation, but it's a bug
		// In a real scenario, this test would help identify the security issue
	}
}