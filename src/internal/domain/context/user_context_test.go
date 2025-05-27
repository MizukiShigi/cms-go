package context

import (
	"context"
	"errors"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestGetUserID_Success(t *testing.T) {
	testUserID := "test-user-123"
	ctx := context.WithValue(context.Background(), UserID, testUserID)
	
	userID, err := GetUserID(ctx)
	if err != nil {
		t.Errorf("GetUserID() error = %v", err)
		return
	}
	
	if userID != testUserID {
		t.Errorf("GetUserID() = %v, want %v", userID, testUserID)
	}
}

func TestGetUserID_NotFound(t *testing.T) {
	ctx := context.Background()
	
	userID, err := GetUserID(ctx)
	if err == nil {
		t.Error("GetUserID() should return error when user ID not found")
	}
	
	if userID != "" {
		t.Errorf("GetUserID() should return empty string on error, got %v", userID)
	}
	
	// Check that error is MyError with UnauthorizedCode
	var myErr *valueobject.MyError
	if !errors.As(err, &myErr) {
		t.Error("Error should be MyError")
	} else if myErr.Code != valueobject.UnauthorizedCode {
		t.Errorf("Error code = %v, want %v", myErr.Code, valueobject.UnauthorizedCode)
	}
}

func TestGetUserID_WrongType(t *testing.T) {
	// Test with wrong type in context
	ctx := context.WithValue(context.Background(), UserID, 123) // int instead of string
	
	userID, err := GetUserID(ctx)
	if err == nil {
		t.Error("GetUserID() should return error when user ID is wrong type")
	}
	
	if userID != "" {
		t.Errorf("GetUserID() should return empty string on error, got %v", userID)
	}
	
	// Check that error is MyError with UnauthorizedCode
	var myErr *valueobject.MyError
	if !errors.As(err, &myErr) {
		t.Error("Error should be MyError")
	} else if myErr.Code != valueobject.UnauthorizedCode {
		t.Errorf("Error code = %v, want %v", myErr.Code, valueobject.UnauthorizedCode)
	}
}

func TestGetUserID_EmptyString(t *testing.T) {
	// Test with empty string
	ctx := context.WithValue(context.Background(), UserID, "")
	
	userID, err := GetUserID(ctx)
	if err != nil {
		t.Errorf("GetUserID() with empty string should not error: %v", err)
	}
	
	if userID != "" {
		t.Errorf("GetUserID() = %v, want empty string", userID)
	}
}

func TestGetUserID_NilContext(t *testing.T) {
	// Test with nil context - this will panic in Go, but let's test with Background
	ctx := context.Background()
	
	userID, err := GetUserID(ctx)
	if err == nil {
		t.Error("GetUserID() should return error with background context")
	}
	
	if userID != "" {
		t.Errorf("GetUserID() should return empty string on error, got %v", userID)
	}
}