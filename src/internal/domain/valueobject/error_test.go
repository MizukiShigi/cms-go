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
		t.Errorf("NewMyError() Code = %v, want %v", err.Code, code)
	}
	if err.Message != message {
		t.Errorf("NewMyError() Message = %v, want %v", err.Message, message)
	}
}

func TestMyError_Error(t *testing.T) {
	message := "Test error message"
	err := NewMyError(InvalidCode, message)
	
	if got := err.Error(); got != message {
		t.Errorf("MyError.Error() = %v, want %v", got, message)
	}
}

func TestMyError_StatusCode(t *testing.T) {
	tests := []struct {
		name string
		code Code
		want int
	}{
		{
			name: "Invalid code",
			code: InvalidCode,
			want: http.StatusBadRequest,
		},
		{
			name: "Internal server error code",
			code: InternalServerErrorCode,
			want: http.StatusInternalServerError,
		},
		{
			name: "Unauthorized code",
			code: UnauthorizedCode,
			want: http.StatusUnauthorized,
		},
		{
			name: "Forbidden code",
			code: ForbiddenCode,
			want: http.StatusForbidden,
		},
		{
			name: "Not found code",
			code: NotFoundCode,
			want: http.StatusNotFound,
		},
		{
			name: "Conflict code",
			code: ConflictCode,
			want: http.StatusConflict,
		},
		{
			name: "Unknown code",
			code: Code("UNKNOWN"),
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewMyError(tt.code, "test message")
			if got := err.StatusCode(); got != tt.want {
				t.Errorf("MyError.StatusCode() = %v, want %v", got, tt.want)
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
			name: "MyError",
			err:  NewMyError(InvalidCode, "test"),
			want: true,
		},
		{
			name: "Standard error",
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
			if got := IsMyError(tt.err); got != tt.want {
				t.Errorf("IsMyError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReturnMyError(t *testing.T) {
	myErr := NewMyError(InvalidCode, "my error")
	fallbackErr := NewMyError(InternalServerErrorCode, "fallback")

	tests := []struct {
		name string
		err  error
		want *MyError
	}{
		{
			name: "Returns MyError when input is MyError",
			err:  myErr,
			want: myErr,
		},
		{
			name: "Returns fallback when input is standard error",
			err:  errors.New("standard error"),
			want: fallbackErr,
		},
		{
			name: "Returns fallback when input is nil",
			err:  nil,
			want: fallbackErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReturnMyError(tt.err, fallbackErr)
			if got != tt.want {
				t.Errorf("ReturnMyError() = %v, want %v", got, tt.want)
			}
		})
	}
}