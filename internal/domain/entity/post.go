package entity

import (
	"slices"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/myerror"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type Post struct {
	ID               valueobject.PostID
	Title            valueobject.PostTitle
	Content          valueobject.PostContent
	Status           valueobject.PostStatus
	FirstPublishedAt time.Time
	ContentUpdatedAt time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Tags             []valueobject.Tag
}

// 新規投稿作成
func NewPost(title valueobject.PostTitle, content valueobject.PostContent) (*Post, error) {
	now := time.Now()
	post := &Post{
		ID:               valueobject.NewPostID(),
		Title:            title,
		Content:          content,
		Status:           valueobject.StatusDraft,
		ContentUpdatedAt: now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	return post, nil
}

// 投稿データ再構築
func ParsePost(
	id valueobject.PostID,
	title valueobject.PostTitle,
	content valueobject.PostContent,
	status valueobject.PostStatus,
	firstPublishedAt time.Time,
	contentUpdatedAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *Post {
	return &Post{
		ID:               id,
		Title:            title,
		Content:          content,
		Status:           status,
		FirstPublishedAt: firstPublishedAt,
		ContentUpdatedAt: contentUpdatedAt,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

func (p *Post) AddTag(tag valueobject.Tag) error {
	// タグの重複チェック
	if slices.Contains(p.Tags, tag) {
		return myerror.NewMyError(myerror.InvalidRequestCode, "Tag already exists")
	}

	// 最大タグ数チェック
	if len(p.Tags) >= 10 {
		return myerror.NewMyError(myerror.InvalidRequestCode, "Maximum number of tags reached")
	}

	p.Tags = append(p.Tags, tag)
	return nil
}

/**
* ステータス遷移許容
* 下書き->（公開、削除）
* 公開->（非公開）
* 非公開->（公開、削除）
 */
func (p *Post) Publish() error {
	if p.Status != valueobject.StatusDraft && p.Status != valueobject.StatusPrivate {
		return myerror.NewMyError(myerror.InvalidRequestCode, "Only draft and private posts can be published")
	}

	p.Status = valueobject.StatusPublished

	if p.FirstPublishedAt.IsZero() {
		p.FirstPublishedAt = time.Now()
	}

	p.Status = valueobject.StatusPublished
	p.ContentUpdatedAt = time.Now()
	return nil
}

func (p *Post) Private() error {
	if p.Status != valueobject.StatusPublished {
		return myerror.NewMyError(myerror.InvalidRequestCode, "Only published posts can be private")
	}

	p.Status = valueobject.StatusPrivate
	p.ContentUpdatedAt = time.Now()
	return nil
}

func (p *Post) Delete() error {
	if p.Status != valueobject.StatusDraft && p.Status != valueobject.StatusPrivate {
		return myerror.NewMyError(myerror.InvalidRequestCode, "Only draft and private posts can be deleted")
	}

	p.Status = valueobject.StatusDeleted
	p.ContentUpdatedAt = time.Now()
	return nil
}
