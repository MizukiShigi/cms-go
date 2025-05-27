package valueobject

import (
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "Valid email with numbers",
			email:   "user123@example123.com",
			wantErr: false,
		},
		{
			name:    "Valid email with plus",
			email:   "user+test@example.com",
			wantErr: false,
		},
		{
			name:    "Invalid email without @",
			email:   "testexample.com",
			wantErr: true,
		},
		{
			name:    "Invalid email without domain",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "Invalid email without local part",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "Invalid email with spaces",
			email:   "test @example.com",
			wantErr: true,
		},
		{
			name:    "Empty email",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && email.String() != tt.email {
				t.Errorf("NewEmail() = %v, want %v", email.String(), tt.email)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	email := Email("test@example.com")
	if got := email.String(); got != "test@example.com" {
		t.Errorf("Email.String() = %v, want %v", got, "test@example.com")
	}
}

func TestEmail_Equals(t *testing.T) {
	email1 := Email("test@example.com")
	email2 := Email("test@example.com")
	email3 := Email("other@example.com")

	if !email1.Equals(email2) {
		t.Error("Expected equal emails to be equal")
	}
	if email1.Equals(email3) {
		t.Error("Expected different emails to not be equal")
	}
}