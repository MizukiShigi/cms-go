package auth

import (
	"context"
	"testing"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/golang-jwt/jwt/v4"
)

func TestNewJWTService(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)
	
	if service.secretKey != secretKey {
		t.Errorf("NewJWTService() secretKey = %v, want %v", service.secretKey, secretKey)
	}
	
	expectedExpiration := time.Hour * 24
	if service.expiresAt != expectedExpiration {
		t.Errorf("NewJWTService() expiresAt = %v, want %v", service.expiresAt, expectedExpiration)
	}
}

func TestJWTService_GenerateToken(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)
	ctx := context.Background()
	
	userID := valueobject.NewUserID()
	email, _ := valueobject.NewEmail("test@example.com")
	
	// Test successful token generation
	tokenString, err := service.GenerateToken(ctx, userID, email)
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
		return
	}
	
	if tokenString == "" {
		t.Error("GenerateToken() should return non-empty token")
	}
	
	// Parse and verify the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	
	if err != nil {
		t.Errorf("Failed to parse generated token: %v", err)
		return
	}
	
	if !token.Valid {
		t.Error("Generated token should be valid")
	}
	
	claims, ok := token.Claims.(*Claims)
	if !ok {
		t.Error("Claims should be of type *Claims")
		return
	}
	
	// Verify claims
	if claims.UserID != userID.String() {
		t.Errorf("Claims UserID = %v, want %v", claims.UserID, userID.String())
	}
	if claims.Email != email.String() {
		t.Errorf("Claims Email = %v, want %v", claims.Email, email.String())
	}
	
	// Verify expiration time (should be approximately 24 hours from now)
	expectedExpiry := time.Now().Add(time.Hour * 24)
	if claims.ExpiresAt == nil {
		t.Error("ExpiresAt should be set")
	} else {
		expiryDiff := claims.ExpiresAt.Time.Sub(expectedExpiry).Abs()
		if expiryDiff > time.Minute {
			t.Errorf("ExpiresAt difference too large: %v", expiryDiff)
		}
	}
	
	// Verify issued at time (should be approximately now)
	if claims.IssuedAt == nil {
		t.Error("IssuedAt should be set")
	} else {
		issuedDiff := claims.IssuedAt.Time.Sub(time.Now()).Abs()
		if issuedDiff > time.Minute {
			t.Errorf("IssuedAt difference too large: %v", issuedDiff)
		}
	}
}

func TestJWTService_GenerateToken_DifferentUsers(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)
	ctx := context.Background()
	
	userID1 := valueobject.NewUserID()
	email1, _ := valueobject.NewEmail("user1@example.com")
	
	userID2 := valueobject.NewUserID()
	email2, _ := valueobject.NewEmail("user2@example.com")
	
	token1, err1 := service.GenerateToken(ctx, userID1, email1)
	token2, err2 := service.GenerateToken(ctx, userID2, email2)
	
	if err1 != nil || err2 != nil {
		t.Errorf("GenerateToken() errors: %v, %v", err1, err2)
		return
	}
	
	if token1 == token2 {
		t.Error("Different users should generate different tokens")
	}
}

func TestJWTService_GenerateToken_EmptySecret(t *testing.T) {
	// Test with empty secret key - should still work but be less secure
	service := NewJWTService("")
	ctx := context.Background()
	
	userID := valueobject.NewUserID()
	email, _ := valueobject.NewEmail("test@example.com")
	
	token, err := service.GenerateToken(ctx, userID, email)
	if err != nil {
		t.Errorf("GenerateToken() with empty secret error = %v", err)
	}
	if token == "" {
		t.Error("Should generate token even with empty secret")
	}
}

func TestClaims_Structure(t *testing.T) {
	// Test that Claims struct has the expected structure
	userID := "test-user-id"
	email := "test@example.com"
	
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	if claims.UserID != userID {
		t.Errorf("Claims UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("Claims Email = %v, want %v", claims.Email, email)
	}
	if claims.ExpiresAt == nil {
		t.Error("ExpiresAt should be set")
	}
	if claims.IssuedAt == nil {
		t.Error("IssuedAt should be set")
	}
}