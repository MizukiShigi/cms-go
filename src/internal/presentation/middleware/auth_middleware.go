package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/MizukiShigi/cms-go/internal/presentation/helper"
)

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authorization ヘッダーからトークンを抽出
			tokenString, err := extractTokenFromHeader(r)
			if err != nil {
				myerr := valueobject.NewMyError(valueobject.UnauthorizedCode, err.Error())
				helper.RespondWithError(w, myerr)
				return
			}

			// JWT トークンを解析・検証
			claims, err := validateToken(tokenString, secretKey)
			if err != nil {
				myerr := valueobject.NewMyError(valueobject.UnauthorizedCode, err.Error())
				helper.RespondWithError(w, myerr)
				return
			}

			claimUserID, ok := claims["user_id"].(string)
			if !ok {
				myerr := valueobject.NewMyError(valueobject.UnauthorizedCode, "user_id is not a string")
				helper.RespondWithError(w, myerr)
				return
			}

			// 検証されたクレームをコンテキストに追加
			ctx := context.WithValue(r.Context(), domaincontext.UserID, claimUserID)

			// 検証されたクレームをログコンテキストに追加
			ctx = domaincontext.WithValue(ctx, "user_id", claimUserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractTokenFromHeader は Authorization ヘッダーから JWT トークンを抽出
func extractTokenFromHeader(r *http.Request) (string, error) {
	// Authorization ヘッダーを取得
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	// Bearer プレフィックスを確認
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format. 'Bearer <token>' format is required")
	}

	return parts[1], nil
}

// validateToken は JWT トークンを検証
func validateToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	// トークンを解析
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムを確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.New("failed to parse token")
	}

	// トークンの有効性を確認
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// クレームを取得
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to get token claims")
	}

	return claims, nil
}
