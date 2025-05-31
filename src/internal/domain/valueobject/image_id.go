package valueobject

import (
	"github.com/google/uuid"
)

type ImageID string

func NewImageID() ImageID {
	return ImageID(uuid.New().String())
}

func (p ImageID) String() string {
	return string(p)
}

func (p ImageID) Equals(other ImageID) bool {
	return p == other
}

func ParseImageID(s string) (ImageID, error) {
	uuid, err := uuid.Parse(s)
	if err != nil {
		return ImageID(""), NewMyError(InvalidCode, "Invalid image ID")
	}

	return ImageID(uuid.String()), nil
}
