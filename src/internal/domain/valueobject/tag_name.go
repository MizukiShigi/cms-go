package valueobject

import (
	"regexp"
	"strings"
)

type TagName string

func NewTagName(tag string) (TagName, error) {
	normalizedTag := strings.ToLower(strings.TrimSpace(tag))

	if len(normalizedTag) > 50 {
		return TagName(""), NewMyError(InvalidCode, "TagName is too long")
	}

	validPattern := regexp.MustCompile(`^[a-z0-9\-_]+$`)
	if !validPattern.MatchString(normalizedTag) {
		return TagName(""), NewMyError(InvalidCode, "TagName can only contain lowercase letters, numbers, hyphens, and underscores")
	}

	return TagName(normalizedTag), nil
}

func (t TagName) String() string {
	return string(t)
}

func (t TagName) Equals(other TagName) bool {
	return t == other
}
