package entity

import (
	"testing"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestNewPost(t *testing.T) {
	title, _ := valueobject.NewPostTitle("Test Title")
	content, _ := valueobject.NewPostContent("Test content")
	userID := valueobject.NewUserID()
	
	tests := []struct {
		name   string
		status valueobject.PostStatus
	}{
		{
			name:   "Draft post",
			status: valueobject.StatusDraft,
		},
		{
			name:   "Published post",
			status: valueobject.StatusPublished,
		},
		{
			name:   "Private post",
			status: valueobject.StatusPrivate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := NewPost(title, content, userID, tt.status)
			if err != nil {
				t.Errorf("NewPost() error = %v", err)
				return
			}
			
			// Check basic fields
			if !post.Title.Equals(title) {
				t.Errorf("NewPost() Title = %v, want %v", post.Title, title)
			}
			if !post.Content.Equals(content) {
				t.Errorf("NewPost() Content = %v, want %v", post.Content, content)
			}
			if !post.UserID.Equals(userID) {
				t.Errorf("NewPost() UserID = %v, want %v", post.UserID, userID)
			}
			if !post.Status.Equals(tt.status) {
				t.Errorf("NewPost() Status = %v, want %v", post.Status, tt.status)
			}
			
			// Check ID is generated
			if post.ID.String() == "" {
				t.Error("Post ID should be generated")
			}
			
			// Check timestamps
			if post.CreatedAt.IsZero() {
				t.Error("CreatedAt should be set")
			}
			if post.UpdatedAt.IsZero() {
				t.Error("UpdatedAt should be set")
			}
			if post.ContentUpdatedAt == nil || post.ContentUpdatedAt.IsZero() {
				t.Error("ContentUpdatedAt should be set")
			}
			
			// Check FirstPublishedAt for published posts
			if tt.status == valueobject.StatusPublished {
				if post.FirstPublishedAt == nil || post.FirstPublishedAt.IsZero() {
					t.Error("FirstPublishedAt should be set for published posts")
				}
			}
			
			// Check Tags is initialized
			if post.Tags == nil {
				t.Error("Tags should be initialized")
			}
			if len(post.Tags) != 0 {
				t.Error("New post should have no tags")
			}
		})
	}
}

func TestPost_AddTag(t *testing.T) {
	title, _ := valueobject.NewPostTitle("Test Title")
	content, _ := valueobject.NewPostContent("Test content")
	userID := valueobject.NewUserID()
	post, _ := NewPost(title, content, userID, valueobject.StatusDraft)
	
	tag1, _ := valueobject.NewTagName("tag1")
	tag2, _ := valueobject.NewTagName("tag2")
	
	// Test adding first tag
	err := post.AddTag(tag1)
	if err != nil {
		t.Errorf("AddTag() error = %v", err)
	}
	if len(post.Tags) != 1 || !post.Tags[0].Equals(tag1) {
		t.Error("Tag should be added")
	}
	
	// Test adding second tag
	err = post.AddTag(tag2)
	if err != nil {
		t.Errorf("AddTag() error = %v", err)
	}
	if len(post.Tags) != 2 {
		t.Error("Second tag should be added")
	}
	
	// Test adding duplicate tag
	err = post.AddTag(tag1)
	if err == nil {
		t.Error("Adding duplicate tag should return error")
	}
	if len(post.Tags) != 2 {
		t.Error("Duplicate tag should not be added")
	}
	
	// Test maximum tags limit
	for i := 3; i <= 10; i++ {
		tag, _ := valueobject.NewTagName("tag" + string(rune('0'+i)))
		err = post.AddTag(tag)
		if err != nil {
			t.Errorf("AddTag() error for tag %d = %v", i, err)
		}
	}
	
	// Try to add 11th tag
	tag11, _ := valueobject.NewTagName("tag11")
	err = post.AddTag(tag11)
	if err == nil {
		t.Error("Adding 11th tag should return error")
	}
	if len(post.Tags) != 10 {
		t.Error("Post should have maximum 10 tags")
	}
}

func TestPost_SetStatus(t *testing.T) {
	title, _ := valueobject.NewPostTitle("Test Title")
	content, _ := valueobject.NewPostContent("Test content")
	userID := valueobject.NewUserID()
	
	tests := []struct {
		name        string
		initStatus  valueobject.PostStatus
		newStatus   valueobject.PostStatus
		wantErr     bool
	}{
		// Draft transitions
		{
			name:       "Draft to Published",
			initStatus: valueobject.StatusDraft,
			newStatus:  valueobject.StatusPublished,
			wantErr:    false,
		},
		{
			name:       "Draft to Deleted",
			initStatus: valueobject.StatusDraft,
			newStatus:  valueobject.StatusDeleted,
			wantErr:    false,
		},
		{
			name:       "Draft to Private (should fail)",
			initStatus: valueobject.StatusDraft,
			newStatus:  valueobject.StatusPrivate,
			wantErr:    true,
		},
		// Published transitions
		{
			name:       "Published to Private",
			initStatus: valueobject.StatusPublished,
			newStatus:  valueobject.StatusPrivate,
			wantErr:    false,
		},
		{
			name:       "Published to Draft (should fail)",
			initStatus: valueobject.StatusPublished,
			newStatus:  valueobject.StatusDraft,
			wantErr:    true,
		},
		{
			name:       "Published to Deleted (should fail)",
			initStatus: valueobject.StatusPublished,
			newStatus:  valueobject.StatusDeleted,
			wantErr:    true,
		},
		// Private transitions
		{
			name:       "Private to Published",
			initStatus: valueobject.StatusPrivate,
			newStatus:  valueobject.StatusPublished,
			wantErr:    false,
		},
		{
			name:       "Private to Deleted",
			initStatus: valueobject.StatusPrivate,
			newStatus:  valueobject.StatusDeleted,
			wantErr:    false,
		},
		{
			name:       "Private to Draft (should fail)",
			initStatus: valueobject.StatusPrivate,
			newStatus:  valueobject.StatusDraft,
			wantErr:    true,
		},
		// Same status
		{
			name:       "Draft to Draft (no change)",
			initStatus: valueobject.StatusDraft,
			newStatus:  valueobject.StatusDraft,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, _ := NewPost(title, content, userID, tt.initStatus)
			
			err := post.SetStatus(tt.newStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if !post.Status.Equals(tt.newStatus) {
					t.Errorf("SetStatus() Status = %v, want %v", post.Status, tt.newStatus)
				}
				
				// Check FirstPublishedAt is set when publishing
				if tt.newStatus == valueobject.StatusPublished && tt.initStatus != valueobject.StatusPublished {
					if post.FirstPublishedAt == nil {
						t.Error("FirstPublishedAt should be set when publishing")
					}
				}
			}
		})
	}
}

func TestParsePost(t *testing.T) {
	postID := valueobject.NewPostID()
	title, _ := valueobject.NewPostTitle("Test Title")
	content, _ := valueobject.NewPostContent("Test content")
	userID := valueobject.NewUserID()
	status := valueobject.StatusPublished
	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()
	firstPublishedAt := time.Now().Add(-30 * time.Minute)
	contentUpdatedAt := time.Now().Add(-15 * time.Minute)
	tag1, _ := valueobject.NewTagName("tag1")
	tag2, _ := valueobject.NewTagName("tag2")
	tags := []valueobject.TagName{tag1, tag2}
	
	post := ParsePost(
		postID,
		title,
		content,
		userID,
		status,
		createdAt,
		updatedAt,
		&firstPublishedAt,
		&contentUpdatedAt,
		tags,
	)
	
	// Verify all fields
	if !post.ID.Equals(postID) {
		t.Errorf("ParsePost() ID = %v, want %v", post.ID, postID)
	}
	if !post.Title.Equals(title) {
		t.Errorf("ParsePost() Title = %v, want %v", post.Title, title)
	}
	if !post.Content.Equals(content) {
		t.Errorf("ParsePost() Content = %v, want %v", post.Content, content)
	}
	if !post.UserID.Equals(userID) {
		t.Errorf("ParsePost() UserID = %v, want %v", post.UserID, userID)
	}
	if !post.Status.Equals(status) {
		t.Errorf("ParsePost() Status = %v, want %v", post.Status, status)
	}
	if !post.CreatedAt.Equal(createdAt) {
		t.Errorf("ParsePost() CreatedAt = %v, want %v", post.CreatedAt, createdAt)
	}
	if !post.UpdatedAt.Equal(updatedAt) {
		t.Errorf("ParsePost() UpdatedAt = %v, want %v", post.UpdatedAt, updatedAt)
	}
	if post.FirstPublishedAt == nil || !post.FirstPublishedAt.Equal(firstPublishedAt) {
		t.Errorf("ParsePost() FirstPublishedAt = %v, want %v", post.FirstPublishedAt, &firstPublishedAt)
	}
	if post.ContentUpdatedAt == nil || !post.ContentUpdatedAt.Equal(contentUpdatedAt) {
		t.Errorf("ParsePost() ContentUpdatedAt = %v, want %v", post.ContentUpdatedAt, &contentUpdatedAt)
	}
	if len(post.Tags) != 2 {
		t.Errorf("ParsePost() Tags length = %v, want %v", len(post.Tags), 2)
	}
}