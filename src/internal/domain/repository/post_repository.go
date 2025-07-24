package repository

import (
	"context"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) error
	Get(ctx context.Context, id valueobject.PostID) (*entity.Post, error)
	Update(ctx context.Context, post *entity.Post) error
	SetTags(ctx context.Context, post *entity.Post, tags []*entity.Tag) error
	List(ctx context.Context, options *ListPostsOptions) ([]*entity.Post, int, error)
}

// ListPostsOptions は投稿一覧取得のオプション
type ListPostsOptions struct {
	Limit  int
	Offset int
	Status *valueobject.PostStatus
	Sort   string
}
