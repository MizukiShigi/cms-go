package logger

import (
	"context"
	"log/slog"
	"sync"

	domaincontext "github.com/MizukiShigi/cms-go/internal/domain/context"
)

type Handler struct {
	handler slog.Handler
}

func NewHandler(handler slog.Handler) slog.Handler {
	return Handler{
		handler: handler,
	}
}

func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	if v, ok := ctx.Value(domaincontext.Fields).(*sync.Map); ok {
		v.Range(func(k, v any) bool {
			if k, ok := k.(string); ok {
				record.AddAttrs(slog.Any(k, v))
			}
			return true
		})
	}
	return h.handler.Handle(ctx, record)
}

func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{
		handler: h.handler.WithAttrs(attrs),
	}
}

func (h Handler) WithGroup(name string) slog.Handler {
	return Handler{
		handler: h.handler.WithGroup(name),
	}
}
