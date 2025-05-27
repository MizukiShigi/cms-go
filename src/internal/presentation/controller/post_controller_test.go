package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/MizukiShigi/cms-go/internal/usecase"
)

type MockCreatePostUsecase struct {
	mock.Mock
}

func (m *MockCreatePostUsecase) Execute(ctx context.Context, input *usecase.CreatePostInput) (*usecase.CreatePostOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.CreatePostOutput), args.Error(1)
}

type MockGetPostUsecase struct {
	mock.Mock
}

func (m *MockGetPostUsecase) Execute(ctx context.Context, input *usecase.GetPostInput) (*usecase.GetPostOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.GetPostOutput), args.Error(1)
}

type MockUpdatePostUsecase struct {
	mock.Mock
}

func (m *MockUpdatePostUsecase) Execute(ctx context.Context, input *usecase.UpdatePostInput) (*usecase.UpdatePostOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.UpdatePostOutput), args.Error(1)
}

type MockPatchPostUsecase struct {
	mock.Mock
}

func (m *MockPatchPostUsecase) Execute(ctx context.Context, input *usecase.PatchPostInput) (*usecase.PatchPostOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.PatchPostOutput), args.Error(1)
}

func TestPostController_CreatePost_Success(t *testing.T) {
	mockCreateUsecase := &MockCreatePostUsecase{}
	mockGetUsecase := &MockGetPostUsecase{}
	mockUpdateUsecase := &MockUpdatePostUsecase{}
	mockPatchUsecase := &MockPatchPostUsecase{}

	controller := NewPostController(mockCreateUsecase, mockGetUsecase, mockUpdateUsecase, mockPatchUsecase)

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Post")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")
	expectedPost, _ := entity.NewPost(userID, title, content, status)

	requestBody := map[string]interface{}{
		"title":   "Test Post",
		"content": "Test content",
		"tags":    []string{"go", "testing"},
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), context.UserIDKey, "123e4567-e89b-12d3-a456-426614174000"))

	rr := httptest.NewRecorder()

	expectedOutput := &usecase.CreatePostOutput{
		Post: expectedPost,
	}

	mockCreateUsecase.On("Execute", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*usecase.CreatePostInput")).Return(expectedOutput, nil)

	controller.CreatePost(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockCreateUsecase.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "post")
}

func TestPostController_CreatePost_MissingUserID(t *testing.T) {
	mockCreateUsecase := &MockCreatePostUsecase{}
	mockGetUsecase := &MockGetPostUsecase{}
	mockUpdateUsecase := &MockUpdatePostUsecase{}
	mockPatchUsecase := &MockPatchPostUsecase{}

	controller := NewPostController(mockCreateUsecase, mockGetUsecase, mockUpdateUsecase, mockPatchUsecase)

	requestBody := map[string]interface{}{
		"title":   "Test Post",
		"content": "Test content",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	controller.CreatePost(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockCreateUsecase.AssertNotCalled(t, "Execute")
}

func TestPostController_CreatePost_InvalidJSON(t *testing.T) {
	mockCreateUsecase := &MockCreatePostUsecase{}
	mockGetUsecase := &MockGetPostUsecase{}
	mockUpdateUsecase := &MockUpdatePostUsecase{}
	mockPatchUsecase := &MockPatchPostUsecase{}

	controller := NewPostController(mockCreateUsecase, mockGetUsecase, mockUpdateUsecase, mockPatchUsecase)

	req := httptest.NewRequest("POST", "/posts", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), context.UserIDKey, "123e4567-e89b-12d3-a456-426614174000"))

	rr := httptest.NewRecorder()

	controller.CreatePost(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockCreateUsecase.AssertNotCalled(t, "Execute")
}

func TestPostController_CreatePost_UsecaseError(t *testing.T) {
	mockCreateUsecase := &MockCreatePostUsecase{}
	mockGetUsecase := &MockGetPostUsecase{}
	mockUpdateUsecase := &MockUpdatePostUsecase{}
	mockPatchUsecase := &MockPatchPostUsecase{}

	controller := NewPostController(mockCreateUsecase, mockGetUsecase, mockUpdateUsecase, mockPatchUsecase)

	requestBody := map[string]interface{}{
		"title":   "Test Post",
		"content": "Test content",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), context.UserIDKey, "123e4567-e89b-12d3-a456-426614174000"))

	rr := httptest.NewRecorder()

	expectedErr := &valueobject.MyError{
		Code:    valueobject.InvalidCode,
		Message: "Invalid title",
	}

	mockCreateUsecase.On("Execute", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*usecase.CreatePostInput")).Return(nil, expectedErr)

	controller.CreatePost(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockCreateUsecase.AssertExpectations(t)
}

func TestPostController_GetPost_Success(t *testing.T) {
	mockCreateUsecase := &MockCreatePostUsecase{}
	mockGetUsecase := &MockGetPostUsecase{}
	mockUpdateUsecase := &MockUpdatePostUsecase{}
	mockPatchUsecase := &MockPatchPostUsecase{}

	controller := NewPostController(mockCreateUsecase, mockGetUsecase, mockUpdateUsecase, mockPatchUsecase)

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Test Post")
	content, _ := valueobject.NewPostContent("Test content")
	status, _ := valueobject.NewPostStatus("draft")
	expectedPost, _ := entity.NewPost(userID, title, content, status)

	req := httptest.NewRequest("GET", "/posts/123e4567-e89b-12d3-a456-426614174000", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})

	rr := httptest.NewRecorder()

	expectedOutput := &usecase.GetPostOutput{
		Post: expectedPost,
	}

	mockGetUsecase.On("Execute", mock.AnythingOfType("*context.backgroundCtx"), mock.AnythingOfType("*usecase.GetPostInput")).Return(expectedOutput, nil)

	controller.GetPost(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockGetUsecase.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "post")
}

func TestPostController_GetPost_MissingID(t *testing.T) {
	mockCreateUsecase := &MockCreatePostUsecase{}
	mockGetUsecase := &MockGetPostUsecase{}
	mockUpdateUsecase := &MockUpdatePostUsecase{}
	mockPatchUsecase := &MockPatchPostUsecase{}

	controller := NewPostController(mockCreateUsecase, mockGetUsecase, mockUpdateUsecase, mockPatchUsecase)

	req := httptest.NewRequest("GET", "/posts/", nil)
	rr := httptest.NewRecorder()

	controller.GetPost(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockGetUsecase.AssertNotCalled(t, "Execute")
}

func TestPostController_GetPost_NotFound(t *testing.T) {
	mockCreateUsecase := &MockCreatePostUsecase{}
	mockGetUsecase := &MockGetPostUsecase{}
	mockUpdateUsecase := &MockUpdatePostUsecase{}
	mockPatchUsecase := &MockPatchPostUsecase{}

	controller := NewPostController(mockCreateUsecase, mockGetUsecase, mockUpdateUsecase, mockPatchUsecase)

	req := httptest.NewRequest("GET", "/posts/123e4567-e89b-12d3-a456-426614174000", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})

	rr := httptest.NewRecorder()

	expectedErr := &valueobject.MyError{
		Code:    valueobject.NotFoundCode,
		Message: "Post not found",
	}

	mockGetUsecase.On("Execute", mock.AnythingOfType("*context.backgroundCtx"), mock.AnythingOfType("*usecase.GetPostInput")).Return(nil, expectedErr)

	controller.GetPost(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockGetUsecase.AssertExpectations(t)
}

func TestPostController_UpdatePost_Success(t *testing.T) {
	mockCreateUsecase := &MockCreatePostUsecase{}
	mockGetUsecase := &MockGetPostUsecase{}
	mockUpdateUsecase := &MockUpdatePostUsecase{}
	mockPatchUsecase := &MockPatchPostUsecase{}

	controller := NewPostController(mockCreateUsecase, mockGetUsecase, mockUpdateUsecase, mockPatchUsecase)

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Updated Post")
	content, _ := valueobject.NewPostContent("Updated content")
	status, _ := valueobject.NewPostStatus("published")
	expectedPost, _ := entity.NewPost(userID, title, content, status)

	requestBody := map[string]interface{}{
		"title":   "Updated Post",
		"content": "Updated content",
		"status":  "published",
		"tags":    []string{"go", "testing", "update"},
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("PUT", "/posts/123e4567-e89b-12d3-a456-426614174000", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})

	rr := httptest.NewRecorder()

	expectedOutput := &usecase.UpdatePostOutput{
		Post: expectedPost,
	}

	mockUpdateUsecase.On("Execute", mock.AnythingOfType("*context.backgroundCtx"), mock.AnythingOfType("*usecase.UpdatePostInput")).Return(expectedOutput, nil)

	controller.UpdatePost(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockUpdateUsecase.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "post")
}

func TestPostController_PatchPost_Success(t *testing.T) {
	mockCreateUsecase := &MockCreatePostUsecase{}
	mockGetUsecase := &MockGetPostUsecase{}
	mockUpdateUsecase := &MockUpdatePostUsecase{}
	mockPatchUsecase := &MockPatchPostUsecase{}

	controller := NewPostController(mockCreateUsecase, mockGetUsecase, mockUpdateUsecase, mockPatchUsecase)

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Patched Post")
	content, _ := valueobject.NewPostContent("Original content")
	status, _ := valueobject.NewPostStatus("draft")
	expectedPost, _ := entity.NewPost(userID, title, content, status)

	requestBody := map[string]interface{}{
		"title": "Patched Post",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("PATCH", "/posts/123e4567-e89b-12d3-a456-426614174000", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})

	rr := httptest.NewRecorder()

	expectedOutput := &usecase.PatchPostOutput{
		Post: expectedPost,
	}

	mockPatchUsecase.On("Execute", mock.AnythingOfType("*context.backgroundCtx"), mock.AnythingOfType("*usecase.PatchPostInput")).Return(expectedOutput, nil)

	controller.PatchPost(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockPatchUsecase.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "post")
}

func TestPostController_PatchPost_EmptyBody(t *testing.T) {
	mockCreateUsecase := &MockCreatePostUsecase{}
	mockGetUsecase := &MockGetPostUsecase{}
	mockUpdateUsecase := &MockUpdatePostUsecase{}
	mockPatchUsecase := &MockPatchPostUsecase{}

	controller := NewPostController(mockCreateUsecase, mockGetUsecase, mockUpdateUsecase, mockPatchUsecase)

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	title, _ := valueobject.NewPostTitle("Original Post")
	content, _ := valueobject.NewPostContent("Original content")
	status, _ := valueobject.NewPostStatus("draft")
	expectedPost, _ := entity.NewPost(userID, title, content, status)

	requestBody := map[string]interface{}{}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("PATCH", "/posts/123e4567-e89b-12d3-a456-426614174000", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})

	rr := httptest.NewRecorder()

	expectedOutput := &usecase.PatchPostOutput{
		Post: expectedPost,
	}

	mockPatchUsecase.On("Execute", mock.AnythingOfType("*context.backgroundCtx"), mock.AnythingOfType("*usecase.PatchPostInput")).Return(expectedOutput, nil)

	controller.PatchPost(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockPatchUsecase.AssertExpectations(t)
}