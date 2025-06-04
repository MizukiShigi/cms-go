package service

import (
	"context"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
	secretKey string
	expiresAt time.Duration
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		expiresAt: time.Hour * 24,
	}
}

func (js *JWTService) GenerateToken(ctx context.Context, userID valueobject.UserID, email valueobject.Email) (string, error) {
	claims := &Claims{
		UserID: userID.String(),
		Email:  email.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.expiresAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(js.secretKey))
}
