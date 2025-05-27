package valueobject

import (
	"regexp"
	"testing"

	"github.com/google/uuid"
)

func TestNewUserID(t *testing.T) {
	userID1 := NewUserID()
	userID2 := NewUserID()

	// Check that IDs are valid UUIDs
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	
	if !uuidRegex.MatchString(userID1.String()) {
		t.Errorf("NewUserID() generated invalid UUID format: %v", userID1.String())
	}
	
	if !uuidRegex.MatchString(userID2.String()) {
		t.Errorf("NewUserID() generated invalid UUID format: %v", userID2.String())
	}

	// Check that different calls generate different IDs
	if userID1.Equals(userID2) {
		t.Errorf("NewUserID() generated the same ID twice: %v", userID1.String())
	}

	// Check that the generated ID can be parsed as UUID
	_, err := uuid.Parse(userID1.String())
	if err != nil {
		t.Errorf("NewUserID() generated ID that cannot be parsed as UUID: %v", err)
	}
}

func TestParseUserID(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	
	tests := []struct {
		name    string
		input   string
		want    UserID
		wantErr bool
	}{
		{
			name:    "Valid UUID",
			input:   validUUID,
			want:    UserID(validUUID),
			wantErr: false,
		},
		{
			name:    "Valid UUID with uppercase",
			input:   "123E4567-E89B-12D3-A456-426614174000",
			want:    UserID("123e4567-e89b-12d3-a456-426614174000"), // UUID normalizes to lowercase
			wantErr: false,
		},
		{
			name:    "Invalid UUID - too short",
			input:   "123e4567-e89b-12d3-a456",
			want:    UserID(""),
			wantErr: true,
		},
		{
			name:    "Invalid UUID - invalid characters",
			input:   "123e4567-e89b-12d3-a456-42661417400g",
			want:    UserID(""),
			wantErr: true,
		},
		{
			name:    "Invalid UUID - missing hyphens",
			input:   "123e4567e89b12d3a456426614174000",
			want:    UserID(""),
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   "",
			want:    UserID(""),
			wantErr: true,
		},
		{
			name:    "Random string",
			input:   "not-a-uuid",
			want:    UserID(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUserID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserID_String(t *testing.T) {
	idString := "123e4567-e89b-12d3-a456-426614174000"
	userID := UserID(idString)
	
	result := userID.String()
	
	if result != idString {
		t.Errorf("UserID.String() = %v, want %v", result, idString)
	}
}

func TestUserID_Equals(t *testing.T) {
	id1 := UserID("123e4567-e89b-12d3-a456-426614174000")
	id2 := UserID("123e4567-e89b-12d3-a456-426614174000")
	id3 := UserID("987e6543-e21b-34d5-c678-123456789abc")
	
	tests := []struct {
		name string
		id1  UserID
		id2  UserID
		want bool
	}{
		{
			name: "Equal IDs",
			id1:  id1,
			id2:  id2,
			want: true,
		},
		{
			name: "Different IDs",
			id1:  id1,
			id2:  id3,
			want: false,
		},
		{
			name: "Empty IDs",
			id1:  UserID(""),
			id2:  UserID(""),
			want: true,
		},
		{
			name: "One empty, one not",
			id1:  id1,
			id2:  UserID(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.id1.Equals(tt.id2)
			if result != tt.want {
				t.Errorf("UserID.Equals() = %v, want %v", result, tt.want)
			}
		})
	}
}