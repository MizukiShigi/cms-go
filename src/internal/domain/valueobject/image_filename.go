package valueobject

type ImageFilename string

func NewImageFilename(filename string) (ImageFilename, error) {
	if len(filename) == 0 || len(filename) > 255 {
		return ImageFilename(""), NewMyError(InvalidCode, "filename length invalid")
	}

	return ImageFilename(filename), nil
}

func (f ImageFilename) String() string {
	return string(f)
}

func (f ImageFilename) Equals(other ImageFilename) bool {
	return f == other
}
