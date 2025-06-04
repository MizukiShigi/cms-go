package service

import (
	"context"
	"io"
)

type StorageService interface {
	Upload(ctx context.Context, bucketName string, path string, data io.Reader) (string, error)
}
