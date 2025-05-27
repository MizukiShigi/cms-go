package valueobject

import (
	"strings"
	"testing"
)

func TestNewPostTitle(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantTitle PostTitle
		wantErr   bool
	}{
		{
			name:      "Valid title",
			title:     "Valid Post Title",
			wantTitle: PostTitle("Valid Post Title"),
			wantErr:   false,
		},
		{
			name:      "Title with leading/trailing spaces",
			title:     "  Trimmed Title  ",
			wantTitle: PostTitle("Trimmed Title"),
			wantErr:   false,
		},
		{
			name:      "Title with special characters that need escaping",
			title:     "Title & More",
			wantTitle: PostTitle("Title &amp; More"),
			wantErr:   false,
		},
		{
			name:      "Empty title after trim",
			title:     "   ",
			wantTitle: PostTitle(""),
			wantErr:   false,
		},
		{
			name:      "Title at max length (200 chars)",
			title:     strings.Repeat("a", 200),
			wantTitle: PostTitle(strings.Repeat("a", 200)),
			wantErr:   false,
		},
		{
			name:      "Title too long (201 chars)",
			title:     strings.Repeat("a", 201),
			wantTitle: PostTitle(""),
			wantErr:   true,
		},
		{
			name:      "Title with forbidden character - <",
			title:     "Title <script>",
			wantTitle: PostTitle(""),
			wantErr:   true,
		},
		{
			name:      "Title with forbidden character - >",
			title:     "Title >text",
			wantTitle: PostTitle(""),
			wantErr:   true,
		},
		{
			name:      "Title with forbidden character - double quote",
			title:     "Title \"quoted\"",
			wantTitle: PostTitle(""),
			wantErr:   true,
		},
		{
			name:      "Title with forbidden character - single quote",
			title:     "Title 'quoted'",
			wantTitle: PostTitle(""),
			wantErr:   true,
		},
		{
			name:      "Title with forbidden character - backslash",
			title:     "Title\\path",
			wantTitle: PostTitle(""),
			wantErr:   true,
		},
		{
			name:      "Title with forbidden character - forward slash",
			title:     "Title/path",
			wantTitle: PostTitle(""),
			wantErr:   true,
		},
		{
			name:      "Title with forbidden character - newline",
			title:     "Title\nNew Line",
			wantTitle: PostTitle(""),
			wantErr:   true,
		},
		{
			name:      "Title with forbidden character - carriage return",
			title:     "Title\rReturn",
			wantTitle: PostTitle(""),
			wantErr:   true,
		},
		{
			name:      "Title with forbidden character - tab",
			title:     "Title\tTab",
			wantTitle: PostTitle(""),
			wantErr:   true,
		},
		{
			name:      "Title with many ampersands that expand when escaped",
			title:     strings.Repeat("&", 60), // Will expand to &amp; (5 chars each)
			wantTitle: PostTitle(""),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitle, err := NewPostTitle(tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPostTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTitle != tt.wantTitle {
				t.Errorf("NewPostTitle() = %v, want %v", gotTitle, tt.wantTitle)
			}
		})
	}
}

func TestPostTitle_Value(t *testing.T) {
	tests := []struct {
		name      string
		postTitle PostTitle
		want      string
	}{
		{
			name:      "Title with escaped ampersand",
			postTitle: PostTitle("Title &amp; More"),
			want:      "Title & More",
		},
		{
			name:      "Title with escaped less than",
			postTitle: PostTitle("Title &lt; More"),
			want:      "Title < More",
		},
		{
			name:      "Title with escaped greater than",
			postTitle: PostTitle("Title &gt; More"),
			want:      "Title > More",
		},
		{
			name:      "Plain title",
			postTitle: PostTitle("Plain Title"),
			want:      "Plain Title",
		},
		{
			name:      "Empty title",
			postTitle: PostTitle(""),
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.postTitle.Value()
			if result != tt.want {
				t.Errorf("PostTitle.Value() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestPostTitle_String(t *testing.T) {
	title := PostTitle("Test Title")
	expected := "Test Title"
	
	result := title.String()
	
	if result != expected {
		t.Errorf("PostTitle.String() = %v, want %v", result, expected)
	}
}

func TestPostTitle_Equals(t *testing.T) {
	tests := []struct {
		name  string
		title PostTitle
		other PostTitle
		want  bool
	}{
		{
			name:  "Equal titles",
			title: PostTitle("Same Title"),
			other: PostTitle("Same Title"),
			want:  true,
		},
		{
			name:  "Different titles",
			title: PostTitle("Title One"),
			other: PostTitle("Title Two"),
			want:  false,
		},
		{
			name:  "Empty titles",
			title: PostTitle(""),
			other: PostTitle(""),
			want:  true,
		},
		{
			name:  "One empty, one not",
			title: PostTitle("Title"),
			other: PostTitle(""),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.title.Equals(tt.other)
			if result != tt.want {
				t.Errorf("PostTitle.Equals() = %v, want %v", result, tt.want)
			}
		})
	}
}