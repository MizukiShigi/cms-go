package controller

import (
	"net/http"

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
}
