// Package middleware は、HTTPリクエストの前処理を行うミドルウェア機能を提供します。
package middleware

import (
	"net/http"
	"os"
)

// CORSMiddleware はCORS設定を行うミドルウェア
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 環境に応じてオリジンを設定
		env := os.Getenv("ENV")
		var allowedOrigin string
		
		switch env {
		case "local", "development":
			// 開発環境では全オリジンを許可
			allowedOrigin = "*"
		case "production":
			// 本番環境では特定のオリジンのみ許可
			allowedOrigin = os.Getenv("FRONTEND_URL")
			if allowedOrigin == "" {
				allowedOrigin = "https://your-frontend-domain.com"
			}
		default:
			// デフォルトは開発環境として扱う
			allowedOrigin = "*"
		}
		
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		
		// 許可するHTTPメソッドを設定
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		
		// 許可するヘッダーを設定
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// 認証情報の送信を許可
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// プリフライトリクエストの場合は200を返す
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// 次のハンドラーに処理を渡す
		next.ServeHTTP(w, r)
	})
}