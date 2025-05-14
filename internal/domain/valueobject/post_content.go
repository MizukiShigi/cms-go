package valueobject

import "github.com/MizukiShigi/cms-go/internal/domain/myerror"

type PostContent string

func NewPostContent(content string) (PostContent, error) {
	if len(content) > 10000 {
		return PostContent(""), myerror.NewMyError(myerror.InvalidCode, "Content is too long")
	}
	return PostContent(content), nil
}

func (p PostContent) String() string {
	return string(p)
}

func (p PostContent) Equals(other PostContent) bool {
	return p == other
}
