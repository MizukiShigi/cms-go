package entity

import (
	"testing"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestNewTagWithName(t *testing.T) {
	tagName, err := valueobject.NewTagName("test-tag")
	if err != nil {
		t.Fatalf("Failed to create tag name: %v", err)
	}
	
	tag := NewTagWithName(tagName)
	
	// Check tag name
	if !tag.Name.Equals(tagName) {
		t.Errorf("NewTagWithName() Name = %v, want %v", tag.Name, tagName)
	}
	
	// Check ID is not set (should be zero value)
	if tag.ID.String() != "" {
		t.Error("NewTagWithName() should not set ID")
	}
	
	// Check timestamps
	if tag.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	if tag.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
	
	// Check timestamps are approximately equal (within 1 second)
	if tag.CreatedAt.Sub(tag.UpdatedAt).Abs() > time.Second {
		t.Error("CreatedAt and UpdatedAt should be approximately equal")
	}
}

func TestParseTag(t *testing.T) {
	tagID := valueobject.NewTagID()
	tagName, _ := valueobject.NewTagName("parsed-tag")
	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()
	
	tag := ParseTag(tagID, tagName, createdAt, updatedAt)
	
	// Verify all fields
	if !tag.ID.String() == tagID.String() {
		t.Errorf("ParseTag() ID = %v, want %v", tag.ID, tagID)
	}
	if !tag.Name.Equals(tagName) {
		t.Errorf("ParseTag() Name = %v, want %v", tag.Name, tagName)
	}
	if !tag.CreatedAt.Equal(createdAt) {
		t.Errorf("ParseTag() CreatedAt = %v, want %v", tag.CreatedAt, createdAt)
	}
	if !tag.UpdatedAt.Equal(updatedAt) {
		t.Errorf("ParseTag() UpdatedAt = %v, want %v", tag.UpdatedAt, updatedAt)
	}
}

func TestTag_Fields(t *testing.T) {
	// Test that Tag struct has all expected fields
	tagID := valueobject.NewTagID()
	tagName, _ := valueobject.NewTagName("field-test")
	now := time.Now()
	
	tag := &Tag{
		ID:        tagID,
		Name:      tagName,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	// Check that all fields are accessible
	if tag.ID.String() != tagID.String() {
		t.Error("ID field not properly set")
	}
	if !tag.Name.Equals(tagName) {
		t.Error("Name field not properly set")
	}
	if !tag.CreatedAt.Equal(now) {
		t.Error("CreatedAt field not properly set")
	}
	if !tag.UpdatedAt.Equal(now) {
		t.Error("UpdatedAt field not properly set")
	}
}

func TestNewTagWithName_MultipleCreations(t *testing.T) {
	tagName1, _ := valueobject.NewTagName("tag-one")
	tagName2, _ := valueobject.NewTagName("tag-two")
	
	tag1 := NewTagWithName(tagName1)
	tag2 := NewTagWithName(tagName2)
	
	// Tags should have different names
	if tag1.Name.Equals(tag2.Name) {
		t.Error("Different tag names should create different tags")
	}
	
	// Both should have timestamps set
	if tag1.CreatedAt.IsZero() || tag2.CreatedAt.IsZero() {
		t.Error("Both tags should have CreatedAt set")
	}
	if tag1.UpdatedAt.IsZero() || tag2.UpdatedAt.IsZero() {
		t.Error("Both tags should have UpdatedAt set")
	}
	
	// Creation times should be close but could be slightly different
	timeDiff := tag2.CreatedAt.Sub(tag1.CreatedAt).Abs()
	if timeDiff > time.Second {
		t.Error("Creation times should be close")
	}
}