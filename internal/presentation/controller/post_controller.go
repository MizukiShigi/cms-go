package controller

import (
	"encoding/json"
	"net/http"
	"time"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/MizukiShigi/cms-go/internal/presentation/helper"
	"github.com/MizukiShigi/cms-go/internal/usecase"

	"github.com/go-playground/validator/v10"
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
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content" validate:"required"`
	Tags    []string `json:"tags"`
}

type CreatePostResponse struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (pc *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctxUserID, err := domaincontext.GetUserID(r.Context())
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	var req CreatePostRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		MyError := valueobject.NewMyError(valueobject.InvalidCode, "Invalid request payload")
		helper.RespondWithError(w, MyError)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			myError := valueobject.NewMyError(valueobject.InvalidCode, err.Error())
			helper.RespondWithError(w, myError)
			return
		}
	}

	title, err := valueobject.NewPostTitle(req.Title)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid title"))
		return
	}

	content, err := valueobject.NewPostContent(req.Content)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid content"))
		return
	}

	userID, err := valueobject.ParseUserID(ctxUserID)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid user ID"))
		return
	}

	var inputTags []valueobject.TagName
	for _, tag := range req.Tags {
		tag, err := valueobject.NewTagName(tag)
		if err != nil {
			helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid tag"))
			return
		}
		inputTags = append(inputTags, tag)
	}

	input := &usecase.CreatePostInput{
		Title:   title,
		Content: content,
		Tags:    inputTags,
		UserID:  userID,
	}

	output, err := pc.createPostUsecase.Execute(r.Context(), input)
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	oiutputTags := []string{}
	for _, tag := range output.Tags {
		oiutputTags = append(oiutputTags, tag.String())
	}

	createPostResponse := CreatePostResponse{
		ID:      output.ID.String(),
		Title:   output.Title.String(),
		Content: output.Content.String(),
		Tags:    oiutputTags,
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
	id, exists := vars["id"]
	if !exists {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Required post ID"))
		return
	}

	postID, err := valueobject.ParsePostID(id)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid post ID"))
		return
	}

	input := &usecase.GetPostInput{ID: postID}
	output, err := pc.getPostUsecase.Execute(r.Context(), input)
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	// 配列で返したいので、nilの場合は空配列を返す
	tags := []string{}
	if output.Tags != nil {
		tags = make([]string, 0, len(output.Tags))
		for _, tag := range output.Tags {
			tags = append(tags, tag.String())
		}
	}

	res := GetPostResponse{
		ID:               output.ID.String(),
		Title:            output.Title.String(),
		Content:          output.Content.String(),
		Status:           output.Status.String(),
		Tags:             tags,
		FirstPublishedAt: output.FirstPublishedAt,
		ContentUpdatedAt: output.ContentUpdatedAt,
	}

	helper.RespondWithJSON(w, http.StatusOK, res)
}
