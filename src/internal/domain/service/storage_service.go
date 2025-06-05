package service

import (
	"context"
	"io"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type UploadResult struct {
	StoredFilename string
	URL            string
}

type StorageService interface {
	UploadImage(ctx context.Context, bucketName string, fileName valueobject.ImageFilename, data io.Reader) (UploadResult, error)
}
