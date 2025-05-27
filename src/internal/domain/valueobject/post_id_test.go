package valueobject

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewPostID(t *testing.T) {
	postID := NewPostID()
	
	// Check if it's a valid UUID
	_, err := uuid.Parse(postID.String())
	if err != nil {
		t.Errorf("NewPostID() generated invalid UUID: %v", err)
	}
	
	// Check if two generated UUIDs are different
	postID2 := NewPostID()
	if postID.Equals(postID2) {
		t.Error("NewPostID() should generate unique IDs")
	}
}

func TestParsePostID(t *testing.T) {
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
			postID, err := ParsePostID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePostID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// For valid UUIDs, check if the string representation is correct
				if postID.String() != tt.input {
					t.Errorf("ParsePostID() = %v, want %v", postID.String(), tt.input)
				}
			}
		})
	}
}

func TestPostID_String(t *testing.T) {
	uuidStr := "550e8400-e29b-41d4-a716-446655440000"
	postID := PostID(uuidStr)
	if got := postID.String(); got != uuidStr {
		t.Errorf("PostID.String() = %v, want %v", got, uuidStr)
	}
}

func TestPostID_Equals(t *testing.T) {
	postID1 := PostID("550e8400-e29b-41d4-a716-446655440000")
	postID2 := PostID("550e8400-e29b-41d4-a716-446655440000")
	postID3 := PostID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	if !postID1.Equals(postID2) {
		t.Error("Expected equal PostIDs to be equal")
	}
	if postID1.Equals(postID3) {
		t.Error("Expected different PostIDs to not be equal")
	}
}