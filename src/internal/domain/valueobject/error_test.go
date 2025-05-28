package valueobject

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewMyError(t *testing.T) {
	tests := []struct {
		name    string
		code    Code
		message string
	}{
		{
			name:    "正常ケース: InvalidCodeでエラー作成",
			code:    InvalidCode,
			message: "無効なリクエスト",
		},
		{
			name:    "正常ケース: InternalServerErrorCodeでエラー作成",
			code:    InternalServerErrorCode,
			message: "内部サーバーエラー",
		},
		{
			name:    "正常ケース: UnauthorizedCodeでエラー作成",
			code:    UnauthorizedCode,
			message: "認証が必要です",
		},
		{
			name:    "正常ケース: ForbiddenCodeでエラー作成",
			code:    ForbiddenCode,
			message: "アクセスが禁止されています",
		},
		{
			name:    "正常ケース: NotFoundCodeでエラー作成",
			code:    NotFoundCode,
			message: "リソースが見つかりません",
		},
		{
			name:    "正常ケース: ConflictCodeでエラー作成",
			code:    ConflictCode,
			message: "競合が発生しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewMyError(tt.code, tt.message)

			if err.Code != tt.code {
				t.Errorf("Code = %v, want %v", err.Code, tt.code)
			}
			if err.Message != tt.message {
				t.Errorf("Message = %v, want %v", err.Message, tt.message)
			}
		})
	}
}

func TestMyError_Error(t *testing.T) {
	message := "テストエラーメッセージ"
	err := NewMyError(InvalidCode, message)

	if err.Error() != message {
		t.Errorf("Error() = %v, want %v", err.Error(), message)
	}
}

func TestMyError_StatusCode(t *testing.T) {
	tests := []struct {
		name           string
		code           Code
		expectedStatus int
	}{
		{
			name:           "InvalidCode",
			code:           InvalidCode,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InternalServerErrorCode",
			code:           InternalServerErrorCode,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "UnauthorizedCode",
			code:           UnauthorizedCode,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "ForbiddenCode",
			code:           ForbiddenCode,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "NotFoundCode",
			code:           NotFoundCode,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ConflictCode",
			code:           ConflictCode,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "未知のコード",
			code:           Code("UNKNOWN"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewMyError(tt.code, "テストメッセージ")
			status := err.StatusCode()

			if status != tt.expectedStatus {
				t.Errorf("StatusCode() = %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestIsMyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "正常ケース: MyErrorの場合",
			err:      NewMyError(InvalidCode, "テストエラー"),
			expected: true,
		},
		{
			name:     "正常ケース: 標準エラーの場合",
			err:      errors.New("標準エラー"),
			expected: false,
		},
		{
			name:     "正常ケース: nilエラーの場合",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMyError(tt.err)

			if result != tt.expected {
				t.Errorf("IsMyError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReturnMyError(t *testing.T) {
	defaultErr := NewMyError(InternalServerErrorCode, "デフォルトエラー")
	domainErr := NewMyError(InvalidCode, "ドメインエラー")

	tests := []struct {
		name     string
		err      error
		myerr    *MyError
		expected *MyError
	}{
		{
			name:     "正常ケース: MyErrorが渡された場合はそれを返す",
			err:      domainErr,
			myerr:    defaultErr,
			expected: domainErr,
		},
		{
			name:     "正常ケース: 標準エラーが渡された場合はdefaultを返す",
			err:      errors.New("標準エラー"),
			myerr:    defaultErr,
			expected: defaultErr,
		},
		{
			name:     "正常ケース: nilが渡された場合はdefaultを返す",
			err:      nil,
			myerr:    defaultErr,
			expected: defaultErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReturnMyError(tt.err, tt.myerr)

			if result.Code != tt.expected.Code {
				t.Errorf("ReturnMyError().Code = %v, want %v", result.Code, tt.expected.Code)
			}
			if result.Message != tt.expected.Message {
				t.Errorf("ReturnMyError().Message = %v, want %v", result.Message, tt.expected.Message)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name          string
		err           *MyError
		expectedCode  Code
		expectedMsg   string
		expectedHTTP  int
	}{
		{
			name:         "InvalidRequestError",
			err:          InvalidRequestError,
			expectedCode: InvalidCode,
			expectedMsg:  "Invalid request",
			expectedHTTP: http.StatusBadRequest,
		},
		{
			name:         "InternalServerError",
			err:          InternalServerError,
			expectedCode: InternalServerErrorCode,
			expectedMsg:  "Internal server error",
			expectedHTTP: http.StatusInternalServerError,
		},
		{
			name:         "ForbiddenError",
			err:          ForbiddenError,
			expectedCode: ForbiddenCode,
			expectedMsg:  "Forbidden",
			expectedHTTP: http.StatusForbidden,
		},
		{
			name:         "NotFoundError",
			err:          NotFoundError,
			expectedCode: NotFoundCode,
			expectedMsg:  "Not found",
			expectedHTTP: http.StatusNotFound,
		},
		{
			name:         "ConflictError",
			err:          ConflictError,
			expectedCode: ConflictCode,
			expectedMsg:  "Conflict",
			expectedHTTP: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.expectedCode {
				t.Errorf("Code = %v, want %v", tt.err.Code, tt.expectedCode)
			}
			if tt.err.Message != tt.expectedMsg {
				t.Errorf("Message = %v, want %v", tt.err.Message, tt.expectedMsg)
			}
			if tt.err.StatusCode() != tt.expectedHTTP {
				t.Errorf("StatusCode() = %v, want %v", tt.err.StatusCode(), tt.expectedHTTP)
			}
		})
	}
}