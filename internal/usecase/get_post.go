package usecase

import (
	"context"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/repository"
)

type GetPostInput struct {
	ID string
}

type GetPostOutput struct {
	ID               string
	Title            string
	Content          string
	Tags             []string
	Status           string
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

	tags := make([]string, 0, len(post.Tags))
	for _, tag := range post.Tags {
		tags = append(tags, tag.String())
	}

	return &GetPostOutput{
		ID:      post.ID.String(),
		Title:            post.Title.String(),
		Content:          post.Content.String(),
		Status:           post.Status.String(),
		Tags:             tags,
		FirstPublishedAt: post.FirstPublishedAt,
		ContentUpdatedAt: post.ContentUpdatedAt,
	}, nil
}
