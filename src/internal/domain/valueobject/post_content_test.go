package valueobject

import (
	"strings"
	"testing"
)

func TestNewPostContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "Valid short content",
			content: "This is a valid post content",
			wantErr: false,
		},
		{
			name:    "Valid long content",
			content: strings.Repeat("a", 9999),
			wantErr: false,
		},
		{
			name:    "Content exactly 10000 characters",
			content: strings.Repeat("a", 10000),
			wantErr: false,
		},
		{
			name:    "Content too long (10001 characters)",
			content: strings.Repeat("a", 10001),
			wantErr: true,
		},
		{
			name:    "Empty content",
			content: "",
			wantErr: false,
		},
		{
			name:    "Content with newlines",
			content: "Line 1\nLine 2\nLine 3",
			wantErr: false,
		},
		{
			name:    "Content with special characters",
			content: "Content with special chars: !@#$%^&*()_+-={}[]|\\:;\"'<>?,./",
			wantErr: false,
		},
		{
			name:    "Content with unicode characters",
			content: "Content with unicode: 日本語 français español 中文",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postContent, err := NewPostContent(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPostContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && postContent.String() != tt.content {
				t.Errorf("NewPostContent() = %v, want %v", postContent.String(), tt.content)
			}
		})
	}
}

func TestPostContent_String(t *testing.T) {
	content := "Test content"
	postContent := PostContent(content)
	if got := postContent.String(); got != content {
		t.Errorf("PostContent.String() = %v, want %v", got, content)
	}
}

func TestPostContent_Equals(t *testing.T) {
	content1 := PostContent("Same content")
	content2 := PostContent("Same content")
	content3 := PostContent("Different content")

	if !content1.Equals(content2) {
		t.Error("Expected equal PostContent to be equal")
	}
	if content1.Equals(content3) {
		t.Error("Expected different PostContent to not be equal")
	}
}