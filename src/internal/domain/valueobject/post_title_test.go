package valueobject

import (
	"strings"
	"testing"
)

func TestNewPostTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			name:    "Valid title",
			title:   "This is a valid title",
			wantErr: false,
		},
		{
			name:    "Valid title with spaces",
			title:   "  Valid title with spaces  ",
			wantErr: false,
		},
		{
			name:    "Title with special characters that get escaped",
			title:   "Title with & and other chars",
			wantErr: false,
		},
		{
			name:    "Title exactly 200 characters",
			title:   strings.Repeat("a", 200),
			wantErr: false,
		},
		{
			name:    "Title too long (201 characters)",
			title:   strings.Repeat("a", 201),
			wantErr: true,
		},
		{
			name:    "Title with forbidden character <",
			title:   "Title with < character",
			wantErr: true,
		},
		{
			name:    "Title with forbidden character >",
			title:   "Title with > character",
			wantErr: true,
		},
		{
			name:    "Title with forbidden character quote",
			title:   "Title with \" character",
			wantErr: true,
		},
		{
			name:    "Title with forbidden character single quote",
			title:   "Title with ' character",
			wantErr: true,
		},
		{
			name:    "Title with forbidden character backslash",
			title:   "Title with \\ character",
			wantErr: true,
		},
		{
			name:    "Title with forbidden character newline",
			title:   "Title with \n character",
			wantErr: true,
		},
		{
			name:    "Title with forbidden character carriage return",
			title:   "Title with \r character",
			wantErr: true,
		},
		{
			name:    "Title with forbidden character tab",
			title:   "Title with \t character",
			wantErr: true,
		},
		{
			name:    "Empty title",
			title:   "",
			wantErr: false,
		},
		{
			name:    "Title with only spaces",
			title:   "   ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postTitle, err := NewPostTitle(tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPostTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Check that the title was trimmed
				expected := strings.TrimSpace(tt.title)
				if postTitle.Value() != expected {
					t.Errorf("NewPostTitle() title = %v, want %v", postTitle.Value(), expected)
				}
			}
		})
	}
}

func TestPostTitle_Value(t *testing.T) {
	title := "Test & Title"
	postTitle, _ := NewPostTitle(title)
	
	// Value() should return unescaped string
	if got := postTitle.Value(); got != title {
		t.Errorf("PostTitle.Value() = %v, want %v", got, title)
	}
}

func TestPostTitle_String(t *testing.T) {
	title := "Test & Title"
	postTitle, _ := NewPostTitle(title)
	
	// String() should return escaped string
	expected := "Test &amp; Title"
	if got := postTitle.String(); got != expected {
		t.Errorf("PostTitle.String() = %v, want %v", got, expected)
	}
}

func TestPostTitle_Equals(t *testing.T) {
	title1, _ := NewPostTitle("Same Title")
	title2, _ := NewPostTitle("Same Title")
	title3, _ := NewPostTitle("Different Title")

	if !title1.Equals(title2) {
		t.Error("Expected equal PostTitles to be equal")
	}
	if title1.Equals(title3) {
		t.Error("Expected different PostTitles to not be equal")
	}
}

func TestPostTitle_LongSpecialCharacters(t *testing.T) {
	// Test title that becomes too long after escaping
	title := strings.Repeat("&", 100) // This will become &amp; repeated 100 times
	_, err := NewPostTitle(title)
	if err == nil {
		t.Error("Expected error for title that becomes too long after HTML escaping")
	}
}