package controller

import (
	"encoding/json"
	"net/http"
	"time"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/MizukiShigi/cms-go/internal/domain/myerror"
	"github.com/MizukiShigi/cms-go/internal/presentation/helper"
	"github.com/MizukiShigi/cms-go/internal/usecase"
	"github.com/gorilla/mux"
)

type PostController struct {
	createPostUsecase *usecase.CreatePostUsecase
	getPostUsecase    *usecase.GetPostUsecase
}

func NewPostController(createPostUsecase *usecase.CreatePostUsecase, getPostUsecase *usecase.GetPostUsecase) *PostController {
	return &PostController{
		createPostUsecase: createPostUsecase,
		getPostUsecase:    getPostUsecase,
	}
}

type Post struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

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

func (pc *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := domaincontext.GetUserID(r.Context())
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	var req CreatePostRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		MyError := myerror.NewMyError(myerror.InvalidCode, "Invalid request payload")
		helper.RespondWithError(w, MyError)
		return
	}

	input := &usecase.CreatePostInput{
		Title:   req.Title,
		Content: req.Content,
		Tags:    req.Tags,
		UserID:  userID,
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

type GetPostResponse struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	Status           string     `json:"status"`
	Tags             []string   `json:"tags"`
	FirstPublishedAt *time.Time `json:"first_published_at"`
	ContentUpdatedAt *time.Time `json:"content_updated_at"`
}

func (pc *PostController) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	input := &usecase.GetPostInput{ID: id}
	output, err := pc.getPostUsecase.Execute(r.Context(), input)
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	// 配列で返したいので、nilの場合は空配列を返す
	tags := []string{}
	if output.Tags != nil {
		tags = output.Tags
	}

	res := GetPostResponse{
		ID:               output.ID,
		Title:            output.Title,
		Content:          output.Content,
		Status:           output.Status,
		Tags:             tags,
		FirstPublishedAt: output.FirstPublishedAt,
		ContentUpdatedAt: output.ContentUpdatedAt,
	}

	helper.RespondWithJSON(w, http.StatusOK, res)
}
