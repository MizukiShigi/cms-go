package usecase

import (
	"context"
	"io"
	"os"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/repository"
	"github.com/MizukiShigi/cms-go/internal/domain/service"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type CreateImageInput struct {
	UserID           valueobject.UserID
	PostID           valueobject.PostID
	File             io.Reader
	OriginalFilename valueobject.ImageFilename
	SortOrder        int
}

type CreateImageOutput struct {
	ID               valueobject.ImageID
	ImageURL         string
	UserID           valueobject.UserID
	PostID           valueobject.PostID
	OriginalFilename valueobject.ImageFilename
	StoredFilename   string
	SortOrder        int
}

type CreateImageUsecase struct {
	imageRepository repository.ImageRepository
	storageService  service.StorageService
}

func NewCreateImageUsecase(imageRepository repository.ImageRepository, storageService service.StorageService) *CreateImageUsecase {
	return &CreateImageUsecase{imageRepository: imageRepository, storageService: storageService}
}

func (u *CreateImageUsecase) Execute(ctx context.Context, input *CreateImageInput) (*CreateImageOutput, error) {
	bucketName := os.Getenv("GCS_IMAGEWBUCKET_NAME")
	uploadResult, err := u.storageService.UploadImage(ctx, bucketName, input.OriginalFilename, input.File)
	if err != nil {
		return nil, err
	}

	image := entity.NewImage(input.OriginalFilename, uploadResult.StoredFilename, uploadResult.URL, input.PostID, input.UserID, input.SortOrder)
	err = u.imageRepository.Create(ctx, image)
	if err != nil {
		return nil, err
	}

	return &CreateImageOutput{
		ID:               image.ID,
		ImageURL:         image.GCSURL,
		UserID:           image.UserID,
		PostID:           image.PostID,
		OriginalFilename: image.OriginalFilename,
		StoredFilename:   image.StoredFilename,
		SortOrder:        image.SortOrder,
	}, nil
}
