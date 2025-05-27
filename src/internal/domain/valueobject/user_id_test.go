package valueobject

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUserID(t *testing.T) {
	userID := NewUserID()
	
	// Check if it's a valid UUID
	_, err := uuid.Parse(userID.String())
	if err != nil {
		t.Errorf("NewUserID() generated invalid UUID: %v", err)
	}
	
	// Check if two generated UUIDs are different
	userID2 := NewUserID()
	if userID.Equals(userID2) {
		t.Error("NewUserID() should generate unique IDs")
	}
}

func TestParseUserID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid UUID",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "Valid UUID with different format",
			input:   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			wantErr: false,
		},
		{
			name:    "Invalid UUID - too short",
			input:   "550e8400-e29b-41d4-a716",
			wantErr: true,
		},
		{
			name:    "Invalid UUID - wrong format",
			input:   "not-a-uuid",
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Invalid UUID - contains invalid characters",
			input:   "550e8400-e29b-41d4-a716-44665544000g",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := ParseUserID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// For valid UUIDs, check if the string representation is correct
				if userID.String() != tt.input {
					t.Errorf("ParseUserID() = %v, want %v", userID.String(), tt.input)
				}
			}
		})
	}
}

func TestUserID_String(t *testing.T) {
	uuidStr := "550e8400-e29b-41d4-a716-446655440000"
	userID := UserID(uuidStr)
	if got := userID.String(); got != uuidStr {
		t.Errorf("UserID.String() = %v, want %v", got, uuidStr)
	}
}

func TestUserID_Equals(t *testing.T) {
	userID1 := UserID("550e8400-e29b-41d4-a716-446655440000")
	userID2 := UserID("550e8400-e29b-41d4-a716-446655440000")
	userID3 := UserID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	if !userID1.Equals(userID2) {
		t.Error("Expected equal UserIDs to be equal")
	}
	if userID1.Equals(userID3) {
		t.Error("Expected different UserIDs to not be equal")
	}
}