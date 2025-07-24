package usecase

import (
	"context"
	"strconv"

	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/repository"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type ListPostsUsecase struct {
	postRepository repository.PostRepository
}

func NewListPostsUsecase(postRepository repository.PostRepository) *ListPostsUsecase {
	return &ListPostsUsecase{
		postRepository: postRepository,
	}
}

// ListPostsRequest は投稿一覧取得のリクエスト
type ListPostsRequest struct {
	Limit  string
	Offset string
	Status string
	Sort   string
}

// ListPostsResponse は投稿一覧取得のレスポンス
type ListPostsResponse struct {
	Posts []*PostSummary `json:"posts"`
	Meta  *PaginationMeta `json:"meta"`
}

// PostSummary は投稿の概要情報
type PostSummary struct {
	ID                string   `json:"id"`
	Title             string   `json:"title"`
	Status            string   `json:"status"`
	Tags              []string `json:"tags"`
	FirstPublishedAt  *string  `json:"first_published_at"`
	ContentUpdatedAt  *string  `json:"content_updated_at"`
}

// PaginationMeta はページネーション情報
type PaginationMeta struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasNext bool `json:"has_next"`
}

func (u *ListPostsUsecase) Execute(ctx context.Context, req *ListPostsRequest) (*ListPostsResponse, error) {
	// パラメータのバリデーションとデフォルト値設定
	limit := 20
	if req.Limit != "" {
		if l, err := strconv.Atoi(req.Limit); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if req.Offset != "" {
		if o, err := strconv.Atoi(req.Offset); err == nil && o >= 0 {
			offset = o
		}
	}

	sort := "created_at_desc"
	if req.Sort != "" {
		switch req.Sort {
		case "created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc":
			sort = req.Sort
		}
	}

	// ステータスフィルタの処理
	var status *valueobject.PostStatus
	if req.Status != "" {
		if s, err := valueobject.NewPostStatus(req.Status); err == nil {
			status = &s
		}
	}

	// リポジトリオプション作成
	options := &repository.ListPostsOptions{
		Limit:  limit,
		Offset: offset,
		Status: status,
		Sort:   sort,
	}

	// 投稿一覧取得
	posts, total, err := u.postRepository.List(ctx, options)
	if err != nil {
		return nil, err
	}

	// レスポンス作成
	summaries := make([]*PostSummary, 0, len(posts))
	for _, post := range posts {
		summaries = append(summaries, u.convertToSummary(post))
	}

	meta := &PaginationMeta{
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasNext: offset+limit < total,
	}

	return &ListPostsResponse{
		Posts: summaries,
		Meta:  meta,
	}, nil
}

func (u *ListPostsUsecase) convertToSummary(post *entity.Post) *PostSummary {
	tags := make([]string, 0, len(post.Tags))
	for _, tag := range post.Tags {
		tags = append(tags, tag.String())
	}

	var firstPublishedAt *string
	if post.FirstPublishedAt != nil {
		iso := post.FirstPublishedAt.Format("2006-01-02T15:04:05Z")
		firstPublishedAt = &iso
	}

	var contentUpdatedAt *string
	if post.ContentUpdatedAt != nil {
		iso := post.ContentUpdatedAt.Format("2006-01-02T15:04:05Z")
		contentUpdatedAt = &iso
	}

	return &PostSummary{
		ID:               post.ID.String(),
		Title:            post.Title.String(),
		Status:           post.Status.String(),
		Tags:             tags,
		FirstPublishedAt: firstPublishedAt,
		ContentUpdatedAt: contentUpdatedAt,
	}
}