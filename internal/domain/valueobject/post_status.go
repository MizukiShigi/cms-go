package valueobject

import "github.com/MizukiShigi/cms-go/internal/domain/myerror"

type PostStatus string

const (
	StatusDraft     PostStatus = "DRAFT"
	StatusPublished PostStatus = "PUBLISHED"
	StatusPrivate   PostStatus = "PRIVATE"
	StatusDeleted   PostStatus = "DELETED"
)

func NewPostStatus(status string) (PostStatus, error) {
	switch PostStatus(status) {
	case StatusDraft, StatusPublished, StatusPrivate, StatusDeleted:
		return PostStatus(status), nil
	default:
		return PostStatus(""), myerror.NewMyError(myerror.InvalidRequestCode, "Invalid post status")
	}
}

func (ps PostStatus) String() string {
	return string(ps)
}

func (ps PostStatus) Equals(other PostStatus) bool {
	return ps == other
}