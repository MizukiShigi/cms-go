package usecase

import (
	"context"
	"log/slog"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/repository"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type CreatePostInput struct {
	Title   valueobject.PostTitle
	Content valueobject.PostContent
	Tags    []valueobject.TagName
	UserID  valueobject.UserID
}

type CreatePostOutput struct {
	ID      valueobject.PostID
	Title   valueobject.PostTitle
	Content valueobject.PostContent
	Tags    []valueobject.TagName
	UserID  valueobject.UserID
}

type CreatePostUsecase struct {
	transactionManager repository.TransactionManager
	postRepository     repository.PostRepository
	tagRepository      repository.TagRepository
}

func NewCreatePostUsecase(transactionManager repository.TransactionManager, postRepository repository.PostRepository, tagRepository repository.TagRepository) *CreatePostUsecase {
	return &CreatePostUsecase{
		transactionManager: transactionManager,
		postRepository:     postRepository,
		tagRepository:      tagRepository,
	}
}

func (u *CreatePostUsecase) Execute(ctx context.Context, input *CreatePostInput) (*CreatePostOutput, error) {
	post, err := entity.NewPost(input.Title, input.Content, input.UserID)
	if err != nil {
		return nil, valueobject.NewMyError(valueobject.InvalidCode, "Invalid content")
	}

	for _, tag := range input.Tags {
		post.AddTag(tag)
	}

	transactionErr := u.transactionManager.Transaction(ctx, func(ctx context.Context) error {
		err = u.postRepository.Create(ctx, post)
		if err != nil {
			errMsg := "Failed to create post"
			slog.ErrorContext(ctx, "error", err)
			return valueobject.NewMyError(valueobject.InternalServerErrorCode, errMsg)
		}

		if post.Tags != nil {
			tags := make([]*entity.Tag, 0, len(post.Tags))
			for _, tagName := range post.Tags {
				retTag, err := u.tagRepository.FindOrCreateByName(ctx, entity.NewTagWithName(tagName))
				if err != nil {
					errMsg := "Failed to create tag"
					slog.ErrorContext(ctx, "error", err)
					return valueobject.NewMyError(valueobject.InternalServerErrorCode, errMsg)
				}
				tags = append(tags, retTag)
			}

			err = u.postRepository.SetTags(ctx, post, tags)
			if err != nil {
				errMsg := "Failed to set tags"
				slog.ErrorContext(ctx, "error", err)
				return valueobject.NewMyError(valueobject.InternalServerErrorCode, errMsg)
			}
		}

		return nil
	})

	if transactionErr != nil {
		return nil, transactionErr
	}

	return &CreatePostOutput{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Tags:    post.Tags,
		UserID:  post.UserID,
	}, nil
}
