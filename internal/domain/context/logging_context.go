package context

import (
	"context"
	"sync"
)

type loggingContextKey string

const Fields loggingContextKey = "logging_context"

func WithValue(parent context.Context, key string, val any) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	if v, ok := parent.Value(Fields).(*sync.Map); ok {
		mapCopy := copySyncMap(v)
		mapCopy.Store(key, val)
		return context.WithValue(parent, Fields, mapCopy)
	}
	v := &sync.Map{}
	v.Store(key, val)
	return context.WithValue(parent, Fields, v)
}

func copySyncMap(m *sync.Map) *sync.Map {
	var cp sync.Map
	m.Range(func(k, v any) bool {
		cp.Store(k, v)
		return true
	})
	return &cp
}
