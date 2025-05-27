package usecase

import (
	"context"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

// Mock implementation of UserRepository
type mockUserRepository struct {
	users           map[string]*entity.User
	createError     error
	findByEmailUser *entity.User
	findByEmailError error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *entity.User) error {
	if m.createError != nil {
		return m.createError
	}
	m.users[user.Email.String()] = user
	return nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	if m.findByEmailError != nil {
		return nil, m.findByEmailError
	}
	if m.findByEmailUser != nil {
		return m.findByEmailUser, nil
	}
	if user, exists := m.users[email.String()]; exists {
		return user, nil
	}
	return nil, valueobject.NewMyError(valueobject.NotFoundCode, "User not found")
}

// Helper methods for setting up mock behavior
func (m *mockUserRepository) SetCreateError(err error) {
	m.createError = err
}

func (m *mockUserRepository) SetFindByEmailResult(user *entity.User, err error) {
	m.findByEmailUser = user
	m.findByEmailError = err
}

func TestNewRegisterUserUsecase(t *testing.T) {
	mockRepo := newMockUserRepository()
	usecase := NewRegisterUserUsecase(mockRepo)
	
	if usecase.userRepository != mockRepo {
		t.Errorf("NewRegisterUserUsecase() did not set repository correctly")
	}
}

func TestRegisterUserUsecase_Execute(t *testing.T) {
	tests := []struct {
		name             string
		input            *RegisterUserInput
		setupMock        func(*mockUserRepository)
		wantErr          bool
		wantErrorCode    valueobject.Code
		expectedName     string
		expectedEmail    string
	}{
		{
			name: "Successful registration",
			input: &RegisterUserInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockUserRepository) {
				// User does not exist
				m.SetFindByEmailResult(nil, valueobject.NewMyError(valueobject.NotFoundCode, "User not found"))
			},
			wantErr:       false,
			expectedName:  "John Doe",
			expectedEmail: "john@example.com",
		},
		{
			name: "Invalid email format",
			input: &RegisterUserInput{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "password123",
			},
			setupMock:     func(m *mockUserRepository) {},
			wantErr:       true,
			wantErrorCode: valueobject.InvalidCode,
		},
		{
			name: "User already exists",
			input: &RegisterUserInput{
				Name:     "John Doe",
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockUserRepository) {
				// User exists
				existingEmail, _ := valueobject.NewEmail("existing@example.com")
				existingUser, _ := entity.NewUser("Existing User", existingEmail, "password")
				m.SetFindByEmailResult(existingUser, nil)
			},
			wantErr:       true,
			wantErrorCode: valueobject.ConflictCode,
		},
		{
			name: "Database error on find",
			input: &RegisterUserInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockUserRepository) {
				// Database error
				m.SetFindByEmailResult(nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Database error"))
			},
			wantErr:       true,
			wantErrorCode: valueobject.InternalServerErrorCode,
		},
		{
			name: "Invalid user data - empty name",
			input: &RegisterUserInput{
				Name:     "",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockUserRepository) {
				// User does not exist
				m.SetFindByEmailResult(nil, valueobject.NewMyError(valueobject.NotFoundCode, "User not found"))
			},
			wantErr:       true,
			wantErrorCode: valueobject.InvalidCode,
		},
		{
			name: "Invalid user data - short password",
			input: &RegisterUserInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "short",
			},
			setupMock: func(m *mockUserRepository) {
				// User does not exist
				m.SetFindByEmailResult(nil, valueobject.NewMyError(valueobject.NotFoundCode, "User not found"))
			},
			wantErr:       true,
			wantErrorCode: valueobject.InvalidCode,
		},
		{
			name: "Database error on create",
			input: &RegisterUserInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockUserRepository) {
				// User does not exist
				m.SetFindByEmailResult(nil, valueobject.NewMyError(valueobject.NotFoundCode, "User not found"))
				// But create fails
				m.SetCreateError(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Create failed"))
			},
			wantErr:       true,
			wantErrorCode: valueobject.InternalServerErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			tt.setupMock(mockRepo)
			
			usecase := NewRegisterUserUsecase(mockRepo)
			
			result, err := usecase.Execute(context.Background(), tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("RegisterUserUsecase.Execute() expected error but got none")
					return
				}
				
				if myErr, ok := err.(*valueobject.MyError); ok {
					if myErr.Code != tt.wantErrorCode {
						t.Errorf("RegisterUserUsecase.Execute() error code = %v, want %v", myErr.Code, tt.wantErrorCode)
					}
				} else {
					t.Errorf("RegisterUserUsecase.Execute() error is not MyError type")
				}
				return
			}
			
			if err != nil {
				t.Errorf("RegisterUserUsecase.Execute() unexpected error = %v", err)
				return
			}
			
			// Verify output
			if result.Name != tt.expectedName {
				t.Errorf("RegisterUserUsecase.Execute() name = %v, want %v", result.Name, tt.expectedName)
			}
			
			if result.Email != tt.expectedEmail {
				t.Errorf("RegisterUserUsecase.Execute() email = %v, want %v", result.Email, tt.expectedEmail)
			}
			
			// Verify ID is generated
			if result.ID == "" {
				t.Errorf("RegisterUserUsecase.Execute() ID should be generated")
			}
			
			// Verify user was created in mock repository
			email, _ := valueobject.NewEmail(tt.input.Email)
			createdUser, exists := mockRepo.users[email.String()]
			if !exists {
				t.Errorf("RegisterUserUsecase.Execute() user should be created in repository")
			} else {
				if createdUser.Name != tt.expectedName {
					t.Errorf("Created user name = %v, want %v", createdUser.Name, tt.expectedName)
				}
				if !createdUser.Email.Equals(email) {
					t.Errorf("Created user email = %v, want %v", createdUser.Email, email)
				}
			}
		})
	}
}