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
	listPostsUsecase  *usecase.ListPostsUsecase
	createPostUsecase *usecase.CreatePostUsecase
	getPostUsecase    *usecase.GetPostUsecase
	updatePostUsecase *usecase.UpdatePostUsecase
	patchPostUsecase  *usecase.PatchPostUsecase
}

func NewPostController(listPostsUsecase *usecase.ListPostsUsecase, createPostUsecase *usecase.CreatePostUsecase, getPostUsecase *usecase.GetPostUsecase, updatePostUsecase *usecase.UpdatePostUsecase, patchPostUsecase *usecase.PatchPostUsecase) *PostController {
	return &PostController{
		listPostsUsecase:  listPostsUsecase,
		createPostUsecase: createPostUsecase,
		getPostUsecase:    getPostUsecase,
		updatePostUsecase: updatePostUsecase,
		patchPostUsecase:  patchPostUsecase,
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
	Status  string   `json:"status" validate:"required,oneof=draft published"`
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

	status, err := valueobject.NewPostStatus(req.Status)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid status"))
		return
	}

	input := &usecase.CreatePostInput{
		Title:   title,
		Content: content,
		Tags:    inputTags,
		UserID:  userID,
		Status:  status,
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

type UpdatePostRequest struct {
	Title   string   `json:"title" validate:"required,min=1"`
	Content string   `json:"content" validate:"required,min=1"`
	Tags    []string `json:"tags"`
}

type UpdatePostResponse struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	Status           string     `json:"status"`
	Tags             []string   `json:"tags"`
	FirstPublishedAt *time.Time `json:"first_published_at"`
	ContentUpdatedAt *time.Time `json:"content_updated_at"`
}

func (pc *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, exists := vars["id"]
	if !exists {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Required post ID"))
		return
	}

	var req UpdatePostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid request payload"))
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

	postID, err := valueobject.ParsePostID(id)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid post ID"))
		return
	}

	content, err := valueobject.NewPostContent(req.Content)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid content"))
		return
	}

	title, err := valueobject.NewPostTitle(req.Title)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid title"))
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

	input := &usecase.UpdatePostInput{
		ID:      postID,
		Title:   title,
		Content: content,
		Tags:    inputTags,
	}

	output, err := pc.updatePostUsecase.Execute(r.Context(), input)
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	tags := []string{}
	if output.Tags != nil {
		tags = make([]string, 0, len(output.Tags))
		for _, tag := range output.Tags {
			tags = append(tags, tag.String())
		}
	}

	res := UpdatePostResponse{
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

type PatchPostRequest struct {
	Title   string   `json:"title" validate:"omitempty,min=1"`
	Content string   `json:"content" validate:"omitempty,min=1"`
	Tags    []string `json:"tags"`
	Status  string   `json:"status" validate:"omitempty,oneof=draft published private deleted"`
}

type PatchPostResponse struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	Status           string     `json:"status"`
	Tags             []string   `json:"tags"`
	FirstPublishedAt *time.Time `json:"first_published_at"`
	ContentUpdatedAt *time.Time `json:"content_updated_at"`
}

func (pc *PostController) PatchPost(w http.ResponseWriter, r *http.Request) {
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

	var req PatchPostRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid request payload"))
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

	if req.Title == "" && req.Content == "" && req.Status == "" && len(req.Tags) == 0 {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "No update fields"))
		return
	}

	var title *valueobject.PostTitle
	if req.Title != "" {
		t, err := valueobject.NewPostTitle(req.Title)
		if err != nil {
			helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid title"))
			return
		}
		title = &t
	}

	var content *valueobject.PostContent
	if req.Content != "" {
		c, err := valueobject.NewPostContent(req.Content)
		if err != nil {
			helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid content"))
			return
		}
		content = &c
	}

	var status *valueobject.PostStatus
	if req.Status != "" {
		s, err := valueobject.NewPostStatus(req.Status)
		if err != nil {
			helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid status"))
			return
		}
		status = &s
	}

	var inputTags []valueobject.TagName
	if len(req.Tags) > 0 {
		tags := make([]valueobject.TagName, 0, len(req.Tags))
		for _, tag := range req.Tags {
			t, err := valueobject.NewTagName(tag)
			if err != nil {
				helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid tag"))
				return
			}
			tags = append(tags, t)
		}
		inputTags = tags
	}

	input := &usecase.PatchPostInput{
		ID:      postID,
		Title:   title,
		Content: content,
		Status:  status,
		Tags:    inputTags,
	}

	output, err := pc.patchPostUsecase.Execute(r.Context(), input)
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	outputTags := []string{}
	if output.Tags != nil {
		outputTags = make([]string, 0, len(output.Tags))
		for _, tag := range output.Tags {
			outputTags = append(outputTags, tag.String())
		}
	}

	res := PatchPostResponse{
		ID:               output.ID.String(),
		Title:            output.Title.String(),
		Content:          output.Content.String(),
		Status:           output.Status.String(),
		Tags:             outputTags,
		FirstPublishedAt: output.FirstPublishedAt,
		ContentUpdatedAt: output.ContentUpdatedAt,
	}

	helper.RespondWithJSON(w, http.StatusOK, res)
}

func (pc *PostController) ListPosts(w http.ResponseWriter, r *http.Request) {
	// ユーザーIDをコンテキストから取得
	userIDValue := r.Context().Value(domaincontext.UserID)
	userIDStr, ok := userIDValue.(string)
	if !ok {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.UnauthorizedCode, "User not authenticated"))
		return
	}

	userID, err := valueobject.ParseUserID(userIDStr)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid user ID"))
		return
	}

	// クエリパラメータを取得
	query := r.URL.Query()
	req := &usecase.ListPostsRequest{
		UserID: userID,
		Limit:  query.Get("limit"),
		Offset: query.Get("offset"),
		Status: query.Get("status"),
		Sort:   query.Get("sort"),
	}

	// ユースケース実行
	response, err := pc.listPostsUsecase.Execute(r.Context(), req)
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, response)
}
