package entity

import (
	"testing"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	email, _ := valueobject.NewEmail("test@example.com")
	
	tests := []struct {
		name     string
		userName string
		email    valueobject.Email
		password string
		wantErr  bool
	}{
		{
			name:     "Valid user",
			userName: "John Doe",
			email:    email,
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Empty name",
			userName: "",
			email:    email,
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "Short password",
			userName: "John Doe",
			email:    email,
			password: "short",
			wantErr:  true,
		},
		{
			name:     "Password exactly 8 characters",
			userName: "John Doe",
			email:    email,
			password: "12345678",
			wantErr:  false,
		},
		{
			name:     "Password 7 characters",
			userName: "John Doe",
			email:    email,
			password: "1234567",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.userName, tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				// Check user fields
				if user.Name != tt.userName {
					t.Errorf("NewUser() Name = %v, want %v", user.Name, tt.userName)
				}
				if !user.Email.Equals(tt.email) {
					t.Errorf("NewUser() Email = %v, want %v", user.Email, tt.email)
				}
				
				// Check password is hashed
				if user.Password == tt.password {
					t.Error("Password should be hashed")
				}
				
				// Check password can be verified
				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.password))
				if err != nil {
					t.Error("Password hash verification failed")
				}
				
				// Check ID is generated
				if user.ID.String() == "" {
					t.Error("User ID should be generated")
				}
				
				// Check timestamps
				if user.CreatedAt.IsZero() {
					t.Error("CreatedAt should be set")
				}
				if user.UpdatedAt.IsZero() {
					t.Error("UpdatedAt should be set")
				}
				
				// Check timestamps are approximately equal (within 1 second)
				if user.CreatedAt.Sub(user.UpdatedAt).Abs() > time.Second {
					t.Error("CreatedAt and UpdatedAt should be approximately equal")
				}
			}
		})
	}
}

func TestUser_Authenticate(t *testing.T) {
	email, _ := valueobject.NewEmail("test@example.com")
	password := "testpassword123"
	user, err := NewUser("Test User", email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "Correct password",
			password: "testpassword123",
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
			name:     "Similar but wrong password",
			password: "testpassword124",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := user.Authenticate(tt.password); got != tt.want {
				t.Errorf("User.Authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_AuthenticateWithHashedPassword(t *testing.T) {
	// Test with manually created user to ensure authentication works with pre-hashed passwords
	email, _ := valueobject.NewEmail("test@example.com")
	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	user := &User{
		ID:        valueobject.NewUserID(),
		Name:      "Test User",
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if !user.Authenticate(password) {
		t.Error("Authentication should succeed with correct password")
	}
	
	if user.Authenticate("wrongpassword") {
		t.Error("Authentication should fail with incorrect password")
	}
}