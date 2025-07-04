package valueobject

type PostStatus string

const (
	StatusDraft     PostStatus = "draft"
	StatusPublished PostStatus = "published"
	StatusPrivate   PostStatus = "private"
	StatusDeleted   PostStatus = "deleted"
)

func NewPostStatus(status string) (PostStatus, error) {
	switch PostStatus(status) {
	case StatusDraft, StatusPublished, StatusPrivate, StatusDeleted:
		return PostStatus(status), nil
	default:
		return PostStatus(""), NewMyError(InvalidCode, "Invalid post status")
	}
}

func (ps PostStatus) String() string {
	return string(ps)
}

func (ps PostStatus) Equals(other PostStatus) bool {
	return ps == other
}
