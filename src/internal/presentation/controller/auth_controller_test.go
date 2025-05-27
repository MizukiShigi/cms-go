package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/MizukiShigi/cms-go/internal/usecase"
)

// Mock implementations for usecases
type mockRegisterUserUsecase struct {
	output *usecase.RegisterUserOutput
	error  error
}

func (m *mockRegisterUserUsecase) Execute(ctx context.Context, input *usecase.RegisterUserInput) (*usecase.RegisterUserOutput, error) {
	if m.error != nil {
		return nil, m.error
	}
	return m.output, nil
}

type mockLoginUserUsecase struct {
	output *usecase.LoginUserOutput
	error  error
}

func (m *mockLoginUserUsecase) Execute(ctx context.Context, input *usecase.LoginUserInput) (*usecase.LoginUserOutput, error) {
	if m.error != nil {
		return nil, m.error
	}
	return m.output, nil
}

func TestNewAuthController(t *testing.T) {
	mockRegister := &mockRegisterUserUsecase{}
	mockLogin := &mockLoginUserUsecase{}
	
	controller := NewAuthController(mockRegister, mockLogin)
	
	if controller.registerUserUsecase != mockRegister {
		t.Errorf("NewAuthController() did not set registerUserUsecase correctly")
	}
	if controller.loginUserUsecase != mockLogin {
		t.Errorf("NewAuthController() did not set loginUserUsecase correctly")
	}
}

func TestAuthController_Register(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		setupMock          func(*mockRegisterUserUsecase)
		expectedStatusCode int
		expectError        bool
	}{
		{
			name: "Successful registration",
			requestBody: RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockRegisterUserUsecase) {
				m.output = &usecase.RegisterUserOutput{
					ID:    "user-123",
					Name:  "John Doe",
					Email: "john@example.com",
				}
				m.error = nil
			},
			expectedStatusCode: http.StatusCreated,
			expectError:        false,
		},
		{
			name: "Invalid JSON payload",
			requestBody: "invalid-json",
			setupMock: func(m *mockRegisterUserUsecase) {
				m.output = nil
				m.error = nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectError:        true,
		},
		{
			name: "Missing required fields",
			requestBody: RegisterRequest{
				Name:     "",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockRegisterUserUsecase) {
				m.output = nil
				m.error = nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectError:        true,
		},
		{
			name: "Invalid email format",
			requestBody: RegisterRequest{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "password123",
			},
			setupMock: func(m *mockRegisterUserUsecase) {
				m.output = nil
				m.error = nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectError:        true,
		},
		{
			name: "User already exists",
			requestBody: RegisterRequest{
				Name:     "John Doe",
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockRegisterUserUsecase) {
				m.output = nil
				m.error = valueobject.NewMyError(valueobject.ConflictCode, "User already exists")
			},
			expectedStatusCode: http.StatusConflict,
			expectError:        true,
		},
		{
			name: "Internal server error from usecase",
			requestBody: RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockRegisterUserUsecase) {
				m.output = nil
				m.error = valueobject.NewMyError(valueobject.InternalServerErrorCode, "Database error")
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegister := &mockRegisterUserUsecase{}
			mockLogin := &mockLoginUserUsecase{}
			tt.setupMock(mockRegister)
			
			controller := NewAuthController(mockRegister, mockLogin)
			
			// Create request body
			var bodyBytes []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}
			
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			
			recorder := httptest.NewRecorder()
			controller.Register(recorder, req)
			
			// Check status code
			if recorder.Code != tt.expectedStatusCode {
				t.Errorf("Register() status code = %v, want %v", recorder.Code, tt.expectedStatusCode)
			}
			
			// Check content type
			if recorder.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Register() Content-Type = %v, want application/json", recorder.Header().Get("Content-Type"))
			}
			
			// Check response body
			if !tt.expectError && recorder.Code == http.StatusCreated {
				var response UserResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal success response: %v", err)
				}
				
				if response.ID != mockRegister.output.ID {
					t.Errorf("Response ID = %v, want %v", response.ID, mockRegister.output.ID)
				}
				if response.Name != mockRegister.output.Name {
					t.Errorf("Response Name = %v, want %v", response.Name, mockRegister.output.Name)
				}
				if response.Email != mockRegister.output.Email {
					t.Errorf("Response Email = %v, want %v", response.Email, mockRegister.output.Email)
				}
			}
		})
	}
}

func TestAuthController_Login(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		setupMock          func(*mockLoginUserUsecase)
		expectedStatusCode int
		expectError        bool
	}{
		{
			name: "Successful login",
			requestBody: LoginRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockLoginUserUsecase) {
				m.output = &usecase.LoginUserOutput{
					Token: "jwt-token-123",
					User: usecase.UserDTO{
						ID:    "user-123",
						Name:  "John Doe",
						Email: "john@example.com",
					},
				}
				m.error = nil
			},
			expectedStatusCode: http.StatusOK,
			expectError:        false,
		},
		{
			name: "Invalid JSON payload",
			requestBody: "invalid-json",
			setupMock: func(m *mockLoginUserUsecase) {
				m.output = nil
				m.error = nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectError:        true,
		},
		{
			name: "Missing email",
			requestBody: LoginRequest{
				Email:    "",
				Password: "password123",
			},
			setupMock: func(m *mockLoginUserUsecase) {
				m.output = nil
				m.error = nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectError:        true,
		},
		{
			name: "Invalid email format",
			requestBody: LoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			setupMock: func(m *mockLoginUserUsecase) {
				m.output = nil
				m.error = nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectError:        true,
		},
		{
			name: "Missing password",
			requestBody: LoginRequest{
				Email:    "john@example.com",
				Password: "",
			},
			setupMock: func(m *mockLoginUserUsecase) {
				m.output = nil
				m.error = nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectError:        true,
		},
		{
			name: "User not found",
			requestBody: LoginRequest{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockLoginUserUsecase) {
				m.output = nil
				m.error = valueobject.NewMyError(valueobject.NotFoundCode, "User not found")
			},
			expectedStatusCode: http.StatusNotFound,
			expectError:        true,
		},
		{
			name: "Invalid credentials",
			requestBody: LoginRequest{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(m *mockLoginUserUsecase) {
				m.output = nil
				m.error = valueobject.NewMyError(valueobject.UnauthorizedCode, "Invalid credentials")
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegister := &mockRegisterUserUsecase{}
			mockLogin := &mockLoginUserUsecase{}
			tt.setupMock(mockLogin)
			
			controller := NewAuthController(mockRegister, mockLogin)
			
			// Create request body
			var bodyBytes []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}
			
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			
			recorder := httptest.NewRecorder()
			controller.Login(recorder, req)
			
			// Check status code
			if recorder.Code != tt.expectedStatusCode {
				t.Errorf("Login() status code = %v, want %v", recorder.Code, tt.expectedStatusCode)
			}
			
			// Check content type
			if recorder.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Login() Content-Type = %v, want application/json", recorder.Header().Get("Content-Type"))
			}
			
			// Check response body for successful login
			if !tt.expectError && recorder.Code == http.StatusOK {
				var response LoginResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal success response: %v", err)
				}
				
				if response.Token != mockLogin.output.Token {
					t.Errorf("Response Token = %v, want %v", response.Token, mockLogin.output.Token)
				}
				if response.User.ID != mockLogin.output.User.ID {
					t.Errorf("Response User.ID = %v, want %v", response.User.ID, mockLogin.output.User.ID)
				}
				if response.User.Name != mockLogin.output.User.Name {
					t.Errorf("Response User.Name = %v, want %v", response.User.Name, mockLogin.output.User.Name)
				}
				if response.User.Email != mockLogin.output.User.Email {
					t.Errorf("Response User.Email = %v, want %v", response.User.Email, mockLogin.output.User.Email)
				}
			}
		})
	}
}