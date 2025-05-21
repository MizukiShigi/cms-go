package valueobject

import (
	"github.com/google/uuid"
)

type UserID string

func NewUserID() UserID {
	return UserID(uuid.New().String())
}

func ParseUserID(s string) (UserID, error) {
	uuid, err := uuid.Parse(s)
	if err != nil {
		return UserID(""), NewMyError(InvalidCode, "Invalid user ID")
	}
	return UserID(uuid.String()), nil
}

func (u UserID) String() string {
	return string(u)
}

func (u UserID) Equals(other UserID) bool {
	return u == other
}
