package valueobject

import (
	"regexp"
)

type Email string

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func NewEmail(email string) (Email, error) {
	if !emailRegex.MatchString(email) {
		return Email(""), NewMyError(InvalidCode, "Invalid email format")
	}
	return Email(email), nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Equals(other Email) bool {
	return e == other
}
