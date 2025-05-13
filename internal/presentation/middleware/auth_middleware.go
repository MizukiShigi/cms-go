package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
	"github.com/MizukiShigi/cms-go/internal/domain/myerror"
	"github.com/MizukiShigi/cms-go/internal/presentation/helper"
)

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authorization ヘッダーからトークンを抽出
			tokenString, err := extractTokenFromHeader(r)
			if err != nil {
				myerr := myerror.NewMyError(myerror.UnauthorizedCode, err.Error())
				helper.RespondWithError(w, myerr)
				return
			}

			// JWT トークンを解析・検証
			claims, err := validateToken(tokenString, secretKey)
			if err != nil {
				myerr := myerror.NewMyError(myerror.UnauthorizedCode, err.Error())
				helper.RespondWithError(w, myerr)
				return
			}

			// 検証されたクレームをログコンテキストに追加
			ctx := domaincontext.WithValue(r.Context(), "user_id", claims["user_id"])

			// 認証済みリクエストで次のハンドラーを呼び出す
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractTokenFromHeader は Authorization ヘッダーから JWT トークンを抽出
func extractTokenFromHeader(r *http.Request) (string, error) {
	// Authorization ヘッダーを取得
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("認証ヘッダーがありません")
	}

	// Bearer プレフィックスを確認
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("認証ヘッダーの形式が無効です。'Bearer <token>' 形式が必要です")
	}

	return parts[1], nil
}

// validateToken は JWT トークンを検証
func validateToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	// トークンを解析
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムを確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("予期しない署名方法")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.New("トークンの解析に失敗しました")
	}

	// トークンの有効性を確認
	if !token.Valid {
		return nil, errors.New("無効なトークンです")
	}

	// クレームを取得
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("トークンクレームの取得に失敗しました")
	}

	return claims, nil
}
