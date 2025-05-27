package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestPatchPostUsecase_Execute_UpdateTitle(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewPatchPostUsecase(mockPostRepo)

	newTitle := "Updated Title"
	input := &PatchPostInput{
		PostID: "123e4567-e89b-12d3-a456-426614174000",
		Title:  &newTitle,
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	originalTitle, _ := valueobject.NewPostTitle("Original Title")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")

	existingPost, _ := entity.NewPost(userID, originalTitle, content, status)

	mockPostRepo.On("Get", ctx, postID).Return(existingPost, nil)
	mockPostRepo.On("Update", ctx, mock.AnythingOfType("*entity.Post")).Return(nil)

	output, err := usecase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, newTitle, output.Post.Title().String())
	mockPostRepo.AssertExpectations(t)
}

func TestPatchPostUsecase_Execute_UpdateContent(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewPatchPostUsecase(mockPostRepo)

	newContent := "Updated content"
	input := &PatchPostInput{
		PostID:  "123e4567-e89b-12d3-a456-426614174000",
		Content: &newContent,
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Title")
	originalContent, _ := valueobject.NewPostContent("Original content")
	status, _ := valueobject.NewPostStatus("draft")

	existingPost, _ := entity.NewPost(userID, title, originalContent, status)

	mockPostRepo.On("Get", ctx, postID).Return(existingPost, nil)
	mockPostRepo.On("Update", ctx, mock.AnythingOfType("*entity.Post")).Return(nil)

	output, err := usecase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, newContent, output.Post.Content().String())
	mockPostRepo.AssertExpectations(t)
}

func TestPatchPostUsecase_Execute_UpdateStatus(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewPatchPostUsecase(mockPostRepo)

	newStatus := "published"
	input := &PatchPostInput{
		PostID: "123e4567-e89b-12d3-a456-426614174000",
		Status: &newStatus,
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Title")
	content, _ := valueobject.NewPostContent("Test content")
	originalStatus, _ := valueobject.NewPostStatus("draft")

	existingPost, _ := entity.NewPost(userID, title, content, originalStatus)

	mockPostRepo.On("Get", ctx, postID).Return(existingPost, nil)
	mockPostRepo.On("Update", ctx, mock.AnythingOfType("*entity.Post")).Return(nil)

	output, err := usecase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, newStatus, output.Post.Status().String())
	mockPostRepo.AssertExpectations(t)
}

func TestPatchPostUsecase_Execute_UpdateMultipleFields(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewPatchPostUsecase(mockPostRepo)

	newTitle := "Updated Title"
	newContent := "Updated content"
	newStatus := "published"
	input := &PatchPostInput{
		PostID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   &newTitle,
		Content: &newContent,
		Status:  &newStatus,
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	originalTitle, _ := valueobject.NewPostTitle("Original Title")
	originalContent, _ := valueobject.NewPostContent("Original content")
	originalStatus, _ := valueobject.NewPostStatus("draft")

	existingPost, _ := entity.NewPost(userID, originalTitle, originalContent, originalStatus)

	mockPostRepo.On("Get", ctx, postID).Return(existingPost, nil)
	mockPostRepo.On("Update", ctx, mock.AnythingOfType("*entity.Post")).Return(nil)

	output, err := usecase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, newTitle, output.Post.Title().String())
	assert.Equal(t, newContent, output.Post.Content().String())
	assert.Equal(t, newStatus, output.Post.Status().String())
	mockPostRepo.AssertExpectations(t)
}

func TestPatchPostUsecase_Execute_InvalidPostID(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewPatchPostUsecase(mockPostRepo)

	newTitle := "Updated Title"
	input := &PatchPostInput{
		PostID: "invalid-uuid",
		Title:  &newTitle,
	}

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InvalidCode, myErr.Code)
	mockPostRepo.AssertNotCalled(t, "Get")
}

func TestPatchPostUsecase_Execute_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewPatchPostUsecase(mockPostRepo)

	invalidStatus := "invalid-status"
	input := &PatchPostInput{
		PostID: "123e4567-e89b-12d3-a456-426614174000",
		Status: &invalidStatus,
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Title")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")

	existingPost, _ := entity.NewPost(userID, title, content, status)

	mockPostRepo.On("Get", ctx, postID).Return(existingPost, nil)

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InvalidCode, myErr.Code)
	mockPostRepo.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "Update")
}

func TestPatchPostUsecase_Execute_PostNotFound(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewPatchPostUsecase(mockPostRepo)

	newTitle := "Updated Title"
	input := &PatchPostInput{
		PostID: "123e4567-e89b-12d3-a456-426614174000",
		Title:  &newTitle,
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	expectedErr := &valueobject.MyError{
		Code:    valueobject.NotFoundCode,
		Message: "Post not found",
	}

	mockPostRepo.On("Get", ctx, postID).Return(nil, expectedErr)

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockPostRepo.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "Update")
}

func TestPatchPostUsecase_Execute_UpdateFails(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewPatchPostUsecase(mockPostRepo)

	newTitle := "Updated Title"
	input := &PatchPostInput{
		PostID: "123e4567-e89b-12d3-a456-426614174000",
		Title:  &newTitle,
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	originalTitle, _ := valueobject.NewPostTitle("Original Title")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")

	existingPost, _ := entity.NewPost(userID, originalTitle, content, status)

	expectedErr := &valueobject.MyError{
		Code:    valueobject.InternalServerErrorCode,
		Message: "Failed to update post",
	}

	mockPostRepo.On("Get", ctx, postID).Return(existingPost, nil)
	mockPostRepo.On("Update", ctx, mock.AnythingOfType("*entity.Post")).Return(expectedErr)

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockPostRepo.AssertExpectations(t)
}

func TestPatchPostUsecase_Execute_NoFieldsToUpdate(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewPatchPostUsecase(mockPostRepo)

	input := &PatchPostInput{
		PostID: "123e4567-e89b-12d3-a456-426614174000",
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Title")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")

	existingPost, _ := entity.NewPost(userID, title, content, status)

	mockPostRepo.On("Get", ctx, postID).Return(existingPost, nil)

	output, err := usecase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, existingPost, output.Post)
	mockPostRepo.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "Update")
}