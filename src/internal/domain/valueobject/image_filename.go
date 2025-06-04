package valueobject

import (
	"path/filepath"
	"slices"
	"strings"
)

type ImageFilename string

func NewImageFilename(filename string) (ImageFilename, error) {
	if err := validateFilename(filename); err != nil {
		return ImageFilename(""), err
	}

	return ImageFilename(filename), nil
}

func validateFilename(name string) error {
	if len(name) == 0 || len(name) > 255 {
		return NewMyError(InvalidCode, "filename length invalid")
	}

	// 禁止文字チェック
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return NewMyError(InvalidCode, "filename contains invalid characters")
		}
	}

	// 画像ファイルの拡張子チェック
	ext := strings.ToLower(filepath.Ext(name))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	if slices.Contains(allowedExts, ext) {
		return nil
	}

	return NewMyError(InvalidCode, "unsupported file extension")
}

func (f ImageFilename) String() string {
	return string(f)
}

func (f ImageFilename) Equals(other ImageFilename) bool {
	return f == other
}
