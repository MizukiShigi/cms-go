package usecase

import (
	"context"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/repository"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type UpdatePostUsecase struct {
	transactionManager repository.TransactionManager
	postRepository     repository.PostRepository
	tagRepository      repository.TagRepository
}

type UpdatePostInput struct {
	ID      valueobject.PostID
	Title   valueobject.PostTitle
	Content valueobject.PostContent
	Tags    []valueobject.TagName
	Status  valueobject.PostStatus
}

type UpdatePostOutput struct {
	ID               valueobject.PostID
	Title            valueobject.PostTitle
	Content          valueobject.PostContent
	Tags             []valueobject.TagName
	Status           valueobject.PostStatus
	FirstPublishedAt *time.Time
	ContentUpdatedAt *time.Time
}

func NewUpdatePostUsecase(transactionManager repository.TransactionManager, postRepository repository.PostRepository, tagRepository repository.TagRepository) *UpdatePostUsecase {
	return &UpdatePostUsecase{transactionManager: transactionManager, postRepository: postRepository, tagRepository: tagRepository}
}

func (u *UpdatePostUsecase) Execute(ctx context.Context, input *UpdatePostInput) (*UpdatePostOutput, error) {
	post, err := u.postRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, valueobject.ReturnMyError(err, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to get post"))
	}

	err = u.transactionManager.Transaction(ctx, func(ctx context.Context) error {
		now := time.Now()
		post.Title = input.Title
		post.Content = input.Content
		post.ContentUpdatedAt = &now

		err = u.postRepository.Update(ctx, post)
		if err != nil {
			return valueobject.ReturnMyError(err, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to update post"))
		}

		tags := make([]*entity.Tag, 0, len(input.Tags))
		for _, tagName := range input.Tags {
			inputTag := entity.NewTagWithName(tagName)
			tag, err := u.tagRepository.FindOrCreateByName(ctx, inputTag)
			if err != nil {
				return valueobject.ReturnMyError(err, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to find or create tag"))
			}
			tags = append(tags, tag)
		}

		err = u.postRepository.SetTags(ctx, post, tags)
		if err != nil {
			return valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to set tags")
		}

		return nil
	})
	if err != nil {
		return nil, valueobject.ReturnMyError(err, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to update post"))
	}

	updatePost, err := u.postRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, valueobject.ReturnMyError(err, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to get post"))
	}
	return &UpdatePostOutput{
		ID:               updatePost.ID,
		Title:            updatePost.Title,
		Content:          updatePost.Content,
		Tags:             updatePost.Tags,
		Status:           updatePost.Status,
		FirstPublishedAt: updatePost.FirstPublishedAt,
		ContentUpdatedAt: updatePost.ContentUpdatedAt,
	}, nil
}
