package valueobject

import (
	"testing"
)

func TestNewPostStatus(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		wantStatus PostStatus
		wantErr    bool
	}{
		{
			name:       "Valid draft status",
			status:     "draft",
			wantStatus: StatusDraft,
			wantErr:    false,
		},
		{
			name:       "Valid published status",
			status:     "published",
			wantStatus: StatusPublished,
			wantErr:    false,
		},
		{
			name:       "Valid private status",
			status:     "private",
			wantStatus: StatusPrivate,
			wantErr:    false,
		},
		{
			name:       "Valid deleted status",
			status:     "deleted",
			wantStatus: StatusDeleted,
			wantErr:    false,
		},
		{
			name:       "Invalid status",
			status:     "invalid",
			wantStatus: PostStatus(""),
			wantErr:    true,
		},
		{
			name:       "Empty status",
			status:     "",
			wantStatus: PostStatus(""),
			wantErr:    true,
		},
		{
			name:       "Case sensitive - uppercase",
			status:     "DRAFT",
			wantStatus: PostStatus(""),
			wantErr:    true,
		},
		{
			name:       "Case sensitive - mixed case",
			status:     "Published",
			wantStatus: PostStatus(""),
			wantErr:    true,
		},
		{
			name:       "Invalid status with spaces",
			status:     "draft ",
			wantStatus: PostStatus(""),
			wantErr:    true,
		},
		{
			name:       "Random string",
			status:     "random_status",
			wantStatus: PostStatus(""),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, err := NewPostStatus(tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPostStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("NewPostStatus() = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestPostStatus_String(t *testing.T) {
	tests := []struct {
		name       string
		postStatus PostStatus
		want       string
	}{
		{
			name:       "Draft status to string",
			postStatus: StatusDraft,
			want:       "draft",
		},
		{
			name:       "Published status to string",
			postStatus: StatusPublished,
			want:       "published",
		},
		{
			name:       "Private status to string",
			postStatus: StatusPrivate,
			want:       "private",
		},
		{
			name:       "Deleted status to string",
			postStatus: StatusDeleted,
			want:       "deleted",
		},
		{
			name:       "Empty status to string",
			postStatus: PostStatus(""),
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.postStatus.String()
			if result != tt.want {
				t.Errorf("PostStatus.String() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestPostStatus_Equals(t *testing.T) {
	tests := []struct {
		name   string
		status PostStatus
		other  PostStatus
		want   bool
	}{
		{
			name:   "Equal draft statuses",
			status: StatusDraft,
			other:  StatusDraft,
			want:   true,
		},
		{
			name:   "Equal published statuses",
			status: StatusPublished,
			other:  StatusPublished,
			want:   true,
		},
		{
			name:   "Different statuses - draft vs published",
			status: StatusDraft,
			other:  StatusPublished,
			want:   false,
		},
		{
			name:   "Different statuses - private vs deleted",
			status: StatusPrivate,
			other:  StatusDeleted,
			want:   false,
		},
		{
			name:   "Empty statuses",
			status: PostStatus(""),
			other:  PostStatus(""),
			want:   true,
		},
		{
			name:   "One empty, one not",
			status: StatusDraft,
			other:  PostStatus(""),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.Equals(tt.other)
			if result != tt.want {
				t.Errorf("PostStatus.Equals() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestPostStatus_Constants(t *testing.T) {
	// Test that the constants have the expected values
	tests := []struct {
		name     string
		constant PostStatus
		expected string
	}{
		{
			name:     "StatusDraft constant",
			constant: StatusDraft,
			expected: "draft",
		},
		{
			name:     "StatusPublished constant",
			constant: StatusPublished,
			expected: "published",
		},
		{
			name:     "StatusPrivate constant",
			constant: StatusPrivate,
			expected: "private",
		},
		{
			name:     "StatusDeleted constant",
			constant: StatusDeleted,
			expected: "deleted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant.String() != tt.expected {
				t.Errorf("Constant %s = %v, want %v", tt.name, tt.constant.String(), tt.expected)
			}
		})
	}
}

func TestPostStatus_AllValidValues(t *testing.T) {
	// Test that all valid constant values can be created successfully
	validStatuses := []string{"draft", "published", "private", "deleted"}
	
	for _, status := range validStatuses {
		t.Run("Valid status: "+status, func(t *testing.T) {
			result, err := NewPostStatus(status)
			if err != nil {
				t.Errorf("NewPostStatus(%v) should not return error, got %v", status, err)
			}
			if result.String() != status {
				t.Errorf("NewPostStatus(%v) = %v, want %v", status, result.String(), status)
			}
		})
	}
}