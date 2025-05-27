package valueobject

import "testing"

func TestNewPostStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
		want    PostStatus
	}{
		{
			name:    "Valid status draft",
			status:  "draft",
			wantErr: false,
			want:    StatusDraft,
		},
		{
			name:    "Valid status published",
			status:  "published",
			wantErr: false,
			want:    StatusPublished,
		},
		{
			name:    "Valid status private",
			status:  "private",
			wantErr: false,
			want:    StatusPrivate,
		},
		{
			name:    "Valid status deleted",
			status:  "deleted",
			wantErr: false,
			want:    StatusDeleted,
		},
		{
			name:    "Invalid status",
			status:  "invalid",
			wantErr: true,
		},
		{
			name:    "Empty status",
			status:  "",
			wantErr: true,
		},
		{
			name:    "Case sensitive - DRAFT should fail",
			status:  "DRAFT",
			wantErr: true,
		},
		{
			name:    "Case sensitive - Published should fail",
			status:  "Published",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPostStatus(tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPostStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("NewPostStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status PostStatus
		want   string
	}{
		{
			name:   "Draft status",
			status: StatusDraft,
			want:   "draft",
		},
		{
			name:   "Published status",
			status: StatusPublished,
			want:   "published",
		},
		{
			name:   "Private status",
			status: StatusPrivate,
			want:   "private",
		},
		{
			name:   "Deleted status",
			status: StatusDeleted,
			want:   "deleted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("PostStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostStatus_Equals(t *testing.T) {
	status1 := StatusDraft
	status2 := StatusDraft
	status3 := StatusPublished

	if !status1.Equals(status2) {
		t.Error("Expected equal PostStatus to be equal")
	}
	if status1.Equals(status3) {
		t.Error("Expected different PostStatus to not be equal")
	}
}

func TestPostStatusConstants(t *testing.T) {
	// Test that constants have the expected values
	tests := []struct {
		name     string
		status   PostStatus
		expected string
	}{
		{"StatusDraft", StatusDraft, "draft"},
		{"StatusPublished", StatusPublished, "published"},
		{"StatusPrivate", StatusPrivate, "private"},
		{"StatusDeleted", StatusDeleted, "deleted"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("Constant %s = %v, want %v", tt.name, string(tt.status), tt.expected)
			}
		})
	}
}