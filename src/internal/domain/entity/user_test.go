package entity

import (
	"strings"
	"testing"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	validEmail, _ := valueobject.NewEmail("test@example.com")
	
	tests := []struct {
		name     string
		userName string
		email    valueobject.Email
		password string
		wantErr  bool
		errCode  valueobject.Code
	}{
		{
			name:     "Valid user",
			userName: "John Doe",
			email:    validEmail,
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Empty name",
			userName: "",
			email:    validEmail,
			password: "password123",
			wantErr:  true,
			errCode:  valueobject.InvalidCode,
		},
		{
			name:     "Password too short",
			userName: "John Doe",
			email:    validEmail,
			password: "pass",
			wantErr:  true,
			errCode:  valueobject.InvalidCode,
		},
		{
			name:     "Password exactly 8 characters",
			userName: "John Doe",
			email:    validEmail,
			password: "password",
			wantErr:  false,
		},
		{
			name:     "Long password",
			userName: "John Doe",
			email:    validEmail,
			password: "very_long_password_with_many_characters",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.userName, tt.email, tt.password)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error but got none")
					return
				}
				
				if myErr, ok := err.(*valueobject.MyError); ok {
					if myErr.Code != tt.errCode {
						t.Errorf("NewUser() error code = %v, want %v", myErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("NewUser() error is not MyError type")
				}
				return
			}
			
			if err != nil {
				t.Errorf("NewUser() unexpected error = %v", err)
				return
			}
			
			// Verify user fields
			if user.Name != tt.userName {
				t.Errorf("NewUser().Name = %v, want %v", user.Name, tt.userName)
			}
			
			if !user.Email.Equals(tt.email) {
				t.Errorf("NewUser().Email = %v, want %v", user.Email, tt.email)
			}
			
			// Verify password is hashed
			if user.Password == tt.password {
				t.Errorf("NewUser() password should be hashed, but got plaintext")
			}
			
			// Verify password can be authenticated
			if !user.Authenticate(tt.password) {
				t.Errorf("NewUser() created user cannot authenticate with original password")
			}
			
			// Verify ID is generated
			if user.ID.String() == "" {
				t.Errorf("NewUser() ID should be generated")
			}
			
			// Verify timestamps are set and reasonable
			now := time.Now()
			if user.CreatedAt.IsZero() {
				t.Errorf("NewUser() CreatedAt should be set")
			}
			if user.UpdatedAt.IsZero() {
				t.Errorf("NewUser() UpdatedAt should be set")
			}
			if user.CreatedAt.After(now) {
				t.Errorf("NewUser() CreatedAt should not be in the future")
			}
			if user.UpdatedAt.After(now) {
				t.Errorf("NewUser() UpdatedAt should not be in the future")
			}
			if !user.CreatedAt.Equal(user.UpdatedAt) {
				t.Errorf("NewUser() CreatedAt and UpdatedAt should be equal for new user")
			}
		})
	}
}

func TestUser_Authenticate(t *testing.T) {
	validEmail, _ := valueobject.NewEmail("test@example.com")
	password := "password123"
	
	user, err := NewUser("John Doe", validEmail, password)
	if err != nil {
		t.Fatalf("Failed to create user for test: %v", err)
	}
	
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "Correct password",
			password: "password123",
			want:     true,
		},
		{
			name:     "Incorrect password",
			password: "wrongpassword",
			want:     false,
		},
		{
			name:     "Empty password",
			password: "",
			want:     false,
		},
		{
			name:     "Similar but incorrect password",
			password: "password124",
			want:     false,
		},
		{
			name:     "Case sensitive password",
			password: "Password123",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.Authenticate(tt.password)
			if result != tt.want {
				t.Errorf("User.Authenticate() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestUser_PasswordHashing(t *testing.T) {
	validEmail, _ := valueobject.NewEmail("test@example.com")
	password := "password123"
	
	user1, err := NewUser("User1", validEmail, password)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}
	
	user2, err := NewUser("User2", validEmail, password)
	if err != nil {
		t.Fatalf("Failed to create second user: %v", err)
	}
	
	// Verify that the same password creates different hashes (due to salt)
	if user1.Password == user2.Password {
		t.Errorf("Same password should create different hashes due to salt")
	}
	
	// Verify both users can authenticate with the same password
	if !user1.Authenticate(password) {
		t.Errorf("User1 should authenticate with original password")
	}
	if !user2.Authenticate(password) {
		t.Errorf("User2 should authenticate with original password")
	}
	
	// Verify hashed password looks like bcrypt hash
	if !strings.HasPrefix(user1.Password, "$2a$") && !strings.HasPrefix(user1.Password, "$2b$") && !strings.HasPrefix(user1.Password, "$2y$") {
		t.Errorf("Password should be bcrypt hash, got: %v", user1.Password[:10])
	}
	
	// Verify hash can be verified with bcrypt directly
	err = bcrypt.CompareHashAndPassword([]byte(user1.Password), []byte(password))
	if err != nil {
		t.Errorf("Stored password hash should be valid bcrypt hash: %v", err)
	}
}