package entity

import (
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type Tag struct {
	ID        valueobject.TagID
	Name      valueobject.TagName
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTagWithName(tagName valueobject.TagName) *Tag {
	now := time.Now()
	return &Tag{
		Name:      tagName,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func ParseTag(
	tagID valueobject.TagID,
	tagName valueobject.TagName,
	createdAt time.Time,
	updatedAt time.Time,
) *Tag {
	return &Tag{
		ID:        tagID,
		Name:      tagName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
