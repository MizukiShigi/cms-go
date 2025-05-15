package valueobject

import (
	"regexp"
	"strings"
)

type Tag string

func NewTag(tag string) (Tag, error) {
	normalizedTag := strings.ToLower(strings.TrimSpace(tag))

	if len(normalizedTag) > 50 {
		return Tag(""), NewMyError(InvalidCode, "Tag is too long")
	}

	validPattern := regexp.MustCompile(`^[a-z0-9\-_]+$`)
	if !validPattern.MatchString(normalizedTag) {
		return Tag(""), NewMyError(InvalidCode, "Tag can only contain lowercase letters, numbers, hyphens, and underscores")
	}

	return Tag(normalizedTag), nil
}

func (t Tag) String() string {
	return string(t)
}

func (t Tag) Equals(other Tag) bool {
	return t == other
}
