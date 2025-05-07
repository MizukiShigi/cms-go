package valueobject

import (
	"fmt"
	"html"
	"strings"

	"github.com/MizukiShigi/cms-go/internal/domain/myerror"
)

type PostTitle string

func NewPostTitle(title string) (PostTitle, error) {
	normalizedTitle := strings.TrimSpace(title)

	if len(normalizedTitle) > 200 {
		return PostTitle(""), myerror.NewMyError(myerror.InvalidRequestCode, "Title is too long")
	}

	sanitizedTitle := html.EscapeString(normalizedTitle)

	if len(sanitizedTitle) > 255 {
		return PostTitle(""), myerror.NewMyError(myerror.InvalidRequestCode, "Title contains too many special characters that expand when escaped")
	}

	forbiddenChars := []rune{'<', '>', '"', '\'', '\\', '/', '\n', '\r', '\t'}
	for _, char := range forbiddenChars {
		if strings.ContainsRune(sanitizedTitle, char) {
			return PostTitle(""), myerror.NewMyError(myerror.InvalidRequestCode, fmt.Sprintf("Title contains forbidden character: %q", char))
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
