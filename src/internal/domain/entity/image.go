package entity

import (
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type Image struct {
	ID               valueobject.ImageID
	OriginalFilename string
	StoredFilename   string
	GCSURL           string
	PostID           valueobject.PostID
	UserID           valueobject.UserID
	SortOrder        int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewImage(originalFilename string, storedFilename string, gcsURL string, postID valueobject.PostID, userID valueobject.UserID, sortOrder int) *Image {
	now := time.Now()
	return &Image{
		ID:               valueobject.NewImageID(),
		OriginalFilename: originalFilename,
		StoredFilename:   storedFilename,
		GCSURL:           gcsURL,
		PostID:           postID,
		UserID:           userID,
		SortOrder:        sortOrder,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func ParseImage(
	id valueobject.ImageID,
	originalFilename string,
	storedFilename string,
	gcsURL string,
	postID valueobject.PostID,
	userID valueobject.UserID,
	sortOrder int,
	createdAt time.Time,
	updatedAt time.Time,
) *Image {
	return &Image{
		ID:               id,
		OriginalFilename: originalFilename,
		StoredFilename:   storedFilename,
		GCSURL:           gcsURL,
		PostID:           postID,
		UserID:           userID,
		SortOrder:        sortOrder,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}