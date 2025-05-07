package usecase

import (
	"context"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/myerror"
	"github.com/MizukiShigi/cms-go/internal/domain/repository"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type CreatePostInput struct {
	Title   string
	Content string
	Tags    []string
}

type CreatePostOutput struct {
	ID      string
	Title   string
	Content string
	Tags    []string
}

type CreatePostUsecase struct {
	postRepository repository.PostRepository
}

func NewCreatePostUsecase(postRepository repository.PostRepository) *CreatePostUsecase {
	return &CreatePostUsecase{postRepository: postRepository}
}

func (u *CreatePostUsecase) Execute(ctx context.Context, input *CreatePostInput) (*CreatePostOutput, error) {
	title, err := valueobject.NewPostTitle(input.Title)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InvalidRequestCode, "Invalid title")
	}

	content, err := valueobject.NewPostContent(input.Content)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InvalidRequestCode, "Invalid content")
	}

	post, err := entity.NewPost(title, content)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InvalidRequestCode, "Invalid content")
	}
	for _, tag := range input.Tags {
		tag, err := valueobject.NewTag(tag)
		if err != nil {
			return nil, myerror.NewMyError(myerror.InvalidRequestCode, "Invalid tag")
		}
		post.AddTag(tag)
	}

	err = u.postRepository.Create(ctx, post)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InternalServerErrorCode, "Failed to create post")
	}

	tags := make([]string, 0, len(post.Tags))
	for _, tag := range post.Tags {
		tags = append(tags, tag.String())
	}

	return &CreatePostOutput{
		ID:      post.ID.String(),
		Title:   post.Title.String(),
		Content: post.Content.String(),
		Tags:    tags,
	}, nil
}
