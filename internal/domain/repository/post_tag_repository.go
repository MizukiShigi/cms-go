package repository

import (
	"context"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type PostTagRepository interface {
	Create(ctx context.Context, postID valueobject.PostID, tagID valueobject.TagID) error
}
