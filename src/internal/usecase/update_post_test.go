package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestUpdatePostUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewUpdatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &UpdatePostInput{
		PostID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "Updated Title",
		Content: "Updated content",
		Status:  "published",
		Tags:    []string{"go", "testing", "update"},
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	originalTitle, _ := valueobject.NewPostTitle("Original Title")
	originalContent, _ := valueobject.NewPostContent("Original content")
	originalStatus, _ := valueobject.NewPostStatus("draft")

	existingPost, _ := entity.NewPost(userID, originalTitle, originalContent, originalStatus)

	mockTag1 := &entity.Tag{}
	mockTag2 := &entity.Tag{}
	mockTag3 := &entity.Tag{}

	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockPostRepo.On("Get", ctx, postID).Return(existingPost, nil)
	mockPostRepo.On("Update", ctx, mock.AnythingOfType("*entity.Post")).Return(nil)
	mockTagRepo.On("FindOrCreateByName", ctx, mock.AnythingOfType("*entity.Tag")).Return(mockTag1, nil).Once()
	mockTagRepo.On("FindOrCreateByName", ctx, mock.AnythingOfType("*entity.Tag")).Return(mockTag2, nil).Once()
	mockTagRepo.On("FindOrCreateByName", ctx, mock.AnythingOfType("*entity.Tag")).Return(mockTag3, nil).Once()
	mockPostRepo.On("SetTags", ctx, mock.AnythingOfType("*entity.Post"), mock.AnythingOfType("[]*entity.Tag")).Return(nil)

	output, err := usecase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, input.Title, output.Post.Title().String())
	assert.Equal(t, input.Content, output.Post.Content().String())
	assert.Equal(t, input.Status, output.Post.Status().String())
	mockTxManager.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
}

func TestUpdatePostUsecase_Execute_WithoutTags(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewUpdatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &UpdatePostInput{
		PostID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "Updated Title",
		Content: "Updated content",
		Status:  "published",
		Tags:    []string{},
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	originalTitle, _ := valueobject.NewPostTitle("Original Title")
	originalContent, _ := valueobject.NewPostContent("Original content")
	originalStatus, _ := valueobject.NewPostStatus("draft")

	existingPost, _ := entity.NewPost(userID, originalTitle, originalContent, originalStatus)

	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockPostRepo.On("Get", ctx, postID).Return(existingPost, nil)
	mockPostRepo.On("Update", ctx, mock.AnythingOfType("*entity.Post")).Return(nil)
	mockPostRepo.On("SetTags", ctx, mock.AnythingOfType("*entity.Post"), mock.AnythingOfType("[]*entity.Tag")).Return(nil)

	output, err := usecase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockTxManager.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTagRepo.AssertNotCalled(t, "FindOrCreateByName")
}

func TestUpdatePostUsecase_Execute_InvalidPostID(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewUpdatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &UpdatePostInput{
		PostID:  "invalid-uuid",
		Title:   "Updated Title",
		Content: "Updated content",
		Status:  "published",
		Tags:    []string{},
	}

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InvalidCode, myErr.Code)
	mockTxManager.AssertNotCalled(t, "WithinTransaction")
}

func TestUpdatePostUsecase_Execute_InvalidTitle(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewUpdatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &UpdatePostInput{
		PostID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "",
		Content: "Updated content",
		Status:  "published",
		Tags:    []string{},
	}

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InvalidCode, myErr.Code)
	mockTxManager.AssertNotCalled(t, "WithinTransaction")
}

func TestUpdatePostUsecase_Execute_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewUpdatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &UpdatePostInput{
		PostID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "Updated Title",
		Content: "Updated content",
		Status:  "invalid-status",
		Tags:    []string{},
	}

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InvalidCode, myErr.Code)
	mockTxManager.AssertNotCalled(t, "WithinTransaction")
}

func TestUpdatePostUsecase_Execute_PostNotFound(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewUpdatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &UpdatePostInput{
		PostID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "Updated Title",
		Content: "Updated content",
		Status:  "published",
		Tags:    []string{},
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	expectedErr := &valueobject.MyError{
		Code:    valueobject.NotFoundCode,
		Message: "Post not found",
	}

	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(expectedErr)
	mockPostRepo.On("Get", ctx, postID).Return(nil, expectedErr)

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockTxManager.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
}

func TestUpdatePostUsecase_Execute_UpdateFails(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewUpdatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &UpdatePostInput{
		PostID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "Updated Title",
		Content: "Updated content",
		Status:  "published",
		Tags:    []string{},
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	originalTitle, _ := valueobject.NewPostTitle("Original Title")
	originalContent, _ := valueobject.NewPostContent("Original content")
	originalStatus, _ := valueobject.NewPostStatus("draft")

	existingPost, _ := entity.NewPost(userID, originalTitle, originalContent, originalStatus)

	expectedErr := &valueobject.MyError{
		Code:    valueobject.InternalServerErrorCode,
		Message: "Failed to update post",
	}

	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(expectedErr)
	mockPostRepo.On("Get", ctx, postID).Return(existingPost, nil)
	mockPostRepo.On("Update", ctx, mock.AnythingOfType("*entity.Post")).Return(expectedErr)

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockTxManager.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
}

func TestUpdatePostUsecase_Execute_TagCreationFails(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewUpdatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &UpdatePostInput{
		PostID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "Updated Title",
		Content: "Updated content",
		Status:  "published",
		Tags:    []string{"go"},
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	originalTitle, _ := valueobject.NewPostTitle("Original Title")
	originalContent, _ := valueobject.NewPostContent("Original content")
	originalStatus, _ := valueobject.NewPostStatus("draft")

	existingPost, _ := entity.NewPost(userID, originalTitle, originalContent, originalStatus)

	expectedErr := &valueobject.MyError{
		Code:    valueobject.InternalServerErrorCode,
		Message: "Failed to create tag",
	}

	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(expectedErr)
	mockPostRepo.On("Get", ctx, postID).Return(existingPost, nil)
	mockPostRepo.On("Update", ctx, mock.AnythingOfType("*entity.Post")).Return(nil)
	mockTagRepo.On("FindOrCreateByName", ctx, mock.AnythingOfType("*entity.Tag")).Return(nil, expectedErr)

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockTxManager.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
}