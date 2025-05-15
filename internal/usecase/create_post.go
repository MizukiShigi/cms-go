package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/repository"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type CreatePostInput struct {
	Title   string
	Content string
	Tags    []string
	UserID  string
}

type CreatePostOutput struct {
	ID      string
	Title   string
	Content string
	Tags    []string
	UserID  string
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
		return nil, valueobject.NewMyError(valueobject.InvalidCode, "Invalid title")
	}

	content, err := valueobject.NewPostContent(input.Content)
	if err != nil {
		return nil, valueobject.NewMyError(valueobject.InvalidCode, "Invalid content")
	}

	userID, err := valueobject.ParseUserID(input.UserID)
	if err != nil {
		return nil, valueobject.NewMyError(valueobject.InvalidCode, "Invalid user ID")
	}

	post, err := entity.NewPost(title, content, userID)
	if err != nil {
		return nil, valueobject.NewMyError(valueobject.InvalidCode, "Invalid content")
	}
	for _, tag := range input.Tags {
		tag, err := valueobject.NewTag(tag)
		if err != nil {
			return nil, valueobject.NewMyError(valueobject.InvalidCode, "Invalid tag")
		}
		post.AddTag(tag)
	}

	err = u.postRepository.Create(ctx, post)
	if err != nil {
		errMsg := "Failed to create post"
		slog.ErrorContext(ctx, fmt.Sprintf("%s: %s", errMsg, err))
		return nil, valueobject.NewMyError(valueobject.InternalServerErrorCode, errMsg)
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
		UserID:  post.UserID.String(),
	}, nil
}
