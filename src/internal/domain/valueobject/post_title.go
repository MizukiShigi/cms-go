package valueobject

import (
	"fmt"
	"html"
	"strings"
)

type PostTitle string

func NewPostTitle(title string) (PostTitle, error) {
	normalizedTitle := strings.TrimSpace(title)
	normalizedTitleRune := []rune(normalizedTitle)

	if len(normalizedTitleRune) > 200 {
		return PostTitle(""), NewMyError(InvalidCode, "Title is too long")
	}

	sanitizedTitle := html.EscapeString(normalizedTitle)
	sanitizedTitleRune := []rune(sanitizedTitle)
	if len(sanitizedTitleRune) > 255 {
		return PostTitle(""), NewMyError(InvalidCode, "Title contains too many special characters that expand when escaped")
	}

	forbiddenChars := []rune{'<', '>', '"', '\'', '\\', '/', '\n', '\r', '\t'}
	for _, char := range forbiddenChars {
		if strings.ContainsRune(sanitizedTitle, char) {
			return PostTitle(""), NewMyError(InvalidCode, fmt.Sprintf("Title contains forbidden character: %q", char))
		}
	}

	return PostTitle(sanitizedTitle), nil
}

func (p PostTitle) Value() string {
	return html.UnescapeString(string(p))
}

func (p PostTitle) String() string {
	return string(p)
}

func (p PostTitle) Equals(other PostTitle) bool {
	return p == other
}
