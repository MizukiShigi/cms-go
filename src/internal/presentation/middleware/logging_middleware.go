package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// リクエストIDとコンテキスト設定
		ctx := domaincontext.WithValue(r.Context(), "request_id", uuid.New().String())
		ctx = domaincontext.WithValue(ctx, "method", r.Method)
		ctx = domaincontext.WithValue(ctx, "url", r.URL.String())
		ctx = domaincontext.WithValue(ctx, "remote_addr", r.RemoteAddr)
		ctx = domaincontext.WithValue(ctx, "user_agent", r.Header.Get("User-Agent"))

		// リクエストボディの読み取りと復元
		var requestBody string
		if shouldLogBody(r) {
			body, err := readAndRestoreBody(r)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to read request body", "error", err)
			} else {
				// 機密情報をマスキング
				maskedBody := maskSensitiveData(body)
				requestBody = maskedBody
				ctx = domaincontext.WithValue(ctx, "request_body", requestBody)
			}
		}

		// レスポンスキャプチャのためのラッパー
		capture := &ResponseCapture{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// リクエストログ
		slog.InfoContext(ctx, "request started")

		// 次のハンドラーを実行
		next.ServeHTTP(capture, r.WithContext(ctx))

		// 処理時間計算
		duration := time.Since(startTime)

		// レスポンス情報をコンテキストに追加
		ctx = domaincontext.WithValue(ctx, "status_code", capture.statusCode)
		ctx = domaincontext.WithValue(ctx, "response_size", capture.body.Len())
		ctx = domaincontext.WithValue(ctx, "duration_ms", duration.Milliseconds())

		// レスポンスボディをログに追加（必要な場合）
		if shouldLogResponseBody(capture) {
			responseBody := maskSensitiveData(capture.body.String())
			ctx = domaincontext.WithValue(ctx, "response_body", responseBody)
		}

		// レスポンスログ
		slog.InfoContext(ctx, "request completed")
	})
}

// ResponseCapture はレスポンスの内容をキャプチャするためのラッパー
type ResponseCapture struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// Write はレスポンスボディをキャプチャしながら書き込む
func (rc *ResponseCapture) Write(data []byte) (int, error) {
	// ボディをキャプチャ
	rc.body.Write(data)
	// 元のResponseWriterにも書き込み
	return rc.ResponseWriter.Write(data)
}

// WriteHeader はステータスコードをキャプチャ
func (rc *ResponseCapture) WriteHeader(statusCode int) {
	rc.statusCode = statusCode
	rc.ResponseWriter.WriteHeader(statusCode)
}

// readAndRestoreBody はリクエストボディを読み取り、読み取り可能な状態に復元する
func readAndRestoreBody(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", nil
	}

	// ボディを読み取り
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	// ボディを復元
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// 文字列として返却（適切な場合のみ）
	bodyStr := string(bodyBytes)

	// サイズ制限（10KB）
	const maxBodySize = 10 * 1024
	if len(bodyStr) > maxBodySize {
		return bodyStr[:maxBodySize] + "...[truncated]", nil
	}

	return bodyStr, nil
}

// shouldLogBody はリクエストボディをログに記録すべきかを判定
func shouldLogBody(r *http.Request) bool {
	// GETリクエストはボディがないため除外
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		return false
	}

	// Content-Typeをチェック
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return true // Content-Typeが不明でも一応ログ
	}

	// テキスト系のContent-Typeのみログ対象
	textContentTypes := []string{
		"application/json",
		"application/xml",
		"text/",
		"application/x-www-form-urlencoded",
	}

	for _, textType := range textContentTypes {
		if strings.Contains(strings.ToLower(contentType), textType) {
			return true
		}
	}

	return false // バイナリデータは除外
}

// shouldLogResponseBody はレスポンスボディをログに記録すべきかを判定
func shouldLogResponseBody(capture *ResponseCapture) bool {
	// エラー系のステータスコードは詳細ログ
	if capture.statusCode >= 400 {
		return true
	}

	// ボディサイズが大きすぎる場合は除外
	const maxResponseSize = 5 * 1024 // 5KB
	if capture.body.Len() > maxResponseSize {
		return false
	}

	// Content-Typeをチェック（可能であれば）
	contentType := capture.Header().Get("Content-Type")
	if contentType != "" && !strings.Contains(strings.ToLower(contentType), "json") &&
		!strings.Contains(strings.ToLower(contentType), "text") {
		return false
	}

	return true
}

// maskSensitiveData は機密情報をマスキングする
func maskSensitiveData(body string) string {
	if body == "" {
		return body
	}

	// パスワード関連のフィールドをマスキング
	sensitiveFields := []string{
		"password",
		"passwd",
		"pwd",
		"secret",
		"token",
		"authorization",
		"auth",
		"key",
		"credential",
	}

	result := body
	for _, field := range sensitiveFields {
		// JSON形式のフィールドを単純に置換
		if strings.Contains(strings.ToLower(result), `"`+strings.ToLower(field)+`":`) {
			result = strings.ReplaceAll(result,
				`"`+field+`":`, `"`+field+`": "***" // masked`)
		}
	}

	return result
}
