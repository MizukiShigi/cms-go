package middleware

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := domaincontext.WithValue(r.Context(), "request_id", uuid.New().String())
		ctx = domaincontext.WithValue(ctx, "method", r.Method)
		ctx = domaincontext.WithValue(ctx, "url", r.URL.String())
		slog.InfoContext(ctx, "request")
		next.ServeHTTP(w, r.WithContext(ctx))
		slog.InfoContext(ctx, "response")
	})
}
