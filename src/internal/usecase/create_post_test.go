package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	if fn != nil {
		return fn(ctx)
	}
	return args.Error(0)
}

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(ctx context.Context, post *entity.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) Get(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostRepository) Update(ctx context.Context, post *entity.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) SetTags(ctx context.Context, post *entity.Post, tags []*entity.Tag) error {
	args := m.Called(ctx, post, tags)
	return args.Error(0)
}

type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) FindByPostID(ctx context.Context, postID valueobject.PostID) ([]*entity.Tag, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Tag), args.Error(1)
}

func (m *MockTagRepository) FindOrCreateByName(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	args := m.Called(ctx, tag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tag), args.Error(1)
}

func TestCreatePostUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewCreatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &CreatePostInput{
		UserID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "Test Post",
		Content: "This is a test post content",
		Tags:    []string{"go", "testing"},
	}

	userID, _ := valueobject.NewUserID(input.UserID)
	mockPost := &entity.Post{}

	mockTag1 := &entity.Tag{}
	mockTag2 := &entity.Tag{}

	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockPostRepo.On("Create", ctx, mock.AnythingOfType("*entity.Post")).Return(nil)
	mockTagRepo.On("FindOrCreateByName", ctx, mock.AnythingOfType("*entity.Tag")).Return(mockTag1, nil).Once()
	mockTagRepo.On("FindOrCreateByName", ctx, mock.AnythingOfType("*entity.Tag")).Return(mockTag2, nil).Once()
	mockPostRepo.On("SetTags", ctx, mock.AnythingOfType("*entity.Post"), mock.AnythingOfType("[]*entity.Tag")).Return(nil)

	output, err := usecase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, userID, output.Post.UserID())
	mockTxManager.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
}

func TestCreatePostUsecase_Execute_WithoutTags(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewCreatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &CreatePostInput{
		UserID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "Test Post",
		Content: "This is a test post content",
		Tags:    []string{},
	}

	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockPostRepo.On("Create", ctx, mock.AnythingOfType("*entity.Post")).Return(nil)

	output, err := usecase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockTxManager.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTagRepo.AssertNotCalled(t, "FindOrCreateByName")
	mockPostRepo.AssertNotCalled(t, "SetTags")
}

func TestCreatePostUsecase_Execute_InvalidUserID(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewCreatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &CreatePostInput{
		UserID:  "invalid-uuid",
		Title:   "Test Post",
		Content: "This is a test post content",
		Tags:    []string{},
	}

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InvalidCode, myErr.Code)
}

func TestCreatePostUsecase_Execute_InvalidTitle(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewCreatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &CreatePostInput{
		UserID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "",
		Content: "This is a test post content",
		Tags:    []string{},
	}

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	myErr, ok := err.(*valueobject.MyError)
	assert.True(t, ok)
	assert.Equal(t, valueobject.InvalidCode, myErr.Code)
}

func TestCreatePostUsecase_Execute_PostCreationFails(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewCreatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &CreatePostInput{
		UserID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "Test Post",
		Content: "This is a test post content",
		Tags:    []string{},
	}

	expectedErr := &valueobject.MyError{
		Code:    valueobject.InternalServerErrorCode,
		Message: "Failed to create post",
	}

	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(expectedErr)
	mockPostRepo.On("Create", ctx, mock.AnythingOfType("*entity.Post")).Return(expectedErr)

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockTxManager.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
}

func TestCreatePostUsecase_Execute_TagCreationFails(t *testing.T) {
	ctx := context.Background()
	mockTxManager := &MockTransactionManager{}
	mockPostRepo := &MockPostRepository{}
	mockTagRepo := &MockTagRepository{}

	usecase := NewCreatePostUsecase(mockTxManager, mockPostRepo, mockTagRepo)

	input := &CreatePostInput{
		UserID:  "123e4567-e89b-12d3-a456-426614174000",
		Title:   "Test Post",
		Content: "This is a test post content",
		Tags:    []string{"go"},
	}

	expectedErr := &valueobject.MyError{
		Code:    valueobject.InternalServerErrorCode,
		Message: "Failed to create tag",
	}

	mockTxManager.On("WithinTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(expectedErr)
	mockPostRepo.On("Create", ctx, mock.AnythingOfType("*entity.Post")).Return(nil)
	mockTagRepo.On("FindOrCreateByName", ctx, mock.AnythingOfType("*entity.Tag")).Return(nil, expectedErr)

	output, err := usecase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockTxManager.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
}