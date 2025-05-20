package entity

import (
	"slices"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type Post struct {
	ID               valueobject.PostID
	Title            valueobject.PostTitle
	Content          valueobject.PostContent
	UserID           valueobject.UserID
	Status           valueobject.PostStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
	FirstPublishedAt *time.Time
	ContentUpdatedAt *time.Time
	Tags             []valueobject.TagName
}

// 新規投稿作成
func NewPost(title valueobject.PostTitle, content valueobject.PostContent, userID valueobject.UserID, status valueobject.PostStatus) (*Post, error) {
	now := time.Now()
	var firstPublishedAt time.Time
	if status == valueobject.StatusPublished {
		firstPublishedAt = now
	}

	post := &Post{
		ID:               valueobject.NewPostID(),
		Title:            title,
		Content:          content,
		UserID:           userID,
		Status:           status,
		FirstPublishedAt: &firstPublishedAt,
		ContentUpdatedAt: &now,
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
	userID valueobject.UserID,
	status valueobject.PostStatus,
	createdAt time.Time,
	updatedAt time.Time,
	firstPublishedAt *time.Time,
	contentUpdatedAt *time.Time,
	tags []valueobject.TagName,
) *Post {
	return &Post{
		ID:               id,
		Title:            title,
		Content:          content,
		UserID:           userID,
		Status:           status,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		FirstPublishedAt: firstPublishedAt,
		ContentUpdatedAt: contentUpdatedAt,
		Tags:             tags,
	}
}

func (p *Post) AddTag(tag valueobject.TagName) error {
	// タグの重複チェック
	if slices.Contains(p.Tags, tag) {
		return valueobject.NewMyError(valueobject.InvalidCode, "Tag already exists")
	}

	// 最大タグ数チェック
	if len(p.Tags) >= 10 {
		return valueobject.NewMyError(valueobject.InvalidCode, "Maximum number of tags reached")
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
func (p *Post) SetStatus(status valueobject.PostStatus) error {
	if p.Status == status {
		return nil
	}	

	switch status {
	case valueobject.StatusDraft:
		return p.draft()
	case valueobject.StatusPublished:
		return p.publish()
	case valueobject.StatusPrivate:
		return p.private()
	case valueobject.StatusDeleted:
		return p.delete()
	default:
		return valueobject.NewMyError(valueobject.InvalidCode, "Invalid post status")
	}
}

func (p *Post) draft() error {
	if p.Status != valueobject.StatusDraft {
		return valueobject.NewMyError(valueobject.InvalidCode, "Only draft posts can be drafted")
	}

	p.Status = valueobject.StatusDraft
	*p.ContentUpdatedAt = time.Now()
	return nil
}

func (p *Post) publish() error {
	if p.Status != valueobject.StatusDraft && p.Status != valueobject.StatusPrivate {
		return valueobject.NewMyError(valueobject.InvalidCode, "Only draft and private posts can be published")
	}

	p.Status = valueobject.StatusPublished

	now := time.Now()
	if p.FirstPublishedAt == nil {
		p.FirstPublishedAt = &now
	}

	p.Status = valueobject.StatusPublished
	p.ContentUpdatedAt = &now
	return nil
}

func (p *Post) private() error {
	if p.Status != valueobject.StatusPublished {
		return valueobject.NewMyError(valueobject.InvalidCode, "Only published posts can be private")
	}

	p.Status = valueobject.StatusPrivate
	now := time.Now()
	p.ContentUpdatedAt = &now
	return nil
}

func (p *Post) delete() error {
	if p.Status != valueobject.StatusDraft && p.Status != valueobject.StatusPrivate {
		return valueobject.NewMyError(valueobject.InvalidCode, "Only draft and private posts can be deleted")
	}

	p.Status = valueobject.StatusDeleted
	now := time.Now()
	p.ContentUpdatedAt = &now
	return nil
}
