package repository

import (
	"context"
	"database/sql"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type PostTagRepository struct {
	db *sql.DB
}

func NewPostTagRepository(db *sql.DB) *PostTagRepository {
	return &PostTagRepository{db: db}
}

func Create(ctx context.Context, postID valueobject.PostID, TagID valueobject.TagID) error {
	
}
