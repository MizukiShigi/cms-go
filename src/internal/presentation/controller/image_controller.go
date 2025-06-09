package controller

import (
	"net/http"
	"strconv"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/MizukiShigi/cms-go/internal/presentation/helper"
	"github.com/MizukiShigi/cms-go/internal/usecase"
)

type ImageController struct {
	createImageUsecase *usecase.CreateImageUsecase
}

func NewImageController(createImageUsecase *usecase.CreateImageUsecase) *ImageController {
	return &ImageController{createImageUsecase: createImageUsecase}
}

type CreateImageResponse struct {
	ID               string `json:"id"`
	ImageURL         string `json:"image_url"`
	UserID           string `json:"user_id"`
	PostID           string `json:"post_id"`
	OriginalFilename string `json:"original_filename"`
	StoredFilename   string `json:"stored_filename"`
	SortOrder        int    `json:"sort_order"`
}

func (c *ImageController) CreateImage(w http.ResponseWriter, r *http.Request) {
	ctxUserID, err := domaincontext.GetUserID(r.Context())
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	userID, err := valueobject.ParseUserID(ctxUserID)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid user ID"))
		return
	}

	postIDStr := r.FormValue("post_id")
	if postIDStr == "" {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Post ID is required"))
		return
	}

	postID, err := valueobject.ParsePostID(postIDStr)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid post ID"))
		return
	}

	sortOrderStr := r.FormValue("sort_order")
	sortOrder := 0
	if sortOrderStr != "" {
		sortOrder, err = strconv.Atoi(sortOrderStr)
		if err != nil {
			helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid sort order"))
			return
		}
	}
	if sortOrder < 0 || sortOrder > 999 {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid sort order"))
		return
	}

	_, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "image required", http.StatusBadRequest)
		return
	}

	filenameStr := header.Filename
	filename, err := valueobject.NewImageFilename(filenameStr)
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Invalid filename"))
		return
	}

	if header.Size > 10*1024*1024 {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "File too large"))
		return
	}

    file, err := header.Open()
	if err != nil {
		helper.RespondWithError(w, valueobject.NewMyError(valueobject.InvalidCode, "Failed to open file"))
		return
	}
	defer file.Close()

	input := &usecase.CreateImageInput{
		UserID: userID,
		PostID: postID,
		File: file,
		OriginalFilename: filename,
		SortOrder: sortOrder,
	}

	output, err := c.createImageUsecase.Execute(r.Context(), input)
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	response := CreateImageResponse{
		ID:               output.ID.String(),
		ImageURL:         output.ImageURL,
		UserID:           output.UserID.String(),
		PostID:           output.PostID.String(),
		OriginalFilename: output.OriginalFilename.String(),
		StoredFilename:   output.StoredFilename,
		SortOrder:        output.SortOrder,
	}

	helper.RespondWithJSON(w, http.StatusCreated, response)
}
