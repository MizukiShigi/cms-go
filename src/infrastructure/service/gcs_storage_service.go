package service

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

type storageService struct {
	client *storage.Client
}

func NewStorageService(client *storage.Client) *storageService {
	return &storageService{client: client}
}

func (s *storageService) Upload(ctx context.Context, bucketName string, path string, data io.Reader) (string, error) {
	return "", nil
}
