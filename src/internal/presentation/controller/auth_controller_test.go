package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/MizukiShigi/cms-go/internal/usecase"
)

// Mock implementations for testing
type mockRegisterUserUsecase struct {
	shouldFail bool
	failError  error
	output     *usecase.RegisterUserOutput
}

func (m *mockRegisterUserUsecase) Execute(ctx context.Context, input *usecase.RegisterUserInput) (*usecase.RegisterUserOutput, error) {
	if m.shouldFail {
		if m.failError != nil {
			return nil, m.failError
		}
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Mock error")
	}
	
	if m.output != nil {
		return m.output, nil
	}
	
	return &usecase.RegisterUserOutput{
		ID:    "test-user-id",
		Name:  input.Name,
		Email: input.Email,
	}, nil
}

type mockLoginUserUsecase struct {
	shouldFail bool
	failError  error
	output     *usecase.LoginUserOutput
}

func (m *mockLoginUserUsecase) Execute(ctx context.Context, input *usecase.LoginUserInput) (*usecase.LoginUserOutput, error) {
	if m.shouldFail {
		if m.failError != nil {
			return nil, m.failError
		}
		return nil, valueobject.NewMyError(valueobject.UnauthorizedCode, "Invalid credentials")
	}
	
	if m.output != nil {
		return m.output, nil
	}
	
	return &usecase.LoginUserOutput{
		Token: "test-jwt-token",
		User: usecase.UserDTO{
			ID:    "test-user-id",
			Name:  "Test User",
			Email: input.Email,
		},
	}, nil
}

func TestNewAuthController(t *testing.T) {
	registerUsecase := &mockRegisterUserUsecase{}
	loginUsecase := &mockLoginUserUsecase{}
	
	controller := NewAuthController(registerUsecase, loginUsecase)
	
	if controller.registerUserUsecase != registerUsecase {
		t.Error("NewAuthController() should set registerUserUsecase")
	}
	if controller.loginUserUsecase != loginUsecase {
		t.Error("NewAuthController() should set loginUserUsecase")
	}
}

func TestAuthController_Register_Success(t *testing.T) {
	registerUsecase := &mockRegisterUserUsecase{}
	loginUsecase := &mockLoginUserUsecase{}
	controller := NewAuthController(registerUsecase, loginUsecase)
	
	reqBody := RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	controller.Register(w, req)
	
	// Check status code
	if w.Code != http.StatusCreated {
		t.Errorf("Register() status = %v, want %v", w.Code, http.StatusCreated)
	}
	
	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Register() content-type = %v, want application/json", contentType)
	}
	
	// Check response body
	var response UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if response.Name != reqBody.Name {
		t.Errorf("Register() response name = %v, want %v", response.Name, reqBody.Name)
	}
	if response.Email != reqBody.Email {
		t.Errorf("Register() response email = %v, want %v", response.Email, reqBody.Email)
	}
	if response.ID == "" {
		t.Error("Register() response should include user ID")
	}
}

func TestAuthController_Register_InvalidJSON(t *testing.T) {
	registerUsecase := &mockRegisterUserUsecase{}
	loginUsecase := &mockLoginUserUsecase{}
	controller := NewAuthController(registerUsecase, loginUsecase)
	
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	controller.Register(w, req)
	
	// Should return bad request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Register() with invalid JSON status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestAuthController_Register_ValidationError(t *testing.T) {
	registerUsecase := &mockRegisterUserUsecase{}
	loginUsecase := &mockLoginUserUsecase{}
	controller := NewAuthController(registerUsecase, loginUsecase)
	
	tests := []struct {
		name string
		req  RegisterRequest
	}{
		{
			name: "Missing name",
			req: RegisterRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
		},
		{
			name: "Invalid email",
			req: RegisterRequest{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "password123",
			},
		},
		{
			name: "Missing password",
			req: RegisterRequest{
				Name:  "John Doe",
				Email: "john@example.com",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.req)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			controller.Register(w, req)
			
			// Should return bad request for validation errors
			if w.Code != http.StatusBadRequest {
				t.Errorf("Register() validation error status = %v, want %v", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestAuthController_Register_UsecaseError(t *testing.T) {
	registerUsecase := &mockRegisterUserUsecase{
		shouldFail: true,
		failError:  valueobject.NewMyError(valueobject.ConflictCode, "User already exists"),
	}
	loginUsecase := &mockLoginUserUsecase{}
	controller := NewAuthController(registerUsecase, loginUsecase)
	
	reqBody := RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	controller.Register(w, req)
	
	// Should return conflict status
	if w.Code != http.StatusConflict {
		t.Errorf("Register() usecase error status = %v, want %v", w.Code, http.StatusConflict)
	}
}

func TestAuthController_Login_Success(t *testing.T) {
	registerUsecase := &mockRegisterUserUsecase{}
	loginUsecase := &mockLoginUserUsecase{}
	controller := NewAuthController(registerUsecase, loginUsecase)
	
	reqBody := LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	controller.Login(w, req)
	
	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Login() status = %v, want %v", w.Code, http.StatusOK)
	}
	
	// Check response body
	var response LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if response.Token == "" {
		t.Error("Login() response should include token")
	}
	if response.User.Email != reqBody.Email {
		t.Errorf("Login() response user email = %v, want %v", response.User.Email, reqBody.Email)
	}
}

func TestAuthController_Login_InvalidJSON(t *testing.T) {
	registerUsecase := &mockRegisterUserUsecase{}
	loginUsecase := &mockLoginUserUsecase{}
	controller := NewAuthController(registerUsecase, loginUsecase)
	
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	controller.Login(w, req)
	
	// Should return bad request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Login() with invalid JSON status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestAuthController_Login_ValidationError(t *testing.T) {
	registerUsecase := &mockRegisterUserUsecase{}
	loginUsecase := &mockLoginUserUsecase{}
	controller := NewAuthController(registerUsecase, loginUsecase)
	
	tests := []struct {
		name string
		req  LoginRequest
	}{
		{
			name: "Invalid email",
			req: LoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
		},
		{
			name: "Missing password",
			req: LoginRequest{
				Email: "john@example.com",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.req)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			controller.Login(w, req)
			
			// Should return bad request for validation errors
			if w.Code != http.StatusBadRequest {
				t.Errorf("Login() validation error status = %v, want %v", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestAuthController_Login_UsecaseError(t *testing.T) {
	registerUsecase := &mockRegisterUserUsecase{}
	loginUsecase := &mockLoginUserUsecase{
		shouldFail: true,
		failError:  valueobject.NewMyError(valueobject.UnauthorizedCode, "Invalid credentials"),
	}
	controller := NewAuthController(registerUsecase, loginUsecase)
	
	reqBody := LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}
	
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	controller.Login(w, req)
	
	// Should return unauthorized status
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Login() usecase error status = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}