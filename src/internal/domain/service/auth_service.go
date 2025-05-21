package service

import (
	"context"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

type AuthService interface {
	GenerateToken(ctx context.Context, userID valueobject.UserID, email valueobject.Email) (string, error)
}
