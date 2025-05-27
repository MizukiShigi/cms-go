package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestGetPostUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewGetPostUsecase(mockPostRepo)

	input := &GetPostInput{
		PostID: "123e4567-e89b-12d3-a456-426614174000",
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	userID, _ := valueobject.NewUserID("456e7890-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Post")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")

	expectedPost, _ := entity.NewPost(userID, title, content, status)

	mockPostRepo.On("Get", ctx, postID).Return(expectedPost, nil)

	output, err := usecase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, expectedPost, output.Post)
	mockPostRepo.AssertExpectations(t)
}

func TestGetPostUsecase_Execute_InvalidPostID(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewGetPostUsecase(mockPostRepo)

	input := &GetPostInput{
		PostID: "invalid-uuid",
	}

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InvalidCode, myErr.Code)
	mockPostRepo.AssertNotCalled(t, "Get")
}

func TestGetPostUsecase_Execute_PostNotFound(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewGetPostUsecase(mockPostRepo)

	input := &GetPostInput{
		PostID: "123e4567-e89b-12d3-a456-426614174000",
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
}

func TestGetPostUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewGetPostUsecase(mockPostRepo)

	input := &GetPostInput{
		PostID: "123e4567-e89b-12d3-a456-426614174000",
	}

	postID, _ := valueobject.NewPostID(input.PostID)
	expectedErr := &valueobject.MyError{
		Code:    valueobject.InternalServerErrorCode,
		Message: "Database error",
	}

	mockPostRepo.On("Get", ctx, postID).Return(nil, expectedErr)

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockPostRepo.AssertExpectations(t)
}

func TestGetPostUsecase_Execute_EmptyPostID(t *testing.T) {
	ctx := context.Background()
	mockPostRepo := &MockPostRepository{}

	usecase := NewGetPostUsecase(mockPostRepo)

	input := &GetPostInput{
		PostID: "",
	}

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InvalidCode, myErr.Code)
	mockPostRepo.AssertNotCalled(t, "Get")
}