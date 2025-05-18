package usecase

import (
	"context"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/repository"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type GetPostInput struct {
	ID valueobject.PostID
}

type GetPostOutput struct {
	ID               valueobject.PostID
	Title            valueobject.PostTitle
	Content          valueobject.PostContent
	Tags             []valueobject.TagName
	Status           valueobject.PostStatus
	FirstPublishedAt *time.Time
	ContentUpdatedAt *time.Time
}

type GetPostUsecase struct {
	postRepository repository.PostRepository
}

func NewGetPostUsecase(postRepository repository.PostRepository) *GetPostUsecase {
	return &GetPostUsecase{postRepository: postRepository}
}

func (u *GetPostUsecase) Execute(ctx context.Context, input *GetPostInput) (*GetPostOutput, error) {
	post, err := u.postRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetPostOutput{
		ID:               post.ID,
		Title:            post.Title,
		Content:          post.Content,
		Status:           post.Status,
		Tags:             post.Tags,
		FirstPublishedAt: post.FirstPublishedAt,
		ContentUpdatedAt: post.ContentUpdatedAt,
	}, nil
}
