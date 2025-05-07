package valueobject

import "github.com/google/uuid"

type PostID string

func NewPostID() PostID {
	return PostID(uuid.New().String())
}

func (p PostID) String() string {
	return string(p)
}

func (p PostID) Equals(other PostID) bool {
	return p == other
}

func ParsePostID(s string) (PostID, error) {
	uuid, err := uuid.Parse(s)
	if err != nil {
		return PostID(""), err
	}

	return PostID(uuid.String()), nil
}
