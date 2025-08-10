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
	Status  valueobject.PostStatus
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
	post, err := entity.NewPost(input.Title, input.Content, input.UserID, input.Status)
	if err != nil {
		return nil, valueobject.NewMyError(valueobject.InvalidCode, "Invalid content")
	}

	for _, tag := range input.Tags {
		err := post.AddTag(tag)
		if err != nil {
			return nil, err
		}
	}

	transactionErr := u.transactionManager.Transaction(ctx, func(ctx context.Context) error {
		err = u.postRepository.Create(ctx, post)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to create post")
		}

		tags := make([]*entity.Tag, 0, len(post.Tags))
		for _, tagName := range post.Tags {
			retTag, err := u.tagRepository.FindOrCreateByName(ctx, entity.NewTagWithName(tagName))
			if err != nil {
				return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to create tag")
			}
			tags = append(tags, retTag)
		}

		err = u.postRepository.SetTags(ctx, post, tags)
		if err != nil {
			return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to set tags")
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
