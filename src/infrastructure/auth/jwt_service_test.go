package auth

import (
	"context"
	"strings"
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

	expectedDuration := time.Hour * 24
	if service.expiresAt != expectedDuration {
		t.Errorf("NewJWTService() expiresAt = %v, want %v", service.expiresAt, expectedDuration)
	}
}

func TestJWTService_GenerateToken(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := valueobject.NewUserID()
	email, _ := valueobject.NewEmail("test@example.com")

	tests := []struct {
		name    string
		userID  valueobject.UserID
		email   valueobject.Email
		wantErr bool
	}{
		{
			name:    "Valid token generation",
			userID:  userID,
			email:   email,
			wantErr: false,
		},
		{
			name:    "Empty userID",
			userID:  valueobject.UserID(""),
			email:   email,
			wantErr: false, // JWT service doesn't validate empty values
		},
		{
			name:    "Empty email",
			userID:  userID,
			email:   valueobject.Email(""),
			wantErr: false, // JWT service doesn't validate empty values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.GenerateToken(context.Background(), tt.userID, tt.email)

			if (err != nil) != tt.wantErr {
				t.Errorf("JWTService.GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify token is not empty
				if token == "" {
					t.Errorf("JWTService.GenerateToken() returned empty token")
				}

				// Verify token can be parsed and validated
				parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(secretKey), nil
				})
				if err != nil {
					t.Errorf("Generated token could not be parsed: %v", err)
					return
				}

				if !parsedToken.Valid {
					t.Errorf("Generated token is not valid")
				}

				// Verify claims
				if claims, ok := parsedToken.Claims.(*Claims); ok {
					if claims.UserID != tt.userID.String() {
						t.Errorf("Token UserID = %v, want %v", claims.UserID, tt.userID.String())
					}
					if claims.Email != tt.email.String() {
						t.Errorf("Token Email = %v, want %v", claims.Email, tt.email.String())
					}

					// Verify expiration is set correctly (approximately 24 hours from now)
					expectedExpiry := time.Now().Add(time.Hour * 24)
					actualExpiry := claims.ExpiresAt.Time

					// Allow for small time differences (1 minute)
					timeDiff := actualExpiry.Sub(expectedExpiry)
					if timeDiff > time.Minute || timeDiff < -time.Minute {
						t.Errorf("Token expiry time difference too large: %v", timeDiff)
					}

					// Verify issued at is recent
					issuedAt := claims.IssuedAt.Time
					timeSinceIssued := time.Since(issuedAt)
					if timeSinceIssued > time.Minute {
						t.Errorf("Token issued at time is too old: %v", timeSinceIssued)
					}
				} else {
					t.Errorf("Could not parse token claims")
				}
			}
		})
	}
}

func TestJWTService_TokenFormat(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewJWTService(secretKey)

	userID := valueobject.NewUserID()
	email, _ := valueobject.NewEmail("test@example.com")

	token, err := service.GenerateToken(context.Background(), userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// JWT tokens should have 3 parts separated by dots
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Token should have 3 parts, got %d", len(parts))
	}

	// Each part should be non-empty
	for i, part := range parts {
		if part == "" {
			t.Errorf("Token part %d should not be empty", i)
		}
	}
}

func TestJWTService_DifferentSecretKeys(t *testing.T) {
	userID := valueobject.NewUserID()
	email, _ := valueobject.NewEmail("test@example.com")

	service1 := NewJWTService("secret1")
	service2 := NewJWTService("secret2")

	token1, err := service1.GenerateToken(context.Background(), userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token with service1: %v", err)
	}

	token2, err := service2.GenerateToken(context.Background(), userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token with service2: %v", err)
	}

	// Tokens should be different even with same input due to different secret keys
	if token1 == token2 {
		t.Errorf("Tokens with different secret keys should be different")
	}

	// Token1 should not be valid with secret2
	_, err = jwt.ParseWithClaims(token1, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret2"), nil
	})

	if err == nil {
		t.Errorf("Token generated with secret1 should not be valid with secret2")
	}
}

func TestClaims_Structure(t *testing.T) {
	claims := &Claims{
		UserID: "user-123",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Verify Claims implements jwt.Claims interface
	var _ jwt.Claims = claims

	// Verify custom fields
	if claims.UserID != "user-123" {
		t.Errorf("Claims.UserID = %v, want %v", claims.UserID, "user-123")
	}

	if claims.Email != "test@example.com" {
		t.Errorf("Claims.Email = %v, want %v", claims.Email, "test@example.com")
	}
}
