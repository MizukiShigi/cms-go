package valueobject

import (
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantEmail Email
		wantErr   bool
	}{
		{
			name:      "Valid email",
			email:     "test@example.com",
			wantEmail: Email("test@example.com"),
			wantErr:   false,
		},
		{
			name:      "Valid email with subdomain",
			email:     "user@mail.example.com",
			wantEmail: Email("user@mail.example.com"),
			wantErr:   false,
		},
		{
			name:      "Valid email with numbers",
			email:     "user123@example123.com",
			wantEmail: Email("user123@example123.com"),
			wantErr:   false,
		},
		{
			name:      "Valid email with special characters",
			email:     "user.name+tag@example.com",
			wantEmail: Email("user.name+tag@example.com"),
			wantErr:   false,
		},
		{
			name:      "Invalid email - missing @",
			email:     "userexample.com",
			wantEmail: Email(""),
			wantErr:   true,
		},
		{
			name:      "Invalid email - missing domain",
			email:     "user@",
			wantEmail: Email(""),
			wantErr:   true,
		},
		{
			name:      "Invalid email - missing local part",
			email:     "@example.com",
			wantEmail: Email(""),
			wantErr:   true,
		},
		{
			name:      "Invalid email - no TLD",
			email:     "user@example",
			wantEmail: Email(""),
			wantErr:   true,
		},
		{
			name:      "Invalid email - invalid characters",
			email:     "user@exam ple.com",
			wantEmail: Email(""),
			wantErr:   true,
		},
		{
			name:      "Empty email",
			email:     "",
			wantEmail: Email(""),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmail, err := NewEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEmail != tt.wantEmail {
				t.Errorf("NewEmail() = %v, want %v", gotEmail, tt.wantEmail)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	email := Email("test@example.com")
	expected := "test@example.com"
	
	result := email.String()
	
	if result != expected {
		t.Errorf("Email.String() = %v, want %v", result, expected)
	}
}

func TestEmail_Equals(t *testing.T) {
	tests := []struct {
		name  string
		email Email
		other Email
		want  bool
	}{
		{
			name:  "Equal emails",
			email: Email("test@example.com"),
			other: Email("test@example.com"),
			want:  true,
		},
		{
			name:  "Different emails",
			email: Email("test@example.com"),
			other: Email("other@example.com"),
			want:  false,
		},
		{
			name:  "Empty emails",
			email: Email(""),
			other: Email(""),
			want:  true,
		},
		{
			name:  "One empty, one not",
			email: Email("test@example.com"),
			other: Email(""),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.email.Equals(tt.other)
			if result != tt.want {
				t.Errorf("Email.Equals() = %v, want %v", result, tt.want)
			}
		})
	}
}