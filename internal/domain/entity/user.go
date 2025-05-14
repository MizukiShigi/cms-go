package entity

import (
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/myerror"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        valueobject.UserID
	Name      string
	Email     valueobject.Email
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name string, email valueobject.Email, password string) (*User, error) {
	if name == "" {
		return nil, myerror.NewMyError(myerror.InvalidCode, "Name cannot be empty")
	}

	if len(password) < 8 {
		return nil, myerror.NewMyError(myerror.InvalidCode, "Password must be at least 8 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InternalServerErrorCode, "Failed to hash password")
	}

	now := time.Now()
	user := &User{
		ID:        valueobject.NewUserID(),
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	return user, nil
}

func (u *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
