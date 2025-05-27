package service

import (
	"context"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

// Mock implementation of AuthService for testing
type mockAuthService struct {
	shouldFail bool
	token      string
	err        error
}

func (m *mockAuthService) GenerateToken(ctx context.Context, userID valueobject.UserID, email valueobject.Email) (string, error) {
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

func TestAuthService_Interface(t *testing.T) {
	// Test that mockAuthService implements AuthService interface
	var service AuthService = &mockAuthService{}
	
	if service == nil {
		t.Error("mockAuthService should implement AuthService interface")
	}
}

func TestMockAuthService_GenerateToken_Success(t *testing.T) {
	service := &mockAuthService{
		token: "test-token-123",
	}
	
	userID := valueobject.NewUserID()
	email, _ := valueobject.NewEmail("test@example.com")
	ctx := context.Background()
	
	token, err := service.GenerateToken(ctx, userID, email)
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
	}
	
	if token != "test-token-123" {
		t.Errorf("GenerateToken() = %v, want %v", token, "test-token-123")
	}
}

func TestMockAuthService_GenerateToken_DefaultToken(t *testing.T) {
	service := &mockAuthService{}
	
	userID := valueobject.NewUserID()
	email, _ := valueobject.NewEmail("test@example.com")
	ctx := context.Background()
	
	token, err := service.GenerateToken(ctx, userID, email)
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
	}
	
	if token != "mock-jwt-token" {
		t.Errorf("GenerateToken() = %v, want %v", token, "mock-jwt-token")
	}
}

func TestMockAuthService_GenerateToken_Failure(t *testing.T) {
	service := &mockAuthService{
		shouldFail: true,
		err:        valueobject.NewMyError(valueobject.InternalServerErrorCode, "Custom error"),
	}
	
	userID := valueobject.NewUserID()
	email, _ := valueobject.NewEmail("test@example.com")
	ctx := context.Background()
	
	token, err := service.GenerateToken(ctx, userID, email)
	if err == nil {
		t.Error("GenerateToken() should return error when shouldFail is true")
	}
	
	if token != "" {
		t.Errorf("GenerateToken() should return empty token on error, got %v", token)
	}
}

func TestMockAuthService_GenerateToken_DefaultFailure(t *testing.T) {
	service := &mockAuthService{
		shouldFail: true,
	}
	
	userID := valueobject.NewUserID()
	email, _ := valueobject.NewEmail("test@example.com")
	ctx := context.Background()
	
	token, err := service.GenerateToken(ctx, userID, email)
	if err == nil {
		t.Error("GenerateToken() should return error when shouldFail is true")
	}
	
	if token != "" {
		t.Errorf("GenerateToken() should return empty token on error, got %v", token)
	}
	
	// Check that error is MyError
	myErr, ok := err.(*valueobject.MyError)
	if !ok {
		t.Error("Error should be MyError type")
	} else if myErr.Code != valueobject.InternalServerErrorCode {
		t.Errorf("Error code = %v, want %v", myErr.Code, valueobject.InternalServerErrorCode)
	}
}