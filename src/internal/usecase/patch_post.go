package usecase

import (
	"context"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/repository"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type PatchPostInput struct {
	ID      valueobject.PostID
	Title   *valueobject.PostTitle
	Content *valueobject.PostContent
	Status  *valueobject.PostStatus
	Tags    []valueobject.TagName
}

type PatchPostOutput struct {
	ID               valueobject.PostID
	Title            valueobject.PostTitle
	Content          valueobject.PostContent
	Status           valueobject.PostStatus
	Tags             []valueobject.TagName
	FirstPublishedAt *time.Time
	ContentUpdatedAt *time.Time
}

type PatchPostUsecase struct {
	postRepository repository.PostRepository
}

func NewPatchPostUsecase(postRepository repository.PostRepository) *PatchPostUsecase {
	return &PatchPostUsecase{
		postRepository: postRepository,
	}
}

func (u *PatchPostUsecase) Execute(ctx context.Context, input *PatchPostInput) (*PatchPostOutput, error) {
	post, err := u.postRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		post.Title = *input.Title
	}

	if input.Content != nil {
		post.Content = *input.Content
	}

	if input.Status != nil {
		if err := post.SetStatus(*input.Status); err != nil {
			return nil, err
		}
	}

	if len(input.Tags) > 0 {
		post.Tags = input.Tags
	}

	if err := u.postRepository.Update(ctx, post); err != nil {
		return nil, valueobject.ReturnMyError(err, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to update post"))
	}

	updatePost, err := u.postRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, valueobject.ReturnMyError(err, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Failed to get post"))
	}
	return &PatchPostOutput{
		ID:               updatePost.ID,
		Title:            updatePost.Title,
		Content:          updatePost.Content,
		Status:           updatePost.Status,
		Tags:             updatePost.Tags,
		FirstPublishedAt: updatePost.FirstPublishedAt,
		ContentUpdatedAt: updatePost.ContentUpdatedAt,
	}, nil
}
