package valueobject

import (
	"strings"
	"testing"
)

func TestNewTagName(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		wantErr bool
		want    string
	}{
		{
			name:    "Valid lowercase tag",
			tag:     "technology",
			wantErr: false,
			want:    "technology",
		},
		{
			name:    "Valid tag with numbers",
			tag:     "tech123",
			wantErr: false,
			want:    "tech123",
		},
		{
			name:    "Valid tag with hyphen",
			tag:     "web-development",
			wantErr: false,
			want:    "web-development",
		},
		{
			name:    "Valid tag with underscore",
			tag:     "go_programming",
			wantErr: false,
			want:    "go_programming",
		},
		{
			name:    "Uppercase tag gets normalized",
			tag:     "GOLANG",
			wantErr: false,
			want:    "golang",
		},
		{
			name:    "Mixed case tag gets normalized",
			tag:     "JavaScript",
			wantErr: false,
			want:    "javascript",
		},
		{
			name:    "Tag with spaces gets trimmed and normalized",
			tag:     "  React  ",
			wantErr: false,
			want:    "react",
		},
		{
			name:    "Tag exactly 50 characters",
			tag:     strings.Repeat("a", 50),
			wantErr: false,
			want:    strings.Repeat("a", 50),
		},
		{
			name:    "Tag too long (51 characters)",
			tag:     strings.Repeat("a", 51),
			wantErr: true,
		},
		{
			name:    "Tag with invalid characters - space",
			tag:     "web development",
			wantErr: true,
		},
		{
			name:    "Tag with invalid characters - special chars",
			tag:     "tag@name",
			wantErr: true,
		},
		{
			name:    "Tag with invalid characters - dot",
			tag:     "tag.name",
			wantErr: true,
		},
		{
			name:    "Empty tag",
			tag:     "",
			wantErr: true,
		},
		{
			name:    "Tag with only spaces",
			tag:     "   ",
			wantErr: true,
		},
		{
			name:    "Single character tag",
			tag:     "a",
			wantErr: false,
			want:    "a",
		},
		{
			name:    "Tag starting with number",
			tag:     "123tag",
			wantErr: false,
			want:    "123tag",
		},
		{
			name:    "Tag starting with hyphen",
			tag:     "-tag",
			wantErr: false,
			want:    "-tag",
		},
		{
			name:    "Tag starting with underscore",
			tag:     "_tag",
			wantErr: false,
			want:    "_tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagName, err := NewTagName(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTagName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tagName.String() != tt.want {
				t.Errorf("NewTagName() = %v, want %v", tagName.String(), tt.want)
			}
		})
	}
}

func TestTagName_String(t *testing.T) {
	tag := "test-tag"
	tagName := TagName(tag)
	if got := tagName.String(); got != tag {
		t.Errorf("TagName.String() = %v, want %v", got, tag)
	}
}

func TestTagName_Equals(t *testing.T) {
	tagName1 := TagName("same-tag")
	tagName2 := TagName("same-tag")
	tagName3 := TagName("different-tag")

	if !tagName1.Equals(tagName2) {
		t.Error("Expected equal TagNames to be equal")
	}
	if tagName1.Equals(tagName3) {
		t.Error("Expected different TagNames to not be equal")
	}
}