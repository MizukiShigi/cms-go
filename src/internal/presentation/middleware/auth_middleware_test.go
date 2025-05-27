package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	appcontext "github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	secretKey := "test-secret-key"
	userID := "123e4567-e89b-12d3-a456-426614174000"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secretKey))

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		extractedUserID := r.Context().Value(appcontext.UserIDKey)
		assert.Equal(t, userID, extractedUserID)
		w.WriteHeader(http.StatusOK)
	})

	middleware := AuthMiddleware(secretKey)
	handler := middleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	secretKey := "test-secret-key"

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	middleware := AuthMiddleware(secretKey)
	handler := middleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_InvalidAuthorizationHeaderFormat(t *testing.T) {
	secretKey := "test-secret-key"

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	middleware := AuthMiddleware(secretKey)
	handler := middleware(nextHandler)

	testCases := []struct {
		name   string
		header string
	}{
		{"No Bearer prefix", "invalid-token"},
		{"Wrong prefix", "Basic invalid-token"},
		{"Bearer only", "Bearer"},
		{"Empty bearer", "Bearer "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tc.header)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	secretKey := "test-secret-key"

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	middleware := AuthMiddleware(secretKey)
	handler := middleware(nextHandler)

	testCases := []struct {
		name        string
		tokenString string
	}{
		{"Malformed token", "invalid.token.format"},
		{"Empty token", ""},
		{"Random string", "random-string"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", "Bearer "+tc.tokenString)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		})
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	secretKey := "test-secret-key"
	userID := "123e4567-e89b-12d3-a456-426614174000"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secretKey))

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	middleware := AuthMiddleware(secretKey)
	handler := middleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_TokenWithWrongSignature(t *testing.T) {
	secretKey := "test-secret-key"
	wrongSecretKey := "wrong-secret-key"
	userID := "123e4567-e89b-12d3-a456-426614174000"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(wrongSecretKey))

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	middleware := AuthMiddleware(secretKey)
	handler := middleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_TokenWithoutUserID(t *testing.T) {
	secretKey := "test-secret-key"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secretKey))

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	middleware := AuthMiddleware(secretKey)
	handler := middleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_TokenWithInvalidUserIDType(t *testing.T) {
	secretKey := "test-secret-key"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 12345,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secretKey))

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	middleware := AuthMiddleware(secretKey)
	handler := middleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestExtractTokenFromHeader_Success(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	token, err := extractTokenFromHeader(req)

	assert.NoError(t, err)
	assert.Equal(t, "valid-token", token)
}

func TestExtractTokenFromHeader_MissingHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	token, err := extractTokenFromHeader(req)

	assert.Error(t, err)
	assert.Empty(t, token)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.UnauthorizedCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Authorization header is required")
}

func TestExtractTokenFromHeader_InvalidFormat(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "invalid-format")

	token, err := extractTokenFromHeader(req)

	assert.Error(t, err)
	assert.Empty(t, token)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.UnauthorizedCode, myErr.Code)
	assert.Contains(t, myErr.Message, "Invalid authorization header format")
}

func TestValidateToken_Success(t *testing.T) {
	secretKey := "test-secret-key"
	userID := "123e4567-e89b-12d3-a456-426614174000"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secretKey))

	claims, err := validateToken(tokenString, secretKey)

	assert.NoError(t, err)
	assert.Equal(t, userID, claims["user_id"])
}

func TestValidateToken_InvalidToken(t *testing.T) {
	secretKey := "test-secret-key"

	claims, err := validateToken("invalid-token", secretKey)

	assert.Error(t, err)
	assert.Nil(t, claims)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.UnauthorizedCode, myErr.Code)
}