package controller

import (
	"encoding/json"
	"net/http"

	"github.com/MizukiShigi/cms-go/internal/domain/myerror"
	"github.com/MizukiShigi/cms-go/internal/presentation/helper"
	"github.com/MizukiShigi/cms-go/internal/usecase"
)

type CreatePostRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type CreatePostResponse struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type PostController struct {
	createPostUsecase *usecase.CreatePostUsecase
}

func NewPostController(createPostUsecase *usecase.CreatePostUsecase) *PostController {
	return &PostController{
		createPostUsecase: createPostUsecase,
	}
}

func (pc *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		MyError := myerror.NewMyError(myerror.InvalidRequestCode, "Invalid request payload")
		helper.RespondWithError(w, MyError)
		return
	}

	input := &usecase.CreatePostInput{
		Title:   req.Title,
		Content: req.Content,
		Tags:    req.Tags,
	}

	output, err := pc.createPostUsecase.Execute(r.Context(), input)
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	createPostResponse := CreatePostResponse{
		ID:      output.ID,
		Title:   output.Title,
		Content: output.Content,
		Tags:    output.Tags,
	}

	helper.RespondWithJSON(w, http.StatusCreated, createPostResponse)
}
