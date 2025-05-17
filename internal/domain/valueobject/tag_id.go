package valueobject

import "github.com/google/uuid"

type TagID string

func NewTagID() TagID {
	uuid := uuid.New().String()
	return TagID(uuid)
}

func (t TagID) String() string {
	return string(t)
}
