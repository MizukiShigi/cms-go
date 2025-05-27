package valueobject

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewMyError(t *testing.T) {
	code := InvalidCode
	message := "Test error message"
	
	err := NewMyError(code, message)
	
	if err.Code != code {
		t.Errorf("NewMyError().Code = %v, want %v", err.Code, code)
	}
	if err.Message != message {
		t.Errorf("NewMyError().Message = %v, want %v", err.Message, message)
	}
}

func TestMyError_Error(t *testing.T) {
	message := "Test error message"
	err := NewMyError(InvalidCode, message)
	
	result := err.Error()
	
	if result != message {
		t.Errorf("MyError.Error() = %v, want %v", result, message)
	}
}

func TestMyError_StatusCode(t *testing.T) {
	tests := []struct {
		name         string
		code         Code
		expectedCode int
	}{
		{
			name:         "Invalid code",
			code:         InvalidCode,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Internal server error code",
			code:         InternalServerErrorCode,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Unauthorized code",
			code:         UnauthorizedCode,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Forbidden code",
			code:         ForbiddenCode,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "Not found code",
			code:         NotFoundCode,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Conflict code",
			code:         ConflictCode,
			expectedCode: http.StatusConflict,
		},
		{
			name:         "Unknown code",
			code:         Code("UNKNOWN"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewMyError(tt.code, "test message")
			result := err.StatusCode()
			
			if result != tt.expectedCode {
				t.Errorf("MyError.StatusCode() = %v, want %v", result, tt.expectedCode)
			}
		})
	}
}

func TestIsMyError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "MyError type",
			err:  NewMyError(InvalidCode, "test"),
			want: true,
		},
		{
			name: "Standard error type",
			err:  errors.New("standard error"),
			want: false,
		},
		{
			name: "Nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMyError(tt.err)
			if result != tt.want {
				t.Errorf("IsMyError() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestReturnMyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		fallback *MyError
		want     *MyError
	}{
		{
			name:     "MyError type - returns original",
			err:      NewMyError(InvalidCode, "original error"),
			fallback: NewMyError(InternalServerErrorCode, "fallback error"),
			want:     NewMyError(InvalidCode, "original error"),
		},
		{
			name:     "Standard error type - returns fallback",
			err:      errors.New("standard error"),
			fallback: NewMyError(InternalServerErrorCode, "fallback error"),
			want:     NewMyError(InternalServerErrorCode, "fallback error"),
		},
		{
			name:     "Nil error - returns fallback",
			err:      nil,
			fallback: NewMyError(InvalidCode, "fallback error"),
			want:     NewMyError(InvalidCode, "fallback error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReturnMyError(tt.err, tt.fallback)
			
			if result.Code != tt.want.Code {
				t.Errorf("ReturnMyError().Code = %v, want %v", result.Code, tt.want.Code)
			}
			if result.Message != tt.want.Message {
				t.Errorf("ReturnMyError().Message = %v, want %v", result.Message, tt.want.Message)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          *MyError
		expectedCode Code
		expectedHTTP int
	}{
		{
			name:         "InvalidRequestError",
			err:          InvalidRequestError,
			expectedCode: InvalidCode,
			expectedHTTP: http.StatusBadRequest,
		},
		{
			name:         "InternalServerError",
			err:          InternalServerError,
			expectedCode: InternalServerErrorCode,
			expectedHTTP: http.StatusInternalServerError,
		},
		{
			name:         "ForbiddenError",
			err:          ForbiddenError,
			expectedCode: ForbiddenCode,
			expectedHTTP: http.StatusForbidden,
		},
		{
			name:         "NotFoundError",
			err:          NotFoundError,
			expectedCode: NotFoundCode,
			expectedHTTP: http.StatusNotFound,
		},
		{
			name:         "ConflictError",
			err:          ConflictError,
			expectedCode: ConflictCode,
			expectedHTTP: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.expectedCode {
				t.Errorf("%s.Code = %v, want %v", tt.name, tt.err.Code, tt.expectedCode)
			}
			if tt.err.StatusCode() != tt.expectedHTTP {
				t.Errorf("%s.StatusCode() = %v, want %v", tt.name, tt.err.StatusCode(), tt.expectedHTTP)
			}
		})
	}
}