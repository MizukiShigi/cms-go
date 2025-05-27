package valueobject

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTagID(t *testing.T) {
	tagID := NewTagID()
	
	// Check if it's a valid UUID
	_, err := uuid.Parse(tagID.String())
	if err != nil {
		t.Errorf("NewTagID() generated invalid UUID: %v", err)
	}
	
	// Check if two generated UUIDs are different
	tagID2 := NewTagID()
	if tagID.String() == tagID2.String() {
		t.Error("NewTagID() should generate unique IDs")
	}
}

func TestParseTagID(t *testing.T) {
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
			tagID, err := ParseTagID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTagID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// For valid UUIDs, check if the string representation is correct
				if tagID.String() != tt.input {
					t.Errorf("ParseTagID() = %v, want %v", tagID.String(), tt.input)
				}
			}
		})
	}
}

func TestTagID_String(t *testing.T) {
	uuidStr := "550e8400-e29b-41d4-a716-446655440000"
	tagID := TagID(uuidStr)
	if got := tagID.String(); got != uuidStr {
		t.Errorf("TagID.String() = %v, want %v", got, uuidStr)
	}
}